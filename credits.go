package coinquant

import (
	"context"
	"fmt"
	"net/http"
)

// GetCredits returns the current credit balance.
func (c *Client) GetCredits(ctx context.Context) (*CreditBalance, error) {
	return doJSON[CreditBalance](ctx, c, http.MethodGet, "/v1/credits", nil, nil)
}

// GetCreditUsage returns aggregated credit usage over time.
func (c *Client) GetCreditUsage(ctx context.Context, opts CreditUsageOptions) (*CreditUsage, error) {
	return doJSON[CreditUsage](ctx, c, http.MethodGet, "/v1/credits/usage", nil, opts)
}

// ListCreditTransactions returns the credit ledger.
func (c *Client) ListCreditTransactions(ctx context.Context, opts CreditTransactionOptions) (*CreditTransactionsList, error) {
	return doJSON[CreditTransactionsList](ctx, c, http.MethodGet, "/v1/credits/transactions", nil, opts)
}

// EstimateTickCredits returns a deterministic credit cost estimate.
func (c *Client) EstimateTickCredits(ctx context.Context, dataDays, timeframeMinutes int) (*TickEstimate, error) {
	if dataDays <= 0 || timeframeMinutes <= 0 {
		return nil, fmt.Errorf("coinquant: data_days and timeframe_minutes must be positive")
	}
	return doJSON[TickEstimate](ctx, c, http.MethodGet, "/v1/credits/estimates/tick", nil, TickEstimateOptions{
		DataDays:         dataDays,
		TimeframeMinutes: timeframeMinutes,
	})
}
