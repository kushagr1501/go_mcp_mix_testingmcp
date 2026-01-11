package main

import (
	"encoding/json"
	"io"
	"os/exec"
	"bufio"
)

type MCPClient struct {
	cmd    *exec.Cmd
	stdin  any
	stdout any
	nextID int
}

func NewMCPClient(allowedpath string) (*MCPClient, error) {
	cmd := exec.Command("mcp-filesystem-server", allowedpath)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	return &MCPClient{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		nextID: 1,
	}, nil
}

//NOTE:ToolCall is a function or method of struct MCPClient

func (c *MCPClient) ToolCall(name string, params map[string]any) ([]byte, error) {
	request := map[string]any{
		"jsonrpc": "2.0",
		"id":      c.nextID,
		"method":  "tools/call",
		"params": map[string]any{ //function name and arguments are passed
			"name":      name,
			"arguments": params,
		},
	}
	c.nextID++

	// Fixed : Add type assertion for io.Writer
	json.NewEncoder(c.stdin.(io.Writer)).Encode(request)
	// Fixed : Add type assertion for io.WriteCloser

	// c.stdin.(io.WriteCloser).Close() //tells the mcp no more request

	//  READ THE FULL RESPONSE FROM THE MCP SERVER
	// buf := make([]byte, 1024) //data comes in chunks
	out := []byte{}

	// for {
	// 	n, err := c.stdout.(io.Reader).Read(buf)
	// 	if n > 0 {
	// 		out = append(out, buf[:n]...)
	// 	}
	// 	if err != nil {
	// 		break
	// 	}
	// }
	scanner := bufio.NewScanner(c.stdout.(io.Reader))
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	if scanner.Scan() {
		out = append(out, scanner.Bytes()...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil

}
