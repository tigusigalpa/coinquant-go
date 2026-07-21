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
	strategyVersionID := os.Getenv("COINQUANT_STRATEGY_VERSION_ID")
	if token == "" || strategyVersionID == "" {
		log.Fatal("COINQUANT_TOKEN and COINQUANT_STRATEGY_VERSION_ID are required")
	}

	client := coinquant.NewClient(token)
	ctx := context.Background()

	result, err := client.CreateBacktestAndWait(ctx, strategyVersionID, 900, 5)
	if err != nil {
		log.Fatalf("backtest failed: %v", err)
	}

	fmt.Printf("Backtest status: %s\n", result.Detail.Status)
	if result.Results != nil {
		fmt.Printf("Metrics: %+v\n", result.Results.Metrics)
		fmt.Printf("Summary CSV rows: %d\n", len(result.SummaryCSV))
		fmt.Printf("Trades CSV rows: %d\n", len(result.TradesCSV))
	}
}
