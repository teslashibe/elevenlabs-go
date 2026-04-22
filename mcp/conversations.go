package mcp

import (
	"context"

	elevenlabs "github.com/teslashibe/elevenlabs-go"
	"github.com/teslashibe/mcptool"
)

// GetConversationInput is the typed input for elevenlabs_get_conversation.
type GetConversationInput struct {
	ConversationID string `json:"conversation_id" jsonschema:"description=ElevenLabs conversation ID,required"`
}

func getConversation(ctx context.Context, c *elevenlabs.Client, in GetConversationInput) (any, error) {
	return c.GetConversation(ctx, in.ConversationID)
}

// ListConversationsInput is the typed input for elevenlabs_list_conversations.
type ListConversationsInput struct {
	AgentID  string `json:"agent_id,omitempty" jsonschema:"description=filter to a specific agent's conversations"`
	PageSize int    `json:"page_size,omitempty" jsonschema:"description=results per page,minimum=1,maximum=100,default=20"`
}

func listConversations(ctx context.Context, c *elevenlabs.Client, in ListConversationsInput) (any, error) {
	res, err := c.ListConversations(ctx, elevenlabs.ListConversationsOptions{
		AgentID:  in.AgentID,
		PageSize: in.PageSize,
	})
	if err != nil {
		return nil, err
	}
	limit := in.PageSize
	if limit <= 0 {
		limit = 20
	}
	return mcptool.PageOf(res, "", limit), nil
}

// GetSignedURLInput is the typed input for elevenlabs_get_signed_url.
type GetSignedURLInput struct {
	AgentID string `json:"agent_id" jsonschema:"description=ElevenLabs agent ID to mint a signed URL for,required"`
}

func getSignedURL(ctx context.Context, c *elevenlabs.Client, in GetSignedURLInput) (any, error) {
	url, err := c.GetSignedURL(ctx, in.AgentID)
	if err != nil {
		return nil, err
	}
	return map[string]any{"signed_url": url, "agent_id": in.AgentID}, nil
}

var conversationTools = []mcptool.Tool{
	mcptool.Define[*elevenlabs.Client, GetConversationInput](
		"elevenlabs_get_conversation",
		"Fetch a single ElevenLabs conversation transcript by ID",
		"GetConversation",
		getConversation,
	),
	mcptool.Define[*elevenlabs.Client, ListConversationsInput](
		"elevenlabs_list_conversations",
		"List recent ElevenLabs conversations, optionally filtered to a specific agent",
		"ListConversations",
		listConversations,
	),
	mcptool.Define[*elevenlabs.Client, GetSignedURLInput](
		"elevenlabs_get_signed_url",
		"Mint a short-lived signed WebSocket URL so a browser/native SDK can talk to an agent directly",
		"GetSignedURL",
		getSignedURL,
	),
}
