package main

import (
	"encoding/json"
	"regexp"
)

func list_directory(client *MCPClient, path string) ([]string, error) {
	data, err := client.ToolCall("list_directory", map[string]any{
		"path": path,
	})

	if err != nil {
		return nil, err
	}
	// Parse MCP response
	var resp MCPResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	var files []string

	// URLs - captures everything until closing parenthesis
	re := regexp.MustCompile(`file://([^\)]+)`)

	for _, item := range resp.Result.Content {
		if item.Type == "text" {
			matches := re.FindAllStringSubmatch(item.Text, -1)
			for _, m := range matches {
				files = append(files, m[1])
			}
		}
	}

	return files, nil

}
