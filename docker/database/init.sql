CREATE EXTENSION vector;

CREATE TABLE IF NOT EXISTS translations_cache (
    id SERIAL PRIMARY KEY,
    source_language text NOT NULL,
    target_language text NOT NULL,
    source_text text NOT NULL,
    target_text text NOT NULL,
    embedding vector(384)
);

CREATE INDEX idx_translations_cache_embedding ON translations_cache USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
