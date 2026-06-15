package main

import (
	"context"
	"log"

	"github.com/Lagwick/order-service/internal/app/config"
	rhealth "github.com/Lagwick/order-service/internal/app/handler/http/health"
	rprocessor "github.com/Lagwick/order-service/internal/app/processor/http"
	rcpostgres "github.com/Lagwick/order-service/internal/app/repository/conn/postgres"
)

func main() {
	config.Load()

	cfg := config.Root

	conn, err := rcpostgres.NewClient(context.Background(), cfg.Repository.Postgres)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	log.Printf("Postgres connection established")
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close postgres connection: %v", err)
		}
	}()
	healthHandler := rhealth.NewHandler()
	httpProc := rprocessor.NewHTTP(healthHandler, cfg.Processor.WebServer)

	if err := httpProc.Serve(); err != nil {
		log.Printf("serve http: %v", err)
		return
	}
}
