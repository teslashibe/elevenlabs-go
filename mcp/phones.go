package mcp

import (
	"context"

	elevenlabs "github.com/teslashibe/elevenlabs-go"
	"github.com/teslashibe/mcptool"
)

// ImportTwilioNumberInput is the typed input for elevenlabs_import_twilio_number.
type ImportTwilioNumberInput struct {
	PhoneNumber string `json:"phone_number" jsonschema:"description=Twilio number in E.164 format (e.g. '+15551234567'),required"`
	Label       string `json:"label" jsonschema:"description=human-readable label shown in the dashboard,required"`
	SID         string `json:"sid" jsonschema:"description=Twilio Account SID,required"`
	Token       string `json:"token" jsonschema:"description=Twilio Auth Token,required"`
	Provider    string `json:"provider,omitempty" jsonschema:"description=phone provider; defaults to 'twilio'"`
}

func importTwilioNumber(ctx context.Context, c *elevenlabs.Client, in ImportTwilioNumberInput) (any, error) {
	return c.ImportTwilioNumber(ctx, elevenlabs.ImportTwilioInput{
		PhoneNumber: in.PhoneNumber,
		Label:       in.Label,
		SID:         in.SID,
		Token:       in.Token,
		Provider:    in.Provider,
	})
}

// ListPhoneNumbersInput is the typed input for elevenlabs_list_phone_numbers.
type ListPhoneNumbersInput struct{}

func listPhoneNumbers(ctx context.Context, c *elevenlabs.Client, _ ListPhoneNumbersInput) (any, error) {
	res, err := c.ListPhoneNumbers(ctx)
	if err != nil {
		return nil, err
	}
	return mcptool.PageOf(res, "", 0), nil
}

// AssignAgentInput is the typed input for elevenlabs_assign_agent.
type AssignAgentInput struct {
	PhoneNumberID string `json:"phone_number_id" jsonschema:"description=phone_number_id from elevenlabs_list_phone_numbers,required"`
	AgentID       string `json:"agent_id" jsonschema:"description=agent_id to attach to the phone number,required"`
}

func assignAgent(ctx context.Context, c *elevenlabs.Client, in AssignAgentInput) (any, error) {
	return c.AssignAgent(ctx, in.PhoneNumberID, in.AgentID)
}

var phoneTools = []mcptool.Tool{
	mcptool.Define[*elevenlabs.Client, ImportTwilioNumberInput](
		"elevenlabs_import_twilio_number",
		"Register an existing Twilio number with ElevenLabs and return its phone_number_id",
		"ImportTwilioNumber",
		importTwilioNumber,
	),
	mcptool.Define[*elevenlabs.Client, ListPhoneNumbersInput](
		"elevenlabs_list_phone_numbers",
		"List every phone number registered to the ElevenLabs workspace",
		"ListPhoneNumbers",
		listPhoneNumbers,
	),
	mcptool.Define[*elevenlabs.Client, AssignAgentInput](
		"elevenlabs_assign_agent",
		"Attach an agent to a phone number so inbound calls route to it",
		"AssignAgent",
		assignAgent,
	),
}
