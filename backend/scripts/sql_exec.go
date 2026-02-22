package main

import (
	"flag"
	"log"
	"os"

	"pos-saas/internal/config"
	"pos-saas/internal/pkg/database"
)

func main() {
	filePath := flag.String("file", "", "path to SQL file to execute")
	flag.Parse()

	if *filePath == "" {
		log.Fatal("missing required -file argument")
	}

	content, err := os.ReadFile(*filePath)
	if err != nil {
		log.Fatalf("failed to read SQL file: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec(string(content)); err != nil {
		log.Fatalf("failed to execute SQL file: %v", err)
	}

	log.Printf("executed SQL file: %s", *filePath)
}
