package elevenlabs

import (
	"context"
	"net/url"
)

// AgentPrompt configures the LLM-facing side of an agent.
type AgentPrompt struct {
	Prompt string `json:"prompt,omitempty"`
	LLM    string `json:"llm,omitempty"`
}

// AgentTTS configures the text-to-speech side of an agent.
type AgentTTS struct {
	VoiceID string `json:"voice_id,omitempty"`
}

// AgentConfig is the conversational config for an agent. Only the fields
// needed for the MVP are exposed; configure anything else via the dashboard.
type AgentConfig struct {
	FirstMessage string       `json:"first_message,omitempty"`
	Language     string       `json:"language,omitempty"`
	Prompt       *AgentPrompt `json:"prompt,omitempty"`
}

// AgentConversationConfig groups the per-feature config blocks accepted by
// the create/update agent endpoints.
type AgentConversationConfig struct {
	Agent *AgentConfig `json:"agent,omitempty"`
	TTS   *AgentTTS    `json:"tts,omitempty"`
}

// CreateAgentInput is the body for POST /v1/convai/agents/create.
type CreateAgentInput struct {
	Name               string                  `json:"name,omitempty"`
	ConversationConfig AgentConversationConfig `json:"conversation_config"`
}

// UpdateAgentInput is the body for PATCH /v1/convai/agents/{agent_id}.
type UpdateAgentInput struct {
	Name               string                   `json:"name,omitempty"`
	ConversationConfig *AgentConversationConfig `json:"conversation_config,omitempty"`
}

// Agent is the minimal subset of the agent record we surface.
type Agent struct {
	AgentID            string                  `json:"agent_id"`
	Name               string                  `json:"name"`
	ConversationConfig AgentConversationConfig `json:"conversation_config"`
}

// CreateAgent provisions a new agent and returns its assigned ID.
func (c *Client) CreateAgent(ctx context.Context, in CreateAgentInput) (Agent, error) {
	var out Agent
	if err := c.do(ctx, "POST", "/v1/convai/agents/create", in, &out); err != nil {
		return Agent{}, err
	}
	return out, nil
}

// GetAgent fetches an agent by ID.
func (c *Client) GetAgent(ctx context.Context, agentID string) (Agent, error) {
	var out Agent
	if err := c.do(ctx, "GET", "/v1/convai/agents/"+url.PathEscape(agentID), nil, &out); err != nil {
		return Agent{}, err
	}
	return out, nil
}

// UpdateAgent applies a partial update.
func (c *Client) UpdateAgent(ctx context.Context, agentID string, in UpdateAgentInput) (Agent, error) {
	var out Agent
	if err := c.do(ctx, "PATCH", "/v1/convai/agents/"+url.PathEscape(agentID), in, &out); err != nil {
		return Agent{}, err
	}
	return out, nil
}

// DeleteAgent removes an agent.
func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	return c.do(ctx, "DELETE", "/v1/convai/agents/"+url.PathEscape(agentID), nil, nil)
}
