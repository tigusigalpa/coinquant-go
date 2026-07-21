package coinquant

import (
	"context"
	"fmt"
	"net/http"
)

// ListTemplates returns reusable strategy templates.
func (c *Client) ListTemplates(ctx context.Context, opts TemplateListOptions) (*TemplatesList, error) {
	return doJSON[TemplatesList](ctx, c, http.MethodGet, "/v1/templates", nil, opts)
}

// GetTemplate retrieves a full template definition.
func (c *Client) GetTemplate(ctx context.Context, templateID string) (*TemplateDetail, error) {
	return doJSON[TemplateDetail](ctx, c, http.MethodGet, fmt.Sprintf("/v1/templates/%s", templateID), nil, nil)
}
