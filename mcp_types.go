package main

// MCPResponse represents the TOP-LEVEL response sent by an MCP server.
//
// This is the direct JSON-RPC envelope.
// Every MCP response follows this structure:
//
// {
//   "jsonrpc": "2.0",
//   "id": 1,
//   "result": { ... }
// }
//
type MCPResponse struct {
	// JSON-RPC protocol version 
	Jsonrpc string `json:"jsonrpc"`

	// Request ID 
	ID int `json:"id"`

	//  response
	Result MCPResult `json:"result"`
}

// MCPResult contains the OUTPUT of a tool execution.
//
// MCP tools do NOT usually return simple key-value JSON.
// Instead, they return a list of "content blocks"
// (text, resources, images, etc).
//
// This struct wraps that list.

type MCPResult struct {
	// Content is an ordered list of blocks returned by the tool.
	//
	// Example:
	// - a text explanation
	// - a file reference
	// - metadata
	Content []MCPContent `json:"content"`
}

// MCPContent represents ONE block of content returned by MCP.
//
// MCP is designed for AI + human consumption,
// so content is intentionally flexible.
//
// A content block can be:
// - text   → human-readable explanation
// - resource → structured reference (file, URI, etc)

type MCPContent struct {
	// Type tells us WHAT kind of content this block is.
	//
	// Common values:
	// - "text"
	// - "resource"
	Type string `json:"type"`

	// Text contains human-readable output.
	//
	// Example:
	// "Size: 68 bytes\nModified: 2026-01-09..."
	//
	// Only present when Type == "text"
	Text string `json:"text,omitempty"`

	// Resource contains structured data like file references.
	//
	// Only present when Type == "resource"
	Resource *MCPResource `json:"resource,omitempty"`
}

// MCPResource represents a STRUCTURED reference returned by MCP.
//  it is metadata that MCP chooses to expose safely.

// Resources are useful when:
// - referencing files
// - pointing to URIs
// - attaching machine-readable context
type MCPResource struct {
	// URI points to the resource location.
	//
	// Example:
	// file://C:\Users\...\resume.pdf
	URI string `json:"uri"`

	// MimeType describes the type of the resource.
	//
	// Example:
	// "text/plain"
	// "application/pdf"
	MimeType string `json:"mimeType"`

	// Text is a short human-readable summary of the resource.
	//
	// Example:
	// "File: resume.pdf (68 bytes)"
	Text string `json:"text"`
}
