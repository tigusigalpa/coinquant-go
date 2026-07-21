package coinquant

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ListBacktests returns paginated backtests.
func (c *Client) ListBacktests(ctx context.Context, opts BacktestListOptions) (*BacktestsList, error) {
	return doJSON[BacktestsList](ctx, c, http.MethodGet, "/v1/backtests", nil, opts)
}

// CreateBacktest creates an async backtest for a strategy version.
func (c *Client) CreateBacktest(ctx context.Context, strategyVersionID string) (*Backtest, error) {
	return doJSON[Backtest](ctx, c, http.MethodPost, "/v1/backtests", CreateBacktestRequest{StrategyVersionID: strategyVersionID}, nil)
}

// GetBacktest retrieves a backtest detail.
func (c *Client) GetBacktest(ctx context.Context, backtestID string) (*Backtest, error) {
	return doJSON[Backtest](ctx, c, http.MethodGet, fmt.Sprintf("/v1/backtests/%s", backtestID), nil, nil)
}

// GetBacktestResults retrieves full results for a completed backtest.
func (c *Client) GetBacktestResults(ctx context.Context, backtestID string) (*BacktestResults, error) {
	return doJSON[BacktestResults](ctx, c, http.MethodGet, fmt.Sprintf("/v1/backtests/%s/results", backtestID), nil, nil)
}

// GetBacktestSummaryCSV downloads the summary CSV for a completed backtest.
func (c *Client) GetBacktestSummaryCSV(ctx context.Context, backtestID string) (string, error) {
	return c.getCSV(ctx, fmt.Sprintf("/v1/backtests/%s/exports/summary.csv", backtestID))
}

// GetBacktestTradesCSV downloads the trades CSV for a completed backtest.
func (c *Client) GetBacktestTradesCSV(ctx context.Context, backtestID string) (string, error) {
	return c.getCSV(ctx, fmt.Sprintf("/v1/backtests/%s/exports/trades.csv", backtestID))
}

func (c *Client) getCSV(ctx context.Context, path string) (string, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("coinquant: read csv: %w", err)
	}
	if resp.StatusCode >= 400 {
		return "", newAPIError(resp.StatusCode, b, resp.Header.Get("X-Request-Id"))
	}
	return string(b), nil
}

// CompareBacktests runs a comparison across completed backtests.
func (c *Client) CompareBacktests(ctx context.Context, req CompareBacktestsRequest) (*CompareBacktestsResponse, error) {
	return doJSON[CompareBacktestsResponse](ctx, c, http.MethodPost, "/v1/backtests/compare", req, nil)
}

// DuplicateBacktest clones a completed backtest.
func (c *Client) DuplicateBacktest(ctx context.Context, backtestID string) (*DuplicateBacktestResponse, error) {
	return doJSON[DuplicateBacktestResponse](ctx, c, http.MethodPost, fmt.Sprintf("/v1/backtests/%s/duplicate", backtestID), struct{}{}, nil)
}

// CreateBacktestAndWait creates a backtest and polls until terminal status.
func (c *Client) CreateBacktestAndWait(ctx context.Context, strategyVersionID string, timeoutSeconds, pollIntervalSeconds int) (*BacktestPollResult, error) {
	bt, err := c.CreateBacktest(ctx, strategyVersionID)
	if err != nil {
		return nil, err
	}

	result := &BacktestPollResult{Detail: bt}
	ticker := time.NewTicker(time.Duration(pollIntervalSeconds) * time.Second)
	defer ticker.Stop()
	timeout := time.After(time.Duration(timeoutSeconds) * time.Second)

	for {
		select {
		case <-timeout:
			return result, fmt.Errorf("coinquant: backtest polling timed out after %d seconds", timeoutSeconds)
		case <-ctx.Done():
			return result, ctx.Err()
		case <-ticker.C:
			detail, err := c.GetBacktest(ctx, bt.BacktestID)
			if err != nil {
				return result, err
			}
			result.Detail = detail
			switch detail.Status {
			case "completed", "failed", "cancelled", "error", "timeout":
				if detail.Status == "completed" {
					results, err := c.GetBacktestResults(ctx, bt.BacktestID)
					if err != nil {
						return result, err
					}
					result.Results = results
					if csv, csvErr := c.GetBacktestSummaryCSV(ctx, bt.BacktestID); csvErr == nil {
						result.SummaryCSV = csv
					}
					if csv, csvErr := c.GetBacktestTradesCSV(ctx, bt.BacktestID); csvErr == nil {
						result.TradesCSV = csv
					}
				}
				return result, nil
			}
		}
	}
}
