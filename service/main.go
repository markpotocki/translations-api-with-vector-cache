package main

import (
	"bytes"
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
)

var (
	client       *http.Client
	conn         *pgx.Conn
	embeddingURL string
	translateURL string
)

func init() {
	// Load environment variables
	embeddingURL = getEnv("EMBEDDING_URL")
	translateURL = getEnv("TRANSLATE_URL")
	connStr := getEnv("DATABASE_URL")

	// Initialize HTTP client
	client = &http.Client{Timeout: 30 * time.Second}

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
}

func main() {
	defer conn.Close(context.Background())

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

type embeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// getEmbedding fetches the embedding for a given text
func getEmbedding(text string) ([]float32, error) {
	request := map[string]string{"text": text}
	response, err := postJSON[embeddingResponse](embeddingURL, request)
	if err != nil {
		return nil, fmt.Errorf("error fetching embedding: %w", err)
	}
	return response.Embedding, nil
}

type translationResponse struct {
	Translation string `json:"translation"`
}

// getTranslation fetches the translation for a given text
func getTranslation(text, sourceLang, targetLang string) (string, error) {
	request := map[string]string{
		"text":            text,
		"source_language": sourceLang,
		"target_language": targetLang,
	}
	response, err := postJSON[translationResponse](translateURL, request)
	if err != nil {
		return "", fmt.Errorf("error fetching translation: %w", err)
	}
	return response.Translation, nil
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

// postJSON sends a POST request with a JSON body and decodes the response
func postJSON[T any](url string, requestBody interface{}) (T, error) {
	data, err := json.Marshal(requestBody)
	if err != nil {
		return *new(T), err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return *new(T), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return *new(T), fmt.Errorf("failed request: %s", resp.Status)
	}

	var response T
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return *new(T), err
	}

	return response, nil
}

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
