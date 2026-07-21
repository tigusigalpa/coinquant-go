//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	coinquant "github.com/tigusigalpa/coinquant-go"
)

func main() {
	token := os.Getenv("COINQUANT_TOKEN")
	if token == "" {
		log.Fatal("COINQUANT_TOKEN is required")
	}

	client := coinquant.NewClient(token)
	ctx := context.Background()

	// Public health check
	health, err := client.Health(ctx)
	if err != nil {
		log.Fatalf("health check failed: %v", err)
	}
	fmt.Printf("API health: %s service=%s\n", health.Status, health.Service)

	// Authenticated credit balance
	credits, err := client.GetCredits(ctx)
	if err != nil {
		log.Fatalf("credits failed: %v", err)
	}
	fmt.Printf("Available credits: %d\n", credits.AvailableCreditsTotal)
}
