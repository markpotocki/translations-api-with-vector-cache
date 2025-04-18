from embedder import Embedder
import grpc
from concurrent import futures
from embed_pb2 import EmbeddingResponse
from embed_pb2_grpc import EmbedderServicer, add_EmbedderServicer_to_server
import logging


logging.basicConfig(level=logging.DEBUG)


class EmbedderService(EmbedderServicer):
    def __init__(self, embedder: Embedder):
        self.embedder = embedder

    def GenerateEmbedding(self, request, context):
        embedding = self.embedder.embed(request.text)
        return EmbeddingResponse(embedding=embedding.tolist())


def serve():
    logging.info("Starting gRPC server...")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    logging.info("Creating Embedder instance...")
    embedder = Embedder()
    add_EmbedderServicer_to_server(EmbedderService(embedder), server)
    server.add_insecure_port('[::]:50051')
    logging.info("Starting server on port 50051...")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.info("Starting embedding service...")
    serve()
