package coinquant

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// StreamEventType describes the high-level classification of an SSE stream.
type StreamEventType string

const (
	StreamTypeError    StreamEventType = "error"
	StreamTypeStrategy StreamEventType = "strategy"
	StreamTypeReport   StreamEventType = "report"
	StreamTypeChat     StreamEventType = "chat"
	StreamTypeUnknown  StreamEventType = "unknown"
)

// StreamEvent is a parsed SSE event from CoinQuant.
type StreamEvent struct {
	Type      StreamEventType `json:"type"`
	Event     string          `json:"event"`
	Data      json.RawMessage `json:"data"`
	RequestID string          `json:"request_id,omitempty"`
	ChatID    *string         `json:"chat_id,omitempty"`
	// Pointers populated depending on event content.
	StrategyID        *string         `json:"strategy_id,omitempty"`
	StrategyVersionID *string         `json:"strategy_version_id,omitempty"`
	ReportID          *string         `json:"report_id,omitempty"`
	Schema            json.RawMessage `json:"schema,omitempty"`
	Text              string          `json:"text,omitempty"`
	ErrorCode         string          `json:"error_code,omitempty"`
	ErrorMessage      string          `json:"error_message,omitempty"`
}

// StreamResult aggregates a complete stream once finished.
type StreamResult struct {
	Type              StreamEventType `json:"type"`
	RequestID         string          `json:"request_id,omitempty"`
	ChatID            *string         `json:"chat_id,omitempty"`
	StrategyID        *string         `json:"strategy_id,omitempty"`
	StrategyVersionID *string         `json:"strategy_version_id,omitempty"`
	ReportID          *string         `json:"report_id,omitempty"`
	Schema            json.RawMessage `json:"schema,omitempty"`
	Text              string          `json:"text"`
	Events            []StreamEvent   `json:"events"`
	Err               error           `json:"-"`
}

// StreamCallback is called for every SSE event.
type StreamCallback func(event StreamEvent) error

// StreamPrompt sends a prompt and streams the result. Pass callback to receive individual events.
func (c *Client) StreamPrompt(ctx context.Context, req StreamingPromptRequest, callback StreamCallback) (*StreamResult, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("coinquant: marshal prompt request: %w", err)
	}
	return c.stream(ctx, "/v1/prompts/stream", body, callback)
}

// StreamChatMessage sends a message to a chat and streams the result. Pass callback to receive individual events.
func (c *Client) StreamChatMessage(ctx context.Context, chatID string, req StreamingChatRequest, callback StreamCallback) (*StreamResult, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("coinquant: marshal chat request: %w", err)
	}
	path := "/v1/chats/" + chatID + "/messages:stream"
	return c.stream(ctx, path, body, callback)
}

func (c *Client) stream(ctx context.Context, path string, body []byte, callback StreamCallback) (*StreamResult, error) {
	u := c.BaseURL + path
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("coinquant: create stream request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("coinquant: stream request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, newAPIError(resp.StatusCode, bodyBytes, resp.Header.Get("X-Request-Id"))
	}

	result := &StreamResult{}
	var requestID string
	if id := resp.Header.Get("X-Request-Id"); id != "" {
		requestID = id
	} else {
		requestID = resp.Header.Get("X-Request-ID")
	}
	result.RequestID = requestID

	reader := bufio.NewReader(resp.Body)
	var currentEvent string
	var currentData bytes.Buffer
	var sawReportEvent bool

	dispatchPending := func() error {
		event, eerr := c.dispatchEvent(currentEvent, currentData.Bytes())
		if eerr != nil {
			return eerr
		}
		if event.Event != "" || currentData.Len() > 0 {
			if event.RequestID == "" {
				event.RequestID = requestID
			}
			result.Events = append(result.Events, event)
			if event.Event == "report_block" || event.Event == "report" {
				sawReportEvent = true
			}
			if callback != nil {
				if cbErr := callback(event); cbErr != nil {
					return cbErr
				}
			}
			c.accumulate(result, event)
		}
		currentEvent = ""
		currentData.Reset()
		return nil
	}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// Flush any pending event that was not terminated by a trailing blank line.
				if currentEvent != "" || currentData.Len() > 0 {
					if dErr := dispatchPending(); dErr != nil {
						result.Err = dErr
						return result, dErr
					}
				}
				break
			}
			result.Err = fmt.Errorf("coinquant: read stream: %w", err)
			return result, result.Err
		}
		line = bytes.TrimSuffix(line, []byte("\n"))
		if bytes.HasSuffix(line, []byte("\r")) {
			line = bytes.TrimSuffix(line, []byte("\r"))
		}
		if len(line) == 0 {
			if dErr := dispatchPending(); dErr != nil {
				result.Err = dErr
				return result, dErr
			}
			continue
		}
		if bytes.HasPrefix(line, []byte("event:")) {
			currentEvent = strings.TrimSpace(string(line[6:]))
			continue
		}
		if bytes.HasPrefix(line, []byte("data:")) {
			// Per the SSE spec, a single leading space after "data:" is stripped;
			// multiple "data:" lines within one event are joined with "\n".
			if currentData.Len() > 0 {
				currentData.WriteByte('\n')
			}
			chunk := line[5:]
			if len(chunk) > 0 && chunk[0] == ' ' {
				chunk = chunk[1:]
			}
			currentData.Write(chunk)
			continue
		}
		if bytes.HasPrefix(line, []byte("id:")) || bytes.HasPrefix(line, []byte(":")) {
			continue
		}
	}
	// Final classification per precedence: error > strategy > report > chat > unknown.
	for _, ev := range result.Events {
		if ev.Type == StreamTypeError {
			result.Type = StreamTypeError
			result.Err = fmt.Errorf("coinquant stream error code=%s message=%s", ev.ErrorCode, ev.ErrorMessage)
			break
		}
	}
	switch {
	case result.Type == StreamTypeError:
		// already set above
	case result.StrategyID != nil || result.StrategyVersionID != nil || result.Schema != nil:
		result.Type = StreamTypeStrategy
	case result.ReportID != nil || sawReportEvent:
		result.Type = StreamTypeReport
	case result.Text != "":
		result.Type = StreamTypeChat
	default:
		result.Type = StreamTypeUnknown
	}
	return result, nil
}

