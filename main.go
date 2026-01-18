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

var deleteMode bool
var undoMode bool
var historyMode bool

func init() {
	flag.StringVar(&filterConfig.ExcludePattern, "exclude", "e", "Exclude files matching pattern")
	flag.StringVar(&filterConfig.IncludePattern, "include", "i", "Include only files matching pattern")
	flag.Int64Var(&filterConfig.MinSizeMB, "min-size", 0, "Minimum file size in MB")
	flag.Int64Var(&filterConfig.MaxSizeMB, "max-size", 0, "Maximum file size in MB")

	flag.StringVar(&Filtertypestr, "filter", "f", "Filter by type (all, pdf, doc, img, code, archive)")

	// Deletion flags
	flag.BoolVar(&deleteMode, "delete", false, "Enable safe file deletion mode")
	flag.BoolVar(&undoMode, "undo", false, "Undo last file deletion")
	flag.BoolVar(&historyMode, "history", false, "Show deletion history")

}

func main() {
	flag.Parse()

	if historyMode {
		handleHistoryMode()
		return
	}

	if undoMode {
		handleUndoMode()
		return
	}

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
	if deleteMode {
		handleDeleteMode(files, client)
		return
	}

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

func handleHistoryMode() {
	PrintHeader("Deletion History")

	history, err := LoadHistory(GetHistoryFilePath())
	if err != nil {
		PrintError("Failed to load history: " + err.Error())
		return
	}

	ShowHistory(history)
}

func handleUndoMode() {
	PrintHeader("Undo Last Deletion")

	history, err := LoadHistory(GetHistoryFilePath())
	if err != nil {
		PrintError("Failed to load history: " + err.Error())
		return
	}

	if err := UndoLastDeletion(history); err != nil {
		PrintError("Failed to undo: " + err.Error())
		return
	}

	// Save updated history
	if err := SaveHistory(history, GetHistoryFilePath()); err != nil {
		PrintError("Failed to save history: " + err.Error())
		return
	}

	PrintSuccess("Undo completed successfully!")
}

func handleDeleteMode(files []string, client *MCPClient) {
	PrintHeader("Safe File Deletion Mode")
	PrintWarning("This will move files to the Recycle Bin - you can restore them later!")

	// Load existing history
	history, err := LoadHistory(GetHistoryFilePath())
	if err != nil {
		PrintError("Failed to load history: " + err.Error())
		return
	}

	deletedCount := 0
	skippedCount := 0

	for _, f := range files {
		// Apply filters
		if !ShouldInclude(f, filterConfig) {
			continue
		}

		info, err := GetFileInfo(client, f)
		if err != nil {
			continue
		}

		// Check size filter
		if !ShouldIncludeSize(f, filterConfig, info.SizeBytes) {
			continue
		}

		// Show file info and ask for confirmation
		fmt.Printf("\n"+ColorCyan+"File %d:"+ColorReset+"\n", deletedCount+skippedCount+1)
		fmt.Printf("Name: %s\n", getFileName(info.Path))
		fmt.Printf("Size: %s\n", formatFileSize(info.SizeBytes))
		fmt.Printf("Type: %s\n", getFileType(info.Path))
		fmt.Printf("Path: %s\n", info.Path)

		// Check if it's unused or zero-byte
		if exp := ExplainUnused(info, 60); exp != nil {
			fmt.Printf(ColorRed+"⚠ Unused: %s"+ColorReset+"\n", exp.Evidence)
		}

		if expzero := ExplainZeroByte(info); expzero != nil {
			fmt.Printf(ColorRed+"⚠ Zero-byte: %s"+ColorReset+"\n", expzero.Reason)
		}

		// Ask for confirmation
		if ConfirmDeletion(*info) {
			if err := DeleteFile(*info, history); err != nil {
				PrintError("Failed to delete file: " + err.Error())
				skippedCount++
				continue
			}
			PrintSuccess("File moved to Recycle Bin!")
			deletedCount++
		} else {
			PrintInfo("Skipped")
			skippedCount++
		}

		// Ask if user wants to continue
		if deletedCount+skippedCount > 0 {
			fmt.Printf("\n" + ColorYellow + "Continue? (Y/n): " + ColorReset)
			var response string
			fmt.Scanln(&response)
			if response == "n" || response == "N" {
				break
			}
		}
	}

	// Save history
	if deletedCount > 0 {
		if err := SaveHistory(history, GetHistoryFilePath()); err != nil {
			PrintError("Failed to save history: " + err.Error())
		} else {
			PrintSuccess(fmt.Sprintf("History saved! %d files deleted.", deletedCount))
		}
	}

	// Show summary
	PrintDivider()
	fmt.Printf(ColorGreen + "Deletion Complete!" + ColorReset + "\n")
	fmt.Printf("Files deleted: %s%d%s\n", ColorYellow+ColorBold, deletedCount, ColorReset)
	fmt.Printf("Files skipped: %s%d%s\n", ColorDim, skippedCount, ColorReset)
	PrintDivider()

	PrintInfo("Use --history to see deleted files")
	PrintInfo("Use --undo to restore the last deleted file")
}
