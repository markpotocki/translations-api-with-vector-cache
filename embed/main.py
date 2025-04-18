from embedder import Embedder
from fastapi import FastAPI
import uvicorn
from pydantic import BaseModel


class EmbeddingRequest(BaseModel):
    text: str


class EmbeddingResponse(BaseModel):
    embedding: list[float]


app = FastAPI()
embedder = Embedder()


@app.post("/embed", response_model=EmbeddingResponse)
async def embed(request: EmbeddingRequest):
    embedding = embedder.embed(request.text)
    return EmbeddingResponse(embedding=embedding.tolist())


if __name__ == "__main__":
    uvicorn.run(app)
