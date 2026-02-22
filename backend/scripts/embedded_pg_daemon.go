package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
)

func main() {
	dataPath := filepath.Join(os.TempDir(), "pos-embedded-pg-data")
	if err := os.RemoveAll(dataPath); err != nil {
		log.Fatalf("failed to reset embedded postgres data path: %v", err)
	}

	config := embeddedpostgres.DefaultConfig().
		Port(5432).
		Locale("C").
		Encoding("UTF8").
		DataPath(dataPath).
		Username("postgres").
		Password("postgres").
		Database("pos_saas")

	db := embeddedpostgres.NewDatabase(config)
	if err := db.Start(); err != nil {
		log.Fatalf("failed to start embedded postgres: %v", err)
	}
	log.Println("[embedded_pg_daemon] PostgreSQL started on 127.0.0.1:5432")

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	log.Println("[embedded_pg_daemon] stopping PostgreSQL")
	if err := db.Stop(); err != nil {
		log.Fatalf("failed to stop embedded postgres: %v", err)
	}
}
