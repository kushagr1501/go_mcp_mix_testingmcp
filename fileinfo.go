package main

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

func GetFileInfo(client *MCPClient, path string) (*FileInfo, error) {
	raw, err := client.ToolCall("get_file_info", map[string]any{
		"path": path,
	})
	if err != nil {
		return nil, err
	}

	var resp MCPResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}

	info := &FileInfo{
		Path: path,
	}

	for _, item := range resp.Result.Content {
		if item.Type != "text" {
			continue
		}

		lines := strings.Split(item.Text, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)

			switch {
			case strings.HasPrefix(line, "Size:"):
				// "Size: 68 bytes"
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					size, err := strconv.ParseInt(parts[1], 10, 64)
					if err == nil {
						info.SizeBytes = size
					}
				}

			case strings.HasPrefix(line, "Created:"):
				t := strings.TrimPrefix(line, "Created:")
				if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(t)); err == nil {
					info.CreatedAt = parsed
				}

			case strings.HasPrefix(line, "Modified:"):
				t := strings.TrimPrefix(line, "Modified:")
				if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(t)); err == nil {
					info.ModifiedAt = parsed
				}

			case strings.HasPrefix(line, "Accessed:"):
				t := strings.TrimPrefix(line, "Accessed:")
				if parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(t)); err == nil {
					info.AccessedAt = parsed
				}

			case strings.HasPrefix(line, "IsFile:"):
				info.IsFile = strings.Contains(line, "true")

			case strings.HasPrefix(line, "IsDirectory:"):
				info.IsDirectory = strings.Contains(line, "true")

			case strings.HasPrefix(line, "MIME Type:"):
				info.MimeType = strings.TrimSpace(strings.TrimPrefix(line, "MIME Type:"))
			}
		}
		if !info.IsDirectory {
			realSize, err := GetRealFileSize(path)
			if err != nil {
				extensions := []string{".pdf", ".doc", ".docx", ".txt"}
				for _, ext := range extensions {
					if size, err2 := GetRealFileSize(path + ext); err2 == nil {
						realSize = size
						err = nil
						info.Path = path + ext
						break
					}
				}
			}
			if err == nil {
				info.SizeBytes = realSize
			}
		}
	}

	return info, nil
}
