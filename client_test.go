package elevenlabs

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewRequiresAPIKey(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Fatal("expected error when api key is empty")
	}
}

func TestCreateAgentSendsAuthAndBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/convai/agents/create" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("xi-api-key"); got != "test-key" {
			t.Fatalf("missing api key header, got %q", got)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"first_message":"hi"`) {
			t.Fatalf("unexpected body: %s", string(body))
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"agent_id": "agent_123", "name": "MVP"})
	}))
	defer srv.Close()

	client, err := New("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	out, err := client.CreateAgent(context.Background(), CreateAgentInput{
		Name: "MVP",
		ConversationConfig: AgentConversationConfig{
			Agent: &AgentConfig{FirstMessage: "hi"},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.AgentID != "agent_123" {
		t.Fatalf("unexpected agent_id: %s", out.AgentID)
	}
}

func TestAPIErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"detail":"bad key"}`))
	}))
	defer srv.Close()

	client, err := New("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	_, err = client.GetAgent(context.Background(), "agent_123")

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %v", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d", apiErr.StatusCode)
	}
}

func TestAssignAgentSendsPatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/v1/convai/phone-numbers/pn_1" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"agent_id":"agent_1"`) {
			t.Fatalf("unexpected body: %s", string(body))
		}
		_ = json.NewEncoder(w).Encode(PhoneNumber{
			PhoneNumberID: "pn_1",
			AssignedAgent: "agent_1",
		})
	}))
	defer srv.Close()

	client, err := New("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	out, err := client.AssignAgent(context.Background(), "pn_1", "agent_1")
	if err != nil {
		t.Fatal(err)
	}
	if out.AssignedAgent != "agent_1" {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestImportTwilioNumberDefaultsProvider(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), `"provider":"twilio"`) {
			t.Fatalf("expected provider=twilio in body, got: %s", string(body))
		}
		_ = json.NewEncoder(w).Encode(PhoneNumber{PhoneNumberID: "pn_1"})
	}))
	defer srv.Close()

	client, err := New("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.ImportTwilioNumber(context.Background(), ImportTwilioInput{
		PhoneNumber: "+15551234567",
		Label:       "test",
		SID:         "AC...",
		Token:       "...",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestListConversationsBuildsQueryString(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("agent_id"); got != "agent_1" {
			t.Fatalf("expected agent_id=agent_1, got %q", got)
		}
		if got := r.URL.Query().Get("page_size"); got != "5" {
			t.Fatalf("expected page_size=5, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(map[string][]ConversationSummary{
			"conversations": {{ConversationID: "c1"}},
		})
	}))
	defer srv.Close()

	client, err := New("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	out, err := client.ListConversations(context.Background(), ListConversationsOptions{
		AgentID:  "agent_1",
		PageSize: 5,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 1 || out[0].ConversationID != "c1" {
		t.Fatalf("unexpected response: %+v", out)
	}
}

func TestOutboundCallShape(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/convai/twilio/outbound-call" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		var got OutboundCallInput
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		if got.AgentID != "a1" || got.AgentPhoneNumberID != "p1" || got.ToNumber != "+15551234567" {
			t.Fatalf("unexpected body: %+v", got)
		}
		_ = json.NewEncoder(w).Encode(OutboundCallResult{
			Success:        true,
			Message:        "queued",
			ConversationID: "conv_1",
			CallSID:        "CA1",
		})
	}))
	defer srv.Close()

	client, err := New("test-key", WithBaseURL(srv.URL))
	if err != nil {
		t.Fatal(err)
	}
	out, err := client.OutboundCall(context.Background(), OutboundCallInput{
		AgentID:            "a1",
		AgentPhoneNumberID: "p1",
		ToNumber:           "+15551234567",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.ConversationID != "conv_1" || out.CallSID != "CA1" {
		t.Fatalf("unexpected response: %+v", out)
	}
}
