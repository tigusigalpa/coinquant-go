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
	// Optional token; community endpoints can be called without auth.
	token := os.Getenv("COINQUANT_TOKEN")
	client := coinquant.NewClient(token)
	ctx := context.Background()

	leaderboards, err := client.ListLeaderboards(ctx, coinquant.LeaderboardOptions{
		Season:  1,
		Limit:   10,
		Include: "entries,current_user",
	})
	if err != nil {
		log.Fatalf("leaderboard failed: %v", err)
	}

	fmt.Printf("Season: %d\n", leaderboards.Season)
	for _, e := range leaderboards.Entries {
		fmt.Printf("#%d %s: %d\n", e.Rank, e.Username, e.TotalPoints)
	}
	if leaderboards.CurrentUser != nil {
		fmt.Printf("You: #%d %s: %d\n", leaderboards.CurrentUser.Rank, leaderboards.CurrentUser.Username, leaderboards.CurrentUser.TotalPoints)
	}
}
