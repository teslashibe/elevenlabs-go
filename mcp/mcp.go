// Package mcp exposes the elevenlabs-go [elevenlabs.Client] surface as a set
// of MCP (Model Context Protocol) tools that any host application can mount
// on its own MCP server.
//
// All tools wrap exported methods on *elevenlabs.Client. Each tool is
// defined via [mcptool.Define] so the JSON input schema is reflected from
// the typed input struct — no hand-maintained schemas, no drift.
//
// Usage from a host application:
//
//	import (
//	    "github.com/teslashibe/mcptool"
//	    elevenlabs "github.com/teslashibe/elevenlabs-go"
//	    elmcp "github.com/teslashibe/elevenlabs-go/mcp"
//	)
//
//	client, _ := elevenlabs.New("xi-api-key-...")
//	for _, tool := range elmcp.Provider{}.Tools() {
//	    // register tool with your MCP server, passing client as the client
//	    // arg when invoking
//	}
//
// The [Excluded] map documents methods on *Client that are intentionally
// not exposed via MCP, with a one-line reason. The coverage test in
// mcp_test.go fails if a new exported method is added without either being
// wrapped by a tool or appearing in [Excluded].
package mcp

import "github.com/teslashibe/mcptool"

// Provider implements [mcptool.Provider] for elevenlabs-go. The zero value
// is ready to use.
type Provider struct{}

// Platform returns "elevenlabs".
func (Provider) Platform() string { return "elevenlabs" }

// Tools returns every elevenlabs-go MCP tool, in registration order.
func (Provider) Tools() []mcptool.Tool {
	out := make([]mcptool.Tool, 0, len(agentTools)+len(conversationTools)+len(callTools)+len(phoneTools))
	out = append(out, agentTools...)
	out = append(out, conversationTools...)
	out = append(out, callTools...)
	out = append(out, phoneTools...)
	return out
}
