# Translations API

This repository provides a modular and scalable translation service built with multiple APIs, a PostgreSQL database with `pgvector` for vector similarity, and Docker Compose for orchestration. The service supports embedding generation, caching translations, and querying translations efficiently. It now includes gRPC-based APIs for translation and embedding services.

---

## Features

- **Translation API**: Handles text translation between languages via gRPC.
- **Embedding API**: Generates embeddings for text via gRPC.
- **Service API**: Orchestrates translation and caching logic.
- **PostgreSQL with `pgvector`**: Stores translations and embeddings for efficient similarity queries.
- **Docker Compose**: Orchestrates the APIs and database.

---

## Architecture

The system consists of the following components:

1. **Translate API**: Handles translation requests via gRPC.
2. **Embedding API**: Generates embeddings for input text via gRPC.
3. **Service API**: Combines translation and embedding logic, caches results in the database, and retrieves cached translations when possible.
4. **PostgreSQL Database**: Stores translations, embeddings, and supports similarity queries using `pgvector`.

### Architectural Diagram

```mermaid
graph TD
    A[Client] -->|gRPC Translate| B[Service API]
    B -->|gRPC Embedding| C[Embedding API]
    B -->|gRPC Translation| D[Translate API]
    B -->|SQL Query| E[PostgreSQL Database]
    E -->|Cached Translation| B
```

---

## Translation Logic Flow

The following flowchart describes the logic of the translation process:

```mermaid
flowchart TD
    A[Start: Client Request] --> B[Split Text into Sentences]
    B --> C[Generate Embedding for Each Sentence]
    C --> D{Check Cache for Translation}
    D -->|Found| E[Return Cached Translation]
    D -->|Not Found| F[Fetch Translation from Translate API]
    F --> G[Save Translation to Cache]
    G --> E
    E --> H[Combine Translations]
    H --> I[Return Response to Client]
```

---

## Getting Started

### Prerequisites

- Docker and Docker Compose installed on your system.
- Python 3.9+ installed for local development of the gRPC APIs.
- Go installed for local development of the Service API.

### Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/translationsapi-py.git
   cd translationsapi-py
   ```

2. Build and start the services:
   ```bash
   docker-compose up --build
   ```

3. Access the services:
   - **Service API**: `http://localhost:8080/translate`
   - **Translate API (gRPC)**: `localhost:50051`
   - **Embedding API (gRPC)**: `localhost:50052`
   - **PostgreSQL**: `localhost:5432`

---

## gRPC Endpoints

### Translate API

#### `Translate`
- **Description**: Translates text between languages.
- **Request**:
  ```proto
  message TranslationRequest {
    string text = 1;
    string source_language = 2;
    string target_language = 3;
  }
  ```
- **Response**:
  ```proto
  message TranslationResponse {
    string translation = 1;
  }
  ```

### Embedding API

#### `GenerateEmbedding`
- **Description**: Generates an embedding for the given text.
- **Request**:
  ```proto
  message EmbeddingRequest {
    string text = 1;
  }
  ```
- **Response**:
  ```proto
  message EmbeddingResponse {
    repeated float embedding = 1;
  }
  ```

---

## Database Schema

The `translations_cache` table is used to store translations and embeddings:

```sql
CREATE TABLE IF NOT EXISTS translations_cache (
    id SERIAL PRIMARY KEY,
    source_language TEXT NOT NULL,
    target_language TEXT NOT NULL,
    source_text TEXT NOT NULL,
    target_text TEXT NOT NULL,
    embedding VECTOR(384)
);

CREATE INDEX IF NOT EXISTS idx_translations_cache_embedding
ON translations_cache USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
```

---

## Environment Variables

The following environment variables are required:

| Variable          | Description                          | Default Value |
|--------------------|--------------------------------------|---------------|
| `DATABASE_URL`     | PostgreSQL connection string         | None          |
| `EMBEDDING_URL`    | URL for the embedding API (gRPC)     | None          |
| `TRANSLATE_URL`    | URL for the translation API (gRPC)   | None          |
| `PORT`             | Port for the Service API             | `8080`        |

---

## Development

### Running Locally

1. Install dependencies for the Translate API:
   ```bash
   cd translate
   pip install -r requirements.txt
   python main.py
   ```

2. Install dependencies for the Embedding API:
   ```bash
   cd embed
   pip install -r requirements.txt
   python main.py
   ```

3. Start the database and Service API using Docker Compose:
   ```bash
   docker-compose up db serviceapi
   ```

---

## Testing

You can test the gRPC APIs using tools like `grpcurl` or by writing a gRPC client in Python or Go.

### Example: Using `grpcurl` to Test the Translate API
```bash
grpcurl -plaintext -d '{
  "text": "Hello world",
  "source_language": "en",
  "target_language": "es"
}' localhost:50051 translate.Translator/Translate
```

### Example: Using `grpcurl` to Test the Embedding API
```bash
grpcurl -plaintext -d '{
  "text": "Hello world"
}' localhost:50052 embed.Embedder/GenerateEmbedding
```

---

## Future Improvements

- Add support for more languages.
- Implement rate limiting for APIs.
- Add monitoring and logging for better observability.
- Add TLS support for gRPC endpoints.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
