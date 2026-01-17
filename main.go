package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var filterConfig FilterConfig
var Filtertypestr string

func init() {
	flag.StringVar(&filterConfig.ExcludePattern, "exclude", "e", "Exclude files matching pattern")
	flag.StringVar(&filterConfig.IncludePattern, "include", "i", "Include only files matching pattern")
	flag.Int64Var(&filterConfig.MinSizeMB, "min-size", 0, "Minimum file size in MB")
	flag.Int64Var(&filterConfig.MaxSizeMB, "max-size", 0, "Maximum file size in MB")

	flag.StringVar(&Filtertypestr, "filter", "f", "Filter by type (all, pdf, doc, img, code, archive)")

}

func main() {
	flag.Parse()

	filterConfig.FileType = FILTERALL

	filterConfig.FileType = parserFilterType(Filtertypestr)

	// PrintLogo()
	PrintHeader("Filesystem Analyzer v2.0")

	reader := bufio.NewReader(os.Stdin)

	PrintSection("Input")
	fmt.Print("  Enter directory path to analyze: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		PrintError("Failed to read input: " + err.Error())
		return
	}

	desiredpath := strings.TrimSpace(input)

	PrintSection("Filter Settings")
	fmt.Printf("  Filter: %s%s%s\n", ColorYellow+ColorBold, filterConfig.FileType.String(), ColorReset)
	fmt.Printf("  Include: %s%s%s\n", ColorDim, filterConfig.IncludePattern, ColorReset)
	fmt.Printf("  Exclude: %s%s%s\n", ColorDim, filterConfig.ExcludePattern, ColorReset)
	fmt.Printf("  Min Size: %d MB%s\n", filterConfig.MinSizeMB, ColorReset)
	fmt.Printf("  Max Size: %d MB%s\n", filterConfig.MaxSizeMB, ColorReset)

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
	matchedCount := 0

	for _, f := range files {
		info, err := GetFileInfo(client, f)
		if err != nil {
			continue
		}
		// Apply filters
		if !ShouldInclude(f, filterConfig) {
			continue
		}

		if !ShouldIncludeSize(f, filterConfig, info.SizeBytes) {
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
		matchedCount++
	}
	PrintDivider()
	fmt.Printf("%sFiles Matching Filter:%s %d%s\n",
		ColorYellow+ColorBold,
		ColorReset,
		matchedCount,
		ColorReset)
	PrintDivider()

	PrintScanComplete(len(files), unusedCount, zeroByteCount)
}