func (c *Client) dispatchEvent(eventName string, data []byte) (StreamEvent, error) {
	ev := StreamEvent{Event: eventName}
	ev.Data = data
	switch eventName {
	case "meta":
		ev.Type = StreamTypeUnknown
		var meta struct {
			RequestID string  `json:"request_id"`
			ChatID    *string `json:"chat_id"`
		}
		_ = json.Unmarshal(data, &meta)
		ev.RequestID = meta.RequestID
		ev.ChatID = meta.ChatID
	case "chunk":
		ev.Type = StreamTypeChat
		ev.Text = string(data)
	case "report_block", "report":
		ev.Type = StreamTypeReport
		var report struct {
			ReportID *string `json:"report_id"`
		}
		_ = json.Unmarshal(data, &report)
		ev.ReportID = report.ReportID
	case "result":
		ev.Type = StreamTypeUnknown
		var result struct {
			ChatID            *string         `json:"chat_id"`
			StrategyID        *string         `json:"strategy_id"`
			StrategyVersionID *string         `json:"strategy_version_id"`
			ReportID          *string         `json:"report_id"`
			Schema            json.RawMessage `json:"schema"`
		}
		_ = json.Unmarshal(data, &result)
		ev.ChatID = result.ChatID
		ev.StrategyID = result.StrategyID
		ev.StrategyVersionID = result.StrategyVersionID
		ev.ReportID = result.ReportID
		ev.Schema = result.Schema
	case "error":
		ev.Type = StreamTypeError
		var apiErr struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		_ = json.Unmarshal(data, &apiErr)
		ev.ErrorCode = apiErr.Code
		ev.ErrorMessage = apiErr.Message
	case "done":
		ev.Type = StreamTypeUnknown
	}
	return ev, nil
}

func (c *Client) accumulate(result *StreamResult, ev StreamEvent) {
	switch ev.Event {
	case "chunk":
		result.Text += ev.Text
	case "result":
		if ev.ChatID != nil {
			result.ChatID = ev.ChatID
		}
		if ev.StrategyID != nil {
			result.StrategyID = ev.StrategyID
		}
		if ev.StrategyVersionID != nil {
			result.StrategyVersionID = ev.StrategyVersionID
		}
		if ev.ReportID != nil {
			result.ReportID = ev.ReportID
		}
		if len(ev.Schema) > 0 && string(ev.Schema) != "null" {
			result.Schema = ev.Schema
		}
	case "report_block", "report":
		if ev.ReportID != nil {
			result.ReportID = ev.ReportID
		}
	case "meta":
		if ev.ChatID != nil && result.ChatID == nil {
			result.ChatID = ev.ChatID
		}
		if ev.RequestID != "" && result.RequestID == "" {
			result.RequestID = ev.RequestID
		}
	}
}
