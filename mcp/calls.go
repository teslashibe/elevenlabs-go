package mcp

import (
	"context"

	elevenlabs "github.com/teslashibe/elevenlabs-go"
	"github.com/teslashibe/mcptool"
)

// OutboundCallInput is the typed input for elevenlabs_outbound_call.
type OutboundCallInput struct {
	AgentID              string `json:"agent_id" jsonschema:"description=ElevenLabs agent ID that should handle the call,required"`
	AgentPhoneNumberID   string `json:"agent_phone_number_id" jsonschema:"description=phone_number_id from elevenlabs_list_phone_numbers / elevenlabs_import_twilio_number,required"`
	ToNumber             string `json:"to_number" jsonschema:"description=destination phone number in E.164 format (e.g. '+15557654321'),required"`
	CallRecordingEnabled *bool  `json:"call_recording_enabled,omitempty" jsonschema:"description=opt in/out of call recording (default: provider default)"`
}

func outboundCall(ctx context.Context, c *elevenlabs.Client, in OutboundCallInput) (any, error) {
	return c.OutboundCall(ctx, elevenlabs.OutboundCallInput{
		AgentID:              in.AgentID,
		AgentPhoneNumberID:   in.AgentPhoneNumberID,
		ToNumber:             in.ToNumber,
		CallRecordingEnabled: in.CallRecordingEnabled,
	})
}

var callTools = []mcptool.Tool{
	mcptool.Define[*elevenlabs.Client, OutboundCallInput](
		"elevenlabs_outbound_call",
		"Place an outbound Twilio call: ElevenLabs dials to_number and the recipient hears the named agent",
		"OutboundCall",
		outboundCall,
	),
}
