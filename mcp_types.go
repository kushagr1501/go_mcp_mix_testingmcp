package main

// MCPResponse is the  JSON-RPC response
type MCPResponse struct {
	Jsonrpc string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  MCPResult `json:"result"`
}

// MCPResult wraps content returned by MCP tools
type MCPResult struct {
	Content []MCPContent `json:"content"`
}

// MCPContent represents each content block
type MCPContent struct {
	Type     string       `json:"type"`
	Text     string       `json:"text,omitempty"`
	Resource *MCPResource `json:"resource,omitempty"`
}

// MCPResource represents resource-type content
type MCPResource struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text"`
}
