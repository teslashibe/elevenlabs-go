package mcp

import (
	"context"

	elevenlabs "github.com/teslashibe/elevenlabs-go"
	"github.com/teslashibe/mcptool"
)

// AgentPromptInput mirrors elevenlabs.AgentPrompt for MCP input.
type AgentPromptInput struct {
	Prompt string `json:"prompt,omitempty" jsonschema:"description=system prompt for the LLM driving the agent"`
	LLM    string `json:"llm,omitempty" jsonschema:"description=LLM model identifier (e.g. 'gpt-4o-mini')"`
}

// AgentTTSInput mirrors elevenlabs.AgentTTS for MCP input.
type AgentTTSInput struct {
	VoiceID string `json:"voice_id,omitempty" jsonschema:"description=ElevenLabs voice ID used for TTS responses"`
}

// AgentConfigInput mirrors elevenlabs.AgentConfig for MCP input.
type AgentConfigInput struct {
	FirstMessage string            `json:"first_message,omitempty" jsonschema:"description=initial message the agent speaks when a call starts"`
	Language     string            `json:"language,omitempty" jsonschema:"description=ISO language code (e.g. 'en')"`
	Prompt       *AgentPromptInput `json:"prompt,omitempty" jsonschema:"description=LLM-facing prompt and model"`
}

// AgentConversationConfigInput mirrors elevenlabs.AgentConversationConfig.
type AgentConversationConfigInput struct {
	Agent *AgentConfigInput `json:"agent,omitempty" jsonschema:"description=conversational behaviour config"`
	TTS   *AgentTTSInput    `json:"tts,omitempty" jsonschema:"description=text-to-speech config"`
}

func (in AgentConversationConfigInput) toClient() elevenlabs.AgentConversationConfig {
	out := elevenlabs.AgentConversationConfig{}
	if in.Agent != nil {
		ac := &elevenlabs.AgentConfig{
			FirstMessage: in.Agent.FirstMessage,
			Language:     in.Agent.Language,
		}
		if in.Agent.Prompt != nil {
			ac.Prompt = &elevenlabs.AgentPrompt{
				Prompt: in.Agent.Prompt.Prompt,
				LLM:    in.Agent.Prompt.LLM,
			}
		}
		out.Agent = ac
	}
	if in.TTS != nil {
		out.TTS = &elevenlabs.AgentTTS{VoiceID: in.TTS.VoiceID}
	}
	return out
}

// CreateAgentInput is the typed input for elevenlabs_create_agent.
type CreateAgentInput struct {
	Name               string                       `json:"name,omitempty" jsonschema:"description=human-readable agent name"`
	ConversationConfig AgentConversationConfigInput `json:"conversation_config" jsonschema:"description=conversational behaviour, prompt, and TTS voice,required"`
}

func createAgent(ctx context.Context, c *elevenlabs.Client, in CreateAgentInput) (any, error) {
	return c.CreateAgent(ctx, elevenlabs.CreateAgentInput{
		Name:               in.Name,
		ConversationConfig: in.ConversationConfig.toClient(),
	})
}

// GetAgentInput is the typed input for elevenlabs_get_agent.
type GetAgentInput struct {
	AgentID string `json:"agent_id" jsonschema:"description=ElevenLabs agent ID,required"`
}

func getAgent(ctx context.Context, c *elevenlabs.Client, in GetAgentInput) (any, error) {
	return c.GetAgent(ctx, in.AgentID)
}

// UpdateAgentInput is the typed input for elevenlabs_update_agent.
type UpdateAgentInput struct {
	AgentID            string                        `json:"agent_id" jsonschema:"description=ElevenLabs agent ID to update,required"`
	Name               string                        `json:"name,omitempty" jsonschema:"description=new agent name (omit to leave unchanged)"`
	ConversationConfig *AgentConversationConfigInput `json:"conversation_config,omitempty" jsonschema:"description=updated conversational config (omit to leave unchanged)"`
}

func updateAgent(ctx context.Context, c *elevenlabs.Client, in UpdateAgentInput) (any, error) {
	body := elevenlabs.UpdateAgentInput{Name: in.Name}
	if in.ConversationConfig != nil {
		cfg := in.ConversationConfig.toClient()
		body.ConversationConfig = &cfg
	}
	return c.UpdateAgent(ctx, in.AgentID, body)
}

// DeleteAgentInput is the typed input for elevenlabs_delete_agent.
type DeleteAgentInput struct {
	AgentID string `json:"agent_id" jsonschema:"description=ElevenLabs agent ID to delete,required"`
}

func deleteAgent(ctx context.Context, c *elevenlabs.Client, in DeleteAgentInput) (any, error) {
	if err := c.DeleteAgent(ctx, in.AgentID); err != nil {
		return nil, err
	}
	return map[string]any{"ok": true, "agent_id": in.AgentID}, nil
}

var agentTools = []mcptool.Tool{
	mcptool.Define[*elevenlabs.Client, CreateAgentInput](
		"elevenlabs_create_agent",
		"Create an ElevenLabs conversational agent with prompt, language, and TTS voice config",
		"CreateAgent",
		createAgent,
	),
	mcptool.Define[*elevenlabs.Client, GetAgentInput](
		"elevenlabs_get_agent",
		"Fetch an ElevenLabs conversational agent by ID",
		"GetAgent",
		getAgent,
	),
	mcptool.Define[*elevenlabs.Client, UpdateAgentInput](
		"elevenlabs_update_agent",
		"Update an ElevenLabs agent's name or conversation config (partial update)",
		"UpdateAgent",
		updateAgent,
	),
	mcptool.Define[*elevenlabs.Client, DeleteAgentInput](
		"elevenlabs_delete_agent",
		"Permanently delete an ElevenLabs conversational agent by ID",
		"DeleteAgent",
		deleteAgent,
	),
}
