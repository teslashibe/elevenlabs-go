package elevenlabs

import (
	"context"
	"net/url"
)

// ImportTwilioInput is the body for POST /v1/convai/phone-numbers (Twilio).
type ImportTwilioInput struct {
	PhoneNumber string `json:"phone_number"`
	Label       string `json:"label"`
	SID         string `json:"sid"`
	Token       string `json:"token"`
	Provider    string `json:"provider,omitempty"` // defaults to "twilio"
}

// PhoneNumber is the minimal subset of an imported phone number record.
type PhoneNumber struct {
	PhoneNumberID string `json:"phone_number_id"`
	PhoneNumber   string `json:"phone_number"`
	Label         string `json:"label"`
	Provider      string `json:"provider"`
	AssignedAgent string `json:"assigned_agent,omitempty"`
}

// ImportTwilioNumber registers an existing Twilio number with ElevenLabs and
// returns its assigned phone_number_id (used for outbound calls).
func (c *Client) ImportTwilioNumber(ctx context.Context, in ImportTwilioInput) (PhoneNumber, error) {
	if in.Provider == "" {
		in.Provider = "twilio"
	}
	var out PhoneNumber
	if err := c.do(ctx, "POST", "/v1/convai/phone-numbers", in, &out); err != nil {
		return PhoneNumber{}, err
	}
	return out, nil
}

// ListPhoneNumbers returns all phone numbers registered to the workspace.
func (c *Client) ListPhoneNumbers(ctx context.Context) ([]PhoneNumber, error) {
	var out []PhoneNumber
	if err := c.do(ctx, "GET", "/v1/convai/phone-numbers", nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// AssignAgent attaches an agent to a phone number. This is what enables
// inbound calls — once assigned, ElevenLabs auto-configures Twilio so dials
// to that number route to the agent.
func (c *Client) AssignAgent(ctx context.Context, phoneNumberID, agentID string) (PhoneNumber, error) {
	var out PhoneNumber
	body := map[string]string{"agent_id": agentID}
	path := "/v1/convai/phone-numbers/" + url.PathEscape(phoneNumberID)
	if err := c.do(ctx, "PATCH", path, body, &out); err != nil {
		return PhoneNumber{}, err
	}
	return out, nil
}
