package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// PrintLogo()
	PrintHeader("Filesystem Analyzer v1.0")

	reader := bufio.NewReader(os.Stdin)

	PrintSection("Input")
	fmt.Print("  Enter directory path to analyze: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		PrintError("Failed to read input: " + err.Error())
		return
	}

	desiredpath := strings.TrimSpace(input)

	PrintSection("Connecting to MCP Server")
	PrintSuccess("Starting mcp-filesystem-server...")

	client, err := NewMCPClient(desiredpath)
	if err != nil {
		PrintError("Failed to connect: " + err.Error())
		return
	}
	PrintSuccess("Connected!")

	PrintSection("Scanning Directory")
	PrintFileInfo("Path", desiredpath)

	files, err := list_directory(client, desiredpath)
	if err != nil {
		PrintError("Failed to scan: " + err.Error())
		return
	}

	PrintFileCount(len(files))

	unusedCount := 0
	zeroByteCount := 0

	for _, f := range files {
		info, err := GetFileInfo(client, f)
		if err != nil {
			continue
		}

		if exp := ExplainUnused(info, 60); exp != nil {
			PrintUnusedFile(info.Path, exp.Evidence)
			unusedCount++
		}

		if expzero := ExplainZeroByte(info); expzero != nil {
			PrintZeroByteFile(info.Path, expzero.Reason, expzero.Evidence)
			zeroByteCount++
		}
	}

	PrintScanComplete(len(files), unusedCount, zeroByteCount)
}
