package coinquant

import (
	"context"
	"fmt"
	"net/http"
)

// ListReports returns paginated report artifacts.
func (c *Client) ListReports(ctx context.Context, opts ReportListOptions) (*ReportsList, error) {
	return doJSON[ReportsList](ctx, c, http.MethodGet, "/v1/reports", nil, opts)
}

// GetReport retrieves a single report.
func (c *Client) GetReport(ctx context.Context, reportID string) (*Report, error) {
	return doJSON[Report](ctx, c, http.MethodGet, fmt.Sprintf("/v1/reports/%s", reportID), nil, nil)
}
