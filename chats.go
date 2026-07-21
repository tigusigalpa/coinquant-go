package coinquant

import (
	"context"
	"fmt"
	"net/http"
)

// ListChats returns paginated chat sessions.
func (c *Client) ListChats(ctx context.Context, opts ChatsListOptions) (*ChatsList, error) {
	return doJSON[ChatsList](ctx, c, http.MethodGet, "/v1/chats", nil, opts)
}

// CreateChat creates a new chat session.
func (c *Client) CreateChat(ctx context.Context, req CreateChatRequest) (*Chat, error) {
	return doJSON[Chat](ctx, c, http.MethodPost, "/v1/chats", req, nil)
}

// GetChat retrieves a single chat session.
func (c *Client) GetChat(ctx context.Context, chatID string) (*Chat, error) {
	return doJSON[Chat](ctx, c, http.MethodGet, fmt.Sprintf("/v1/chats/%s", chatID), nil, nil)
}

// UpdateChat updates mutable chat metadata.
func (c *Client) UpdateChat(ctx context.Context, chatID string, req UpdateChatRequest) (*Chat, error) {
	return doJSON[Chat](ctx, c, http.MethodPatch, fmt.Sprintf("/v1/chats/%s", chatID), req, nil)
}

// DeleteChat permanently deletes a chat session.
func (c *Client) DeleteChat(ctx context.Context, chatID string) (*ChatDeleteResponse, error) {
	return doJSON[ChatDeleteResponse](ctx, c, http.MethodDelete, fmt.Sprintf("/v1/chats/%s", chatID), nil, nil)
}

// ListMessages returns paginated messages for a chat.
func (c *Client) ListMessages(ctx context.Context, chatID string, opts PaginatedOptions) (*MessagesList, error) {
	return doJSON[MessagesList](ctx, c, http.MethodGet, fmt.Sprintf("/v1/chats/%s/messages", chatID), nil, opts)
}

// AppendMessage appends a user message without invoking the AI.
func (c *Client) AppendMessage(ctx context.Context, chatID string, req AppendMessageRequest) (*Message, error) {
	return doJSON[Message](ctx, c, http.MethodPost, fmt.Sprintf("/v1/chats/%s/messages", chatID), req, nil)
}
