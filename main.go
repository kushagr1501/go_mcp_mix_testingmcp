// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os/exec"
// )

// // ---- MCP response structs ----

// type MCPResponse struct {
// 	Jsonrpc string    `json:"jsonrpc"`
// 	ID      int       `json:"id"`
// 	Result  MCPResult `json:"result"`
// }

// type MCPResult struct {
// 	Content []MCPContent `json:"content"`
// }

// type MCPContent struct {
// 	Type     string       `json:"type"`
// 	Text     string       `json:"text,omitempty"`
// 	Resource *MCPResource `json:"resource,omitempty"`
// }

// type MCPResource struct {
// 	URI      string `json:"uri"`
// 	MimeType string `json:"mimeType"`
// 	Text     string `json:"text"`
// }

// func main() {
// 	cmd := exec.Command(
// 		"mcp-filesystem-server",
// 		"C:\\Users\\skush\\OneDrive\\Desktop\\testingmcp",
// 	)
// 	//YOU  --->  stdin  --->  OTHER PROGRAM(MCP)
// 	stdin, _ := cmd.StdinPipe() // a pipe ,when u write anything ,other program receives that input

// 	//YOU  <---  stdout  <---  OTHER PROGRAM(MCP)
// 	stdout, _ := cmd.StdoutPipe() // a pipe ,when other prog  write anything ,your  program receives that and can read  it

// 	err := cmd.Start()
// 	if err != nil { //nil means success
// 		panic(err) //Something went very wrong. Stop everything
// 	}

// 	// MCP request: list_directory
// 	request := map[string]any{
// 		"jsonrpc": "2.0",
// 		"id":      1,
// 		"method":  "tools/call",
// 		"params": map[string]any{ //function name and arguments are passed
// 			"name": "list_directory",
// 			"arguments": map[string]any{
// 				"path": "C:\\Users\\skush\\OneDrive\\Desktop\\testingmcp",
// 			},
// 		},
// 	}

// 	json.NewEncoder(stdin).Encode(request)
// 	stdin.Close() //tells the mcp no more request

// 	// ---- READ THE FULL RESPONSE FROM THE MCP SERVER  ----
// 	buf := make([]byte, 1024) //data comes in chunks
// 	out := []byte{}

// 	for {
// 		n, err := stdout.Read(buf)
// 		if n > 0 {
// 			out = append(out, buf[:n]...)
// 		}
// 		if err != nil {
// 			break
// 		}
// 	}

// 	var resp MCPResponse
// 	err1 := json.Unmarshal(out, &resp) //json data from out array and put them in the resp
// 	if err1 != nil {
// 		panic(err1)
// 	}

// 	//print resp
// 	// fmt.Println(resp)

// 	//  CLEAN RESPONSE
// 	for _, item := range resp.Result.Content {
// 		if item.Type == "text" {
// 			fmt.Println("---- TEXT BLOCK ----")
// 			fmt.Println(item.Text)
// 		}
// 	}
// }

package main

func main() {
	client, err := NewMCPClient("C:\\Users\\skush\\OneDrive\\Desktop\\testingmcp")

	if err!=nil {
		panic(err)
	}

	files := []string{
		//extracted paths
	}

	for _, f := range files {
		info, err := GetFileInfo(client, f)
		if err!=nil {
			continue
		}

		//we have to do something with the info variable later 
		

	}
}
