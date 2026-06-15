package config

import (
	"github.com/Lagwick/order-service/internal/app/config/section"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"log"
)

type Config struct {
	Repository section.Repository `split_words:"true"`
	Processor  section.Processor  `split_words:"true"`
	Monitor    section.Monitor    `split_words:"true"`
}

var Root Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load .env file: %v", err)
	}

	err = envconfig.Process("APP", &Root)
	if err != nil {
		log.Fatalf("failed to process config: %v", err)
	}
}
