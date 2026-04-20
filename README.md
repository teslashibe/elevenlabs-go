# elevenlabs-go

A small, dependency-free Go client for the [ElevenLabs Agents](https://elevenlabs.io/docs/eleven-agents/overview) REST API.

Scope is deliberately narrow — just enough to test a voice agent end-to-end over telephony:

- create / get / update / delete an agent
- import a Twilio phone number
- place an outbound call
- fetch a conversation transcript
- mint a signed WebSocket URL for browser / native testing

No streaming, no realtime audio, no knowledge-base management. For everything else, use the dashboard or call the API directly.

## Install

```bash
go get github.com/teslashibe/elevenlabs-go
```

Requires Go 1.25+.

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/teslashibe/elevenlabs-go"
)

func main() {
	client, err := elevenlabs.New("xi-api-key-...")
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	agent, err := client.CreateAgent(ctx, elevenlabs.CreateAgentInput{
		Name: "Support Agent",
		ConversationConfig: elevenlabs.AgentConversationConfig{
			Agent: &elevenlabs.AgentConfig{
				FirstMessage: "Hi! How can I help?",
				Language:     "en",
				Prompt: &elevenlabs.AgentPrompt{
					Prompt: "You are a friendly support agent.",
					LLM:    "gpt-4o-mini",
				},
			},
			TTS: &elevenlabs.AgentTTS{VoiceID: "XB0fDUnXU5powFXDhCwa"},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	phone, err := client.ImportTwilioNumber(ctx, elevenlabs.ImportTwilioInput{
		PhoneNumber: "+15551234567",
		Label:       "MVP test line",
		SID:         "AC...",
		Token:       "...",
	})
	if err != nil {
		log.Fatal(err)
	}

	call, err := client.OutboundCall(ctx, elevenlabs.OutboundCallInput{
		AgentID:            agent.AgentID,
		AgentPhoneNumberID: phone.PhoneNumberID,
		ToNumber:           "+15557654321",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("conversation:", call.ConversationID)
}
```

## Errors

Non-2xx responses return `*elevenlabs.APIError`:

```go
var apiErr *elevenlabs.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode, apiErr.Body)
}
```

## Options

```go
client, _ := elevenlabs.New(
	"xi-api-key-...",
	elevenlabs.WithBaseURL("https://api.eu.residency.elevenlabs.io"),
	elevenlabs.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
)
```

## License

MIT
