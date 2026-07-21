package coinquant

import (
	"context"
	"fmt"
	"net/http"
)

// ListStrategies returns paginated strategy containers.
func (c *Client) ListStrategies(ctx context.Context, opts StrategyListOptions) (*StrategiesList, error) {
	return doJSON[StrategiesList](ctx, c, http.MethodGet, "/v1/strategies", nil, opts)
}

// CreateStrategy materialises a strategy from a chat.
func (c *Client) CreateStrategy(ctx context.Context, req CreateStrategyRequest) (*Strategy, error) {
	return doJSON[Strategy](ctx, c, http.MethodPost, "/v1/strategies", req, nil)
}

// FinalizeChat is a helper that materialises a schema-only stream result into a backtestable strategy.
func (c *Client) FinalizeChat(ctx context.Context, chatID, name, description string) (*Strategy, error) {
	return c.CreateStrategy(ctx, CreateStrategyRequest{
		ChatID:      chatID,
		Name:        name,
		Description: description,
	})
}

// GetStrategy retrieves a single strategy container.
func (c *Client) GetStrategy(ctx context.Context, strategyID string) (*Strategy, error) {
	return doJSON[Strategy](ctx, c, http.MethodGet, fmt.Sprintf("/v1/strategies/%s", strategyID), nil, nil)
}

// UpdateStrategy updates a strategy container.
func (c *Client) UpdateStrategy(ctx context.Context, strategyID string, req UpdateStrategyRequest) (*Strategy, error) {
	return doJSON[Strategy](ctx, c, http.MethodPatch, fmt.Sprintf("/v1/strategies/%s", strategyID), req, nil)
}

// ListStrategyVersions returns paginated version summaries.
func (c *Client) ListStrategyVersions(ctx context.Context, strategyID string, opts PaginatedOptions) (*StrategyVersionsList, error) {
	return doJSON[StrategyVersionsList](ctx, c, http.MethodGet, fmt.Sprintf("/v1/strategies/%s/versions", strategyID), nil, opts)
}

// GetStrategyVersion retrieves a full strategy version snapshot.
func (c *Client) GetStrategyVersion(ctx context.Context, strategyID, versionID string) (*StrategyVersion, error) {
	return doJSON[StrategyVersion](ctx, c, http.MethodGet, fmt.Sprintf("/v1/strategies/%s/versions/%s", strategyID, versionID), nil, nil)
}

// UpdateStrategyVersion updates version-level metadata.
func (c *Client) UpdateStrategyVersion(ctx context.Context, strategyID, versionID string, req UpdateStrategyVersionRequest) (*StrategyVersion, error) {
	return doJSON[StrategyVersion](ctx, c, http.MethodPatch, fmt.Sprintf("/v1/strategies/%s/versions/%s", strategyID, versionID), req, nil)
}
