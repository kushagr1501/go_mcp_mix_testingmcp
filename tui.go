package main

import (
	"fmt"
	"strings"
)

// Colors for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

// func PrintLogo() {
//     // Replaced with a proper sitting cat ASCII
//     logo := `
//         /\_____/\
//        /  o   o  \
//        ( ==  ^  ==  )
//        )         (
//       (           )
//      (__(__)___(__)__
//     `
//     fmt.Printf("%s%s%s%s\n", ColorCyan+ColorBold, logo, ColorReset)
// }

func PrintHeader(title string) {
	width := 60
	padding := (width - len(title)) / 2
	fmt.Printf("\n%s%s%s\n",
		ColorCyan+ColorBold,
		strings.Repeat("═", width),
		ColorReset)
	fmt.Printf("%s%s%s%s%s\n",
		ColorCyan+ColorBold,
		strings.Repeat(" ", padding),
		title,
		strings.Repeat(" ", padding),
		ColorReset)
	fmt.Printf("%s%s%s\n\n",
		ColorCyan+ColorBold,
		strings.Repeat("═", width),
		ColorReset)
}

func PrintSection(title string) {
	fmt.Printf("\n%s%s● %s%s\n",
		ColorYellow+ColorBold,
		strings.Repeat(" ", 2),
		title,
		ColorReset)
}

func PrintFileInfo(label, value string) {
	fmt.Printf("  %s%s%s: %s%s\n",
		ColorBlue,
		strings.Repeat(" ", 4),
		label,
		ColorReset,
		value)
}

func PrintSuccess(message string) {
	fmt.Printf("%s%s✓ %s%s\n",
		ColorGreen,
		strings.Repeat(" ", 4),
		message,
		ColorReset)
}

func PrintWarning(message string) {
	fmt.Printf("%s%s⚠ %s%s\n",
		ColorYellow,
		strings.Repeat(" ", 4),
		message,
		ColorReset)
}

func PrintError(message string) {
	fmt.Printf("%s%s✗ %s%s\n",
		ColorRed,
		strings.Repeat(" ", 4),
		message,
		ColorReset)
}

func PrintFileCount(count int) {
	fmt.Printf("\n%s%s📁 Total Files Scanned: %d%s\n\n",
		ColorBold+ColorCyan,
		strings.Repeat(" ", 2),
		count,
		ColorReset)
}

func PrintUnusedFile(path string, evidence []string) {
	fmt.Printf("\n%s[UNUSED]%s %s\n",
		ColorRed+ColorBold,
		ColorReset,
		ColorBold+path)
	for _, e := range evidence {
		fmt.Printf("  %s%s▸%s %s\n",
			ColorRed,
			strings.Repeat(" ", 2),
			ColorReset,
			e)
	}
}

func PrintZeroByteFile(path, reason string, evidence []string) {
	fmt.Printf("\n%s[ZERO-BYTE]%s %s\n",
		ColorYellow+ColorBold,
		ColorReset,
		ColorBold+path)
	fmt.Printf("  %sReason: %s%s\n",
		ColorYellow,
		ColorReset,
		reason)
	for _, e := range evidence {
		fmt.Printf("  %s%s▸%s %s\n",
			ColorYellow,
			strings.Repeat(" ", 2),
			ColorReset,
			e)
	}
}

func PrintDivider() {
	fmt.Printf("%s%s%s\n",
		ColorCyan,
		strings.Repeat("─", 60),
		ColorReset)
}

func PrintScanComplete(totalFiles, unusedCount, zeroByteCount int) {
	fmt.Printf("\n")
	PrintDivider()
	fmt.Printf("%sScan Summary:%s\n",
		ColorBold,
		ColorReset)
	PrintFileInfo("Files Scanned", fmt.Sprintf("%d", totalFiles))
	PrintFileInfo("Unused Files", fmt.Sprintf("%d", unusedCount))
	PrintFileInfo("Zero-Byte Files", fmt.Sprintf("%d", zeroByteCount))
	PrintDivider()
	fmt.Printf("\n")
}

func PrintInfo(message string) {
	fmt.Printf("%s%sℹ %s%s\n",
		ColorCyan,
		strings.Repeat(" ", 4),
		message,
		ColorReset)
}

func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

func getFileType(path string) string {
	if strings.Contains(strings.ToLower(path), ".pdf") {
		return "PDF"
	}
	if strings.Contains(strings.ToLower(path), ".doc") {
		return "Document"
	}
	if strings.Contains(strings.ToLower(path), ".txt") {
		return "Text"
	}
	if strings.Contains(strings.ToLower(path), ".jpg") || strings.Contains(strings.ToLower(path), ".png") {
		return "Image"
	}
	if strings.Contains(strings.ToLower(path), ".mp4") || strings.Contains(strings.ToLower(path), ".avi") {
		return "Video"
	}
	return "Unknown"
}
