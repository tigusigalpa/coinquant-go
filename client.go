package coinquant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.coinquant.ai"
	defaultTimeout = 30 * time.Second
)

// Client is the CoinQuant Public API HTTP client.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a new CoinQuant API client. Use WithHTTPClient and WithBaseURL options as needed.
func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: defaultTimeout,
		},
		Token: token,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithHTTPClient replaces the default HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) { c.HTTPClient = hc }
}

// WithBaseURL overrides the default production base URL.
func WithBaseURL(u string) ClientOption {
	return func(c *Client) { c.BaseURL = strings.TrimRight(u, "/") }
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = d
	}
}

// Health checks the API status without authentication.
func (c *Client) Health(ctx context.Context) (*HealthResponse, error) {
	return doJSONNoAuth[HealthResponse](ctx, c, http.MethodGet, "/health", nil)
}

// doJSON performs a JSON request and decodes the response.
func doJSON[T any](ctx context.Context, c *Client, method, path string, body any, params any) (*T, error) {
	resp, err := c.doRequest(ctx, method, path, body, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("coinquant: read response: %w", err)
	}
	requestID := resp.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = resp.Header.Get("X-Request-ID")
	}
	if resp.StatusCode >= 400 {
		return nil, newAPIError(resp.StatusCode, bodyBytes, requestID)
	}
	var apiResp APIResponse[T]
	if err := json.Unmarshal(bodyBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("coinquant: decode response: %w", err)
	}
	if apiResp.Meta != nil && requestID == "" {
		requestID = apiResp.Meta.RequestID
	}
	return &apiResp.Data, nil
}

// doJSONNoAuth performs a request without an Authorization header.
func doJSONNoAuth[T any](ctx context.Context, c *Client, method, path string, body any) (*T, error) {
	return doJSON[T](ctx, c, method, path, body, nil)
}

// doRaw returns the raw HTTP response body and headers.
func (c *Client) doRaw(ctx context.Context, method, path string, body any, params any) (*http.Response, error) {
	return c.doRequest(ctx, method, path, body, params)
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any, params any) (*http.Response, error) {
	u := c.BaseURL + path
	if params != nil {
		q, err := buildQuery(params)
		if err != nil {
			return nil, err
		}
		if q != "" {
			u = u + "?" + q
		}
	}

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("coinquant: marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("coinquant: create request: %w", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("coinquant: do request: %w", err)
	}
	return resp, nil
}

// buildQuery encodes struct fields using url tag values. Supports strings, ints, booleans and *bool.
func buildQuery(params any) (string, error) {
	values := url.Values{}
	v := reflect.ValueOf(params)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", nil
	}
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("url")
		if tag == "" || tag == "-" {
			continue
		}
		parts := strings.Split(tag, ",")
		name := parts[0]
		omitempty := len(parts) > 1 && parts[1] == "omitempty"

		fv := v.Field(i)
		if omitempty && isZero(fv) {
			continue
		}

		switch fv.Kind() {
		case reflect.String:
			if fv.String() != "" {
				values.Set(name, fv.String())
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fv.Int() != 0 {
				values.Set(name, strconv.FormatInt(fv.Int(), 10))
			}
		case reflect.Bool:
			values.Set(name, strconv.FormatBool(fv.Bool()))
		case reflect.Ptr:
			if !fv.IsNil() && fv.Elem().Kind() == reflect.Bool {
				values.Set(name, strconv.FormatBool(fv.Elem().Bool()))
			}
		}
	}
	return values.Encode(), nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Ptr:
		return v.IsNil()
	case reflect.Slice, reflect.Map:
		return v.Len() == 0
	}
	return false
}
