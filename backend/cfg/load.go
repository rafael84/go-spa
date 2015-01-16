package cfg

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/i18n"
)

func MustLoad() {
	// .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment settings: %s", err)
	}

	// translations
	i18n.MustLoadTranslationFile("i18n/en-us.all.json")
	i18n.MustLoadTranslationFile("i18n/pt-br.all.json")
}
