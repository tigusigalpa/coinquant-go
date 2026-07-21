# coinquant-go

![CoinQuant Golang SDK](https://i.postimg.cc/SKx6Fdhs/coinquant-golang.jpg)

### A Go client for the [CoinQuant](https://www.coinquant.ai/) Public API

*This is an unofficial community package. It is not maintained or endorsed by CoinQuant.*

Build, backtest, and ship algorithmic trading strategies from your Go code, with AI in the loop.

[![Go Reference](https://pkg.go.dev/badge/github.com/tigusigalpa/coinquant-go.svg)](https://pkg.go.dev/github.com/tigusigalpa/coinquant-go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/tigusigalpa/coinquant-go)](go.mod)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![API Docs](https://img.shields.io/badge/API-docs.coinquant.ai-6f42c1)](https://docs.coinquant.ai/public-api-reference)

[Installation](#installation) · [Quick start](#quick-start) · [Streaming](#streaming--the-heart-of-coinquant) · [Backtesting](#backtesting-without-the-boilerplate) · [API reference](#full-api-reference) · [FAQ](#faq)

---

## Related Projects and Documentation

- **[coinquant-go Wiki](https://github.com/tigusigalpa/coinquant-go/wiki)** - Complete guides, API notes, and practical
  usage examples.
- **[coinquant-php](https://github.com/tigusigalpa/coinquant-php)** - The companion CoinQuant SDK for PHP/Laravel
  projects.

## Why this package?

I built this because I wanted to talk to the CoinQuant API from Go without writing the same HTTP and SSE plumbing every
time. CoinQuant lets you describe a trading idea in plain English and turns it into a backtestable strategy; this client
wraps the public API so you can do that from your Go program.

It is intentionally thin: you talk to the AI, get a strategy, backtest it, and read the metrics, all in idiomatic Go,
with `context.Context` first and strongly typed requests and responses.

## What you get

- **Idiomatic Go** — every call takes `context.Context` first, requests and responses are concrete structs with `json`
  tags, no `map[string]any` guesswork.
- **Full endpoint coverage** — all 37 public endpoints: health, chats, strategies, versions, backtests, reports,
  templates, credits, and community.
- **SSE streaming that actually works** — callback-driven streaming for `POST /v1/prompts/stream` and
  `POST /v1/chats/{chat_id}/messages:stream`, with events classified into `error`, `strategy`, `report`, `chat`, or
  `unknown` using CoinQuant's own precedence rules.
- **Backtests without busywork** — `CreateBacktestAndWait` submits a run, polls until it reaches a terminal state, and
  returns the results plus CSV exports in one call.
- **Schema-only materialization** — when the AI returns a strategy schema but no version ID yet, `FinalizeChat` turns it
  into a real versioned strategy you can backtest.
- **Actionable errors** — a typed `APIError` carries HTTP status, `request_id`, machine-readable `code`, and a human
  message, so you know what went wrong and can quote a request ID to support.

## Installation

```bash
go get github.com/tigusigalpa/coinquant-go
```

Requires Go 1.22 or newer.

## Quick start

Grab a token, create a client, and check your balance in a dozen lines:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    coinquant "github.com/tigusigalpa/coinquant-go"
)

func main() {
    client := coinquant.NewClient(os.Getenv("COINQUANT_TOKEN"))
    ctx := context.Background()

    credits, err := client.GetCredits(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("available:", credits.AvailableCreditsTotal)
}
```

## Authentication

CoinQuant uses a JWT bearer token. To get one:

1. Open the CoinQuant web app.
2. Go to **Settings → Service Accounts → New key**.
3. Copy the one-time `access_token` (you will not see it again).

Service-account tokens are valid for **30 days** and work on every endpoint except the human-profile surfaces (`/v1/me`,
`/v1/billing/*`). When a token expires, generate a new one in the frontend.

```go
client := coinquant.NewClient("YOUR_SECRET_TOKEN")
```

> **Keep secrets out of source control.** Read the token from an environment variable or your secret manager, never
> hard-code it.

## Configuring the client

The client works out of the box, but you can tweak it with functional options:

```go
client := coinquant.NewClient(
    os.Getenv("COINQUANT_TOKEN"),
    coinquant.WithTimeout(60*time.Second),          // default is 30s
    coinquant.WithBaseURL("https://api.coinquant.ai"), // override for staging/mocks
    coinquant.WithHTTPClient(myInstrumentedClient),  // bring your own *http.Client
)
```

`WithHTTPClient` is the escape hatch for retries, tracing, proxies, or custom TLS — pass any `*http.Client` and the
client will use it.

## Streaming — the heart of CoinQuant

The AI engine replies over Server-Sent Events. Instead of making you parse raw event frames, the client reads the
stream, hands each event to your callback as it arrives, and returns a classified `StreamResult` when the stream closes.

```go
result, err := client.StreamPrompt(ctx, coinquant.StreamingPromptRequest{
    Message: "Generate a BTCUSDT 1h EMA crossover strategy.",
}, func(ev coinquant.StreamEvent) error {
    // Called for every event as it streams in — useful for live UIs or logs.
    fmt.Print(ev.Text)
    return nil
})
if err != nil {
    log.Fatal(err)
}

fmt.Println("\nclassified as:", result.Type)
```

Pass `nil` as the callback if you only care about the final aggregated result.

### How responses are classified

CoinQuant can answer a prompt with a chat reply, a research report, or a full strategy. Even a "simple" question can
come back as a report, so the client does not guess from your intent. It classifies from the actual events using this
precedence:

| Priority | Condition                                                                   | `result.Type`        |
|----------|-----------------------------------------------------------------------------|----------------------|
| 1        | Any `error` event, or HTTP status ≥ 400                                     | `StreamTypeError`    |
| 2        | A `result` event carrying `strategy_id`, `strategy_version_id`, or `schema` | `StreamTypeStrategy` |
| 3        | Any `report_block` or `report` event                                        | `StreamTypeReport`   |
| 4        | Only `chunk` text events                                                    | `StreamTypeChat`     |
| 5        | None of the above                                                           | `StreamTypeUnknown`  |

### Continuing a conversation

Attach to an existing chat to keep context, or stream a message into a specific chat:

```go
// Continue a prompt thread
result, _ := client.StreamPrompt(ctx, coinquant.StreamingPromptRequest{
    Message: "Now make the exit tighter.",
    ChatID:  result.ChatID,
}, nil)

// Or stream directly into a chat
result, _ = client.StreamChatMessage(ctx, chatID, coinquant.StreamingChatRequest{
    Content: "Add a 2% stop loss.",
}, nil)
```

### Schema-only strategies

Sometimes the AI returns a strategy schema but no `strategy_version_id` yet — it is a blueprint, not something you can
backtest. Materialize it into a real versioned strategy with one call:

```go
if result.Type == coinquant.StreamTypeStrategy && result.StrategyVersionID == nil && result.ChatID != nil {
    strategy, err := client.FinalizeChat(ctx, *result.ChatID, "EMA Crossover", "My first strategy")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("ready to backtest:", strategy.LatestVersion.ID)
}
```

## Backtesting without the boilerplate

Backtests are asynchronous: you submit a strategy version, then poll until it finishes. `CreateBacktestAndWait` wraps
the whole loop — submit, poll every N seconds, stop on any terminal status (`completed`, `failed`, `cancelled`, `error`,
`timeout`), and, on success, fetch the results plus both CSV exports.

```go
result, err := client.CreateBacktestAndWait(ctx, "STRATEGY_VERSION_ID", 900 /* timeout s */, 5 /* poll s */)
if err != nil {
    log.Fatal(err)
}

fmt.Println("status:", result.Detail.Status)
if result.Results != nil {
    fmt.Printf("Total Return: %v%%\n", result.Results.Metrics["Total Return"])
    fmt.Printf("Sharpe:       %v\n", result.Results.Metrics["Sharpe Ratio"])
    // result.SummaryCSV and result.TradesCSV hold the raw CSV exports.
}
```

Prefer to drive the loop yourself? The building blocks are all public: `CreateBacktest`, `GetBacktest`,
`GetBacktestResults`, `GetBacktestSummaryCSV`, and `GetBacktestTradesCSV`.

## Handling errors

Every non-2xx response becomes a typed `*APIError`. Use `errors.As` to inspect it:

```go
credits, err := client.GetCredits(ctx)
if err != nil {
    var apiErr *coinquant.APIError
    if errors.As(err, &apiErr) {
        log.Printf("status=%d code=%s request_id=%s: %s",
            apiErr.Status, apiErr.Code, apiErr.RequestID, apiErr.Message)
        if apiErr.Status == 402 {
            log.Println("out of credits — top up and retry")
        }
    }
    return err
}
```

The `request_id` is the thing to include when you reach out to CoinQuant support; they can trace your exact call with
it.

## A complete workflow

From idea to metrics, end to end:

```go
// 1. Describe the idea and let the AI draft a strategy.
res, _ := client.StreamPrompt(ctx, coinquant.StreamingPromptRequest{
    Message: "Long BTCUSDT when price crosses above the 200 EMA on 1h, exit on cross below.",
}, nil)

// 2. Materialize it if it came back schema-only.
versionID := ""
if res.StrategyVersionID != nil {
    versionID = *res.StrategyVersionID
} else if res.ChatID != nil {
    s, _ := client.FinalizeChat(ctx, *res.ChatID, "EMA 200 Crossover", "")
    versionID = s.LatestVersion.ID
}

// 3. Backtest and wait for the verdict.
bt, _ := client.CreateBacktestAndWait(ctx, versionID, 900, 5)
fmt.Println("done:", bt.Detail.Status, bt.Results.Metrics)
```

## Full API reference

| Endpoint                                     | Method                                          |
|----------------------------------------------|-------------------------------------------------|
| `GET /health`                                | `client.Health`                                 |
| `GET /v1/chats`                              | `client.ListChats`                              |
| `POST /v1/chats`                             | `client.CreateChat`                             |
| `GET /v1/chats/{id}`                         | `client.GetChat`                                |
| `PATCH /v1/chats/{id}`                       | `client.UpdateChat`                             |
| `DELETE /v1/chats/{id}`                      | `client.DeleteChat`                             |
| `GET /v1/chats/{id}/messages`                | `client.ListMessages`                           |
| `POST /v1/chats/{id}/messages`               | `client.AppendMessage`                          |
| `POST /v1/chats/{id}/messages:stream`        | `client.StreamChatMessage`                      |
| `POST /v1/prompts/stream`                    | `client.StreamPrompt`                           |
| `GET /v1/strategies`                         | `client.ListStrategies`                         |
| `POST /v1/strategies`                        | `client.CreateStrategy` / `client.FinalizeChat` |
| `GET /v1/strategies/{id}`                    | `client.GetStrategy`                            |
| `PATCH /v1/strategies/{id}`                  | `client.UpdateStrategy`                         |
| `GET /v1/strategies/{id}/versions`           | `client.ListStrategyVersions`                   |
| `GET /v1/strategies/{id}/versions/{vid}`     | `client.GetStrategyVersion`                     |
| `PATCH /v1/strategies/{id}/versions/{vid}`   | `client.UpdateStrategyVersion`                  |
| `GET /v1/backtests`                          | `client.ListBacktests`                          |
| `POST /v1/backtests`                         | `client.CreateBacktest`                         |
| `GET /v1/backtests/{id}`                     | `client.GetBacktest`                            |
| `GET /v1/backtests/{id}/results`             | `client.GetBacktestResults`                     |
| `GET /v1/backtests/{id}/exports/summary.csv` | `client.GetBacktestSummaryCSV`                  |
| `GET /v1/backtests/{id}/exports/trades.csv`  | `client.GetBacktestTradesCSV`                   |
| `POST /v1/backtests/compare`                 | `client.CompareBacktests`                       |
| `POST /v1/backtests/{id}/duplicate`          | `client.DuplicateBacktest`                      |
| `GET /v1/reports`                            | `client.ListReports`                            |
| `GET /v1/reports/{id}`                       | `client.GetReport`                              |
| `GET /v1/templates`                          | `client.ListTemplates`                          |
| `GET /v1/templates/{id}`                     | `client.GetTemplate`                            |
| `GET /v1/credits`                            | `client.GetCredits`                             |
| `GET /v1/credits/usage`                      | `client.GetCreditUsage`                         |
| `GET /v1/credits/transactions`               | `client.ListCreditTransactions`                 |
| `GET /v1/credits/estimates/tick`             | `client.EstimateTickCredits`                    |
| `GET /v1/community/leaderboards`             | `client.ListLeaderboards`                       |
| `GET /v1/community/backtests`                | `client.ListCommunityBacktests`                 |
| `GET /v1/community/backtests/{id}`           | `client.GetCommunityBacktest`                   |
| `GET /v1/community/me`                       | `client.GetCommunityMe`                         |
| `GET /v1/community/me/activities`            | `client.ListCommunityActivities`                |

## Runnable examples

Each file in [`examples/`](examples/) is a self-contained program you can run with `go run`:

| File                       | What it shows                                          |
|----------------------------|--------------------------------------------------------|
| `auth_check.go`            | Verify connectivity and read your credit balance       |
| `send_prompt.go`           | Stream a prompt and materialize a schema-only strategy |
| `backtest_poll.go`         | Kick off a backtest and wait for the results           |
| `community_leaderboard.go` | Browse the public leaderboard (no token required)      |

```bash
export COINQUANT_TOKEN=your_token_here
go run examples/auth_check.go
```

## FAQ

**Do I need a token for everything?**
No. The health check and the community endpoints (`leaderboards`, `backtests`) work anonymously. Everything else needs a
service-account token.

**Why did my "simple" prompt come back as a report?**
That is expected — CoinQuant decides the response shape, not the caller. Always branch on `result.Type` rather than on
what you asked for.

**How long do tokens last?**
30 days. When one expires you will get a `401`; mint a fresh key in **Settings → Service Accounts**.

**A backtest keeps returning `409` on results.**
The run is not finished yet. Use `CreateBacktestAndWait`, which polls for you, or keep calling `GetBacktest` until the
status is terminal before fetching results.

**Can I plug in retries or tracing?**
Yes — pass your own `*http.Client` via `WithHTTPClient` and wrap it however you like.

## Contributing

Issues and pull requests are welcome. Please keep the code idiomatic, run `go build ./...` and `go vet ./...` before
submitting, and describe the change clearly.

## License

Released under the MIT License.

## Author

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: [sovletig@gmail.com](mailto:sovletig@gmail.com)

---

<div align="center">

Built for quants who ship. If this package saves you time, consider giving the repo a star.

[CoinQuant](https://www.coinquant.ai/) · [API Reference](https://docs.coinquant.ai/public-api-reference) · [Go package](https://github.com/tigusigalpa/coinquant-go)

</div>
