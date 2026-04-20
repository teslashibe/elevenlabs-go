package elevenlabs

import "context"

// OutboundCallInput is the body for POST /v1/convai/twilio/outbound-call.
type OutboundCallInput struct {
	AgentID              string `json:"agent_id"`
	AgentPhoneNumberID   string `json:"agent_phone_number_id"`
	ToNumber             string `json:"to_number"`
	CallRecordingEnabled *bool  `json:"call_recording_enabled,omitempty"`
}

// OutboundCallResult mirrors the TwilioOutboundCallResponse schema.
type OutboundCallResult struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	ConversationID string `json:"conversation_id"`
	CallSID        string `json:"callSid"`
}

// OutboundCall asks ElevenLabs to dial ToNumber via Twilio with the given
// agent. The recipient hears the agent when they pick up.
func (c *Client) OutboundCall(ctx context.Context, in OutboundCallInput) (OutboundCallResult, error) {
	var out OutboundCallResult
	if err := c.do(ctx, "POST", "/v1/convai/twilio/outbound-call", in, &out); err != nil {
		return OutboundCallResult{}, err
	}
	return out, nil
}
