package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"

	embedpb "service/embeddingapi/service"
	translatepb "service/translationsapi/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	conn            *pgx.Conn
	embeddingURL    string
	translateURL    string
	embedConn       *grpc.ClientConn
	translateConn   *grpc.ClientConn
	translateClient translatepb.TranslatorClient
	embedClient     embedpb.EmbedderClient
)

func init() {
	// Load environment variables
	embeddingURL = getEnv("EMBEDDING_URL")
	translateURL = getEnv("TRANSLATE_URL")
	connStr := getEnv("DATABASE_URL")

	// Connect to the database
	var err error
	conn, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	// Register pgvector types
	if err := pgxvec.RegisterTypes(context.TODO(), conn); err != nil {
		log.Fatalf("Failed to register pgvector types: %v\n", err)
	}

	// Create grpc clients
	translateConn, err = grpc.NewClient(translateURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("error connecting to translation service: %v", err)
	}
	translateClient = translatepb.NewTranslatorClient(translateConn)

	embedConn, err = grpc.NewClient(embeddingURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Panicf("error connecting to embedding service: %v", err)
	}
	embedClient = embedpb.NewEmbedderClient(embedConn)
	log.Println("Service initialized successfully")
}

func main() {
	defer conn.Close(context.Background())
	defer translateConn.Close()
	defer embedConn.Close()

	http.HandleFunc("/translate", handleTranslate)

	port := getEnvWithDefault("PORT", "8080")
	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}

// handleTranslate handles the /translate endpoint
func handleTranslate(w http.ResponseWriter, r *http.Request) {
	defer recoverFromPanic(w)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Text           string `json:"text"`
		SourceLanguage string `json:"source_language"`
		TargetLanguage string `json:"target_language"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil || request.Text == "" || request.SourceLanguage == "" || request.TargetLanguage == "" {
		http.Error(w, "Invalid or missing fields in request body", http.StatusBadRequest)
		return
	}

	responseText, err := processTranslation(request.Text, request.SourceLanguage, request.TargetLanguage)
	if err != nil {
		log.Printf("Error processing translation: %v", err)
		http.Error(w, "Error processing translation", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, http.StatusOK, map[string]string{"translation": responseText})
}

// processTranslation handles the translation logic
func processTranslation(text, sourceLang, targetLang string) (string, error) {
	sentences := strings.Split(text, ". ")
	embeddings := make([][]float32, len(sentences))
	var translations []string

	// Get embeddings for each sentence
	for i, sentence := range sentences {
		embedding, err := getEmbedding(sentence)
		if err != nil {
			return "", fmt.Errorf("error getting embedding: %w", err)
		}
		embeddings[i] = embedding
	}

	// Process each sentence
	for i, sentence := range sentences {
		embedding := embeddings[i]
		cachedTranslation, found, err := getFromCache(conn, sourceLang, targetLang, embedding)
		if err != nil {
			return "", fmt.Errorf("error accessing cache: %w", err)
		}

		if found {
			log.Printf("Using cached translation for: %s", sentence)
			translations = append(translations, cachedTranslation)
		} else {
			log.Printf("No cache found for: %s, fetching translation", sentence)
			translation, err := getTranslation(sentence, sourceLang, targetLang)
			if err != nil {
				return "", fmt.Errorf("error getting translation: %w", err)
			}

			if err := saveToCache(conn, sourceLang, targetLang, embedding, translation, sentence); err != nil {
				return "", fmt.Errorf("error saving to cache: %w", err)
			}
			translations = append(translations, translation)
		}
	}

	return strings.Join(translations, ". "), nil
}

// getEmbedding fetches the embedding for a given text
func getEmbedding(text string) ([]float32, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	// Create the request for embedding
	req := &embedpb.EmbeddingRequest{
		Text: text,
	}
	res, err := embedClient.GenerateEmbedding(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling embedding service: %w", err)
	}

	return res.Embedding, nil
}

// getTranslation fetches the translation for a given text
func getTranslation(text, sourceLang, targetLang string) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	req := &translatepb.TranslationRequest{
		Text:           text,
		SourceLanguage: sourceLang,
		TargetLanguage: targetLang,
	}

	res, err := translateClient.Translate(ctx, req)
	if err != nil {
		return "", fmt.Errorf("error calling translation service: %w", err)
	}

	return res.Translation, nil
}

// getFromCache retrieves a cached translation from the database
func getFromCache(conn *pgx.Conn, sourceLang, targetLang string, embedding []float32) (string, bool, error) {
	query := `
        SELECT target_text
        FROM translations_cache
        WHERE source_language = $1
        AND target_language = $2
        AND embedding <=> $3 <= 0.1
        LIMIT 1;
    `

	var targetText string
	err := conn.QueryRow(context.TODO(), query, sourceLang, targetLang, pgvector.NewVector(embedding)).Scan(&targetText)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", false, nil
		}
		return "", false, err
	}
	return targetText, true, nil
}

// saveToCache saves a translation to the database
func saveToCache(conn *pgx.Conn, sourceLang, targetLang string, embedding []float32, targetText, sourceText string) error {
	query := `
        INSERT INTO translations_cache (source_language, target_language, embedding, target_text, source_text)
        VALUES ($1, $2, $3, $4, $5);
    `
	_, err := conn.Exec(context.TODO(), query, sourceLang, targetLang, pgvector.NewVector(embedding), targetText, sourceText)
	return err
}

// Utility Functions

// recoverFromPanic handles panics and sends an error response
func recoverFromPanic(w http.ResponseWriter) {
	if r := recover(); r != nil {
		log.Printf("Recovered from panic: %v", r)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// writeJSONResponse writes a JSON response to the client
func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// getEnv retrieves an environment variable or panics if not set
func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("%s environment variable is not set", key))
	}
	return value
}

// getEnvWithDefault retrieves an environment variable or returns a default value
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
