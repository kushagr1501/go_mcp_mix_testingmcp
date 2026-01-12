package main

import (
	"fmt"
	"time"
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
	if time.Since(info.AccessedAt) < time.Duration(days)*24*time.Hour {
		return nil
	}

	return &Explanation{
		Reason: "File apprears unused",
		Evidence: []string{
			fmt.Sprintf("Not accessed in last %d days", days),
			fmt.Sprintf("Last accessed: %s", info.AccessedAt.Format("2006-01-02")),
			fmt.Sprintf("Size: %d bytes", info.SizeBytes),
		},
	}

}

func ExplainZeroByte(info *FileInfo) *Explanation {
	if info.IsDirectory {
		return nil
	}

	if info.SizeBytes != 0 {
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
