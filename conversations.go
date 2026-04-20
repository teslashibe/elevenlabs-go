package elevenlabs

import (
	"context"
	"net/url"
	"strconv"
)

// TranscriptTurn is a single line in a conversation transcript.
type TranscriptTurn struct {
	Role           string  `json:"role"`
	Message        string  `json:"message"`
	TimeInCallSecs int     `json:"time_in_call_secs"`
	Feedback       *string `json:"feedback,omitempty"`
}

// Conversation is the minimal subset of the conversation detail response.
type Conversation struct {
	ConversationID string           `json:"conversation_id"`
	AgentID        string           `json:"agent_id"`
	Status         string           `json:"status"`
	Transcript     []TranscriptTurn `json:"transcript"`
}

// ConversationSummary is a row from the list endpoint.
type ConversationSummary struct {
	ConversationID  string `json:"conversation_id"`
	AgentID         string `json:"agent_id"`
	Status          string `json:"status"`
	StartTimeUnix   int64  `json:"start_time_unix_secs"`
	CallDurationSec int    `json:"call_duration_secs"`
}

// ListConversationsOptions filters the list endpoint.
type ListConversationsOptions struct {
	AgentID  string
	PageSize int
}

// GetConversation fetches one conversation by ID.
func (c *Client) GetConversation(ctx context.Context, conversationID string) (Conversation, error) {
	var out Conversation
	if err := c.do(ctx, "GET", "/v1/convai/conversations/"+url.PathEscape(conversationID), nil, &out); err != nil {
		return Conversation{}, err
	}
	return out, nil
}

// ListConversations returns recent conversations, optionally filtered.
func (c *Client) ListConversations(ctx context.Context, opts ListConversationsOptions) ([]ConversationSummary, error) {
	q := url.Values{}
	if opts.AgentID != "" {
		q.Set("agent_id", opts.AgentID)
	}
	if opts.PageSize > 0 {
		q.Set("page_size", strconv.Itoa(opts.PageSize))
	}
	path := "/v1/convai/conversations"
	if encoded := q.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var raw struct {
		Conversations []ConversationSummary `json:"conversations"`
	}
	if err := c.do(ctx, "GET", path, nil, &raw); err != nil {
		return nil, err
	}
	return raw.Conversations, nil
}

// GetSignedURL returns a short-lived signed WebSocket URL for browser/native
// SDKs to talk to an agent directly (no telephony required).
func (c *Client) GetSignedURL(ctx context.Context, agentID string) (string, error) {
	var raw struct {
		SignedURL string `json:"signed_url"`
	}
	path := "/v1/convai/conversation/get-signed-url?agent_id=" + url.QueryEscape(agentID)
	if err := c.do(ctx, "GET", path, nil, &raw); err != nil {
		return "", err
	}
	return raw.SignedURL, nil
}
