package main

import (
	"fmt"
	"time"
	// "strings"
)

// func IsLikelyUnused(info *FileInfo, days int) bool {
// 	if info.IsDirectory {
// 		return false
// 	}

// 	threshold := time.Duration(days) * 24 * time.Hour //covert given days into 24hrs based
// 	return time.Since(info.AccessedAt) > threshold
// }

func ExplainUnused(info *FileInfo, days int) *Explanation {
	if info.IsDirectory {
		return nil
	}

	if info.AccessedAt.IsZero() {
		return nil
	}

	if time.Since(info.ModifiedAt) < time.Duration(days)*24*time.Hour {
		return nil
	}

	return &Explanation{
		Reason: "File apprears unused",
		Evidence: []string{
			fmt.Sprintf("Not accessed in last %d days", days),
			fmt.Sprintf("Last modified: %s", info.ModifiedAt.Format("2006-01-02")),
			fmt.Sprintf("Size: %d bytes", info.SizeBytes),
		},
	}

}

func ExplainZeroByte(info *FileInfo) *Explanation {
	if info.IsDirectory {
		return nil
	}

	// Get the REAL file size (works with OneDrive)
	realSize, err := GetRealFileSize(info.Path)
	if err == nil && realSize > 0 {
		return nil // File actually has content
	}

	// Also skip if it's a cloud placeholder (size is in cloud)
	if IsCloudPlaceholder(info.Path) {
		return nil
	}

	// If we get here, file is genuinely 0 bytes
	if info.SizeBytes != 0 && realSize != 0 {
		return nil
	}

	return &Explanation{
		Reason: "File is empty (0 bytes)",
		Evidence: []string{
			"File size is 0 bytes",
			"Likely placeholder or incomplete file",
		},
	}
}
