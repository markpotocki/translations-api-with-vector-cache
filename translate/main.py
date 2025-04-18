import grpc
from concurrent import futures
from translate_pb2 import TranslationRequest, TranslationResponse
from translate_pb2_grpc import TranslatorServicer, add_TranslatorServicer_to_server
from translator import ArgosTranslator
import logging


logging.basicConfig(level=logging.DEBUG)


class TranslatorService(TranslatorServicer):
    def __init__(self, translator: ArgosTranslator):
        self.translator = translator

    def Translate(self, request: TranslationRequest, context):
        text = self.translator.translate(
            text=request.text,
            target_language=request.target_language,
            source_language=request.source_language
        )
        return TranslationResponse(translation=text)
    

def serve():
    logging.info("Starting gRPC server...")
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))

    logging.info("Creating ArgosTranslator instance...")
    translator = ArgosTranslator(
        packages=[
            ("en", "es"),  # English to Spanish
            ("es", "en"),  # Spanish to English
            ("fr", "en"),  # French to English
            ("en", "fr"),  # English to French
            ("de", "en"),  # German to English
            ("en", "de"),   # English to German
            ("en", "zh")   # English to Chinese
        ]
    )
    add_TranslatorServicer_to_server(TranslatorService(translator), server)
    server.add_insecure_port('[::]:50051')
    logging.info("Starting server on port 50051...")
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.info("Starting translation service...")
    serve()
