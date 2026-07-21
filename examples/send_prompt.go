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

	req := coinquant.StreamingPromptRequest{
		Message: "Generate a simple BTCUSDT 1h EMA crossover strategy.",
	}

	result, err := client.StreamPrompt(ctx, req, func(ev coinquant.StreamEvent) error {
		fmt.Printf("[%s] %s\n", ev.Event, ev.Text)
		return nil
	})
	if err != nil {
		log.Fatalf("stream failed: %v", err)
	}

	fmt.Printf("Stream type: %s\n", result.Type)
	fmt.Printf("Chat ID: %v\n", result.ChatID)
	fmt.Printf("Strategy version ID: %v\n", result.StrategyVersionID)

	// If the stream produced a schema-only strategy, finalize it.
	if result.Type == coinquant.StreamTypeStrategy && result.StrategyVersionID == nil && result.ChatID != nil {
		strategy, err := client.FinalizeChat(ctx, *result.ChatID, "EMA Crossover", "Schema-only materialization")
		if err != nil {
			log.Fatalf("finalize failed: %v", err)
		}
		fmt.Printf("Materialized strategy %s, latest version %s\n", strategy.ID, strategy.LatestVersion.ID)
	}
}
