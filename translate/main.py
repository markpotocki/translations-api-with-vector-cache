from fastapi import FastAPI
import uvicorn
from pydantic import BaseModel
from translator import ArgosTranslator


class TranslationRequest(BaseModel):
    text: str
    target_language: str
    source_language: str


app = FastAPI()

translator = ArgosTranslator(
    packages=[
        ("en", "es"),  # English to Spanish
        ("es", "en"),  # Spanish to English
        ("fr", "en"),  # French to English
        ("en", "fr"),  # English to French
        ("de", "en"),  # German to English
        ("en", "de"),   # English to German
        ("en", "zh")
    ]
)


@app.post("/translate")
async def translate(request: TranslationRequest):
    text = translator.translate(
        text=request.text,
        target_language=request.target_language,
        source_language=request.source_language
    )

    return {"translation": text}


if __name__ == "__main__":
    uvicorn.run(app)