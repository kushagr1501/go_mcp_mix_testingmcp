package main

import "time"


func IsLikelyUnused(info *FileInfo, days int) bool {
	if info.IsDirectory {
		return false
	}

	threshold := time.Duration(days) * 24 * time.Hour //covert given days into 24hrs based
	return time.Since(info.AccessedAt) > threshold
}