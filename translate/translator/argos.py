import argostranslate.package, argostranslate.translate


class ArgosTranslator:
    def __init__(self, packages: list[tuple[str, str]] = [("en", "es")]):
        for source, target in packages:
            print(f"Installing translation package from {source} to {target}")
            self.download_language_package(source, target)

        
    def translate(self, text: str, target_language: str, source_language: str) -> str:
        installed_languages = argostranslate.translate.get_installed_languages()
        print(f"Installed languages: {[lang.code for lang in installed_languages]}")

        from_lang = list(filter(
            lambda x: x.code == source_language, installed_languages
        ))
        if len(from_lang) == 0:
            raise ValueError("Source language not installed")
        from_lang = from_lang[0]

        to_lang = list(filter(
            lambda x: x.code == target_language, installed_languages
        ))
        if len(to_lang) == 0:
            raise ValueError("Target language not installed")
        to_lang = to_lang[0]
        
        translation = from_lang.get_translation(to_lang)
        return translation.translate(text)


    def download_language_package(self, source: str, target: str):
        """
        Downloads and installs the language package for translation.
        """
        available_packages = argostranslate.package.get_available_packages()
        available_package = list(
            filter(
                lambda x: x.from_code == source and x.to_code == target,
                available_packages
            )
        )[0]
        download_path = available_package.download()
        argostranslate.package.install_from_path(download_path)

available_packages = argostranslate.package.get_available_packages()