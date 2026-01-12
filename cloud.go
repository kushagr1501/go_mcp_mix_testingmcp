//go:build windows

package main

import (
	"syscall"
	// "unsafe"
)

const (
	FILE_ATTRIBUTE_RECALL_ON_OPEN        = 0x00040000
	FILE_ATTRIBUTE_RECALL_ON_DATA_ACCESS = 0x00400000
)

// IsCloudPlaceholder returns true if the file is a OneDrive/cloud
// placeholder that hasn't been downloaded yet
func IsCloudPlaceholder(path string) bool {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false
	}

	attrs, err := syscall.GetFileAttributes(pathPtr)
	if err != nil {
		return false
	}

	return attrs&FILE_ATTRIBUTE_RECALL_ON_OPEN != 0 ||
		attrs&FILE_ATTRIBUTE_RECALL_ON_DATA_ACCESS != 0
}

// GetRealFileSize returns the actual file size using Windows API
// This works correctly for OneDrive files unlike os.Stat()
func GetRealFileSize(path string) (int64, error) {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	// Use FindFirstFile which returns correct size for OneDrive files
	var findData syscall.Win32finddata
	handle, err := syscall.FindFirstFile(pathPtr, &findData)
	if err != nil {
		return 0, err
	}
	syscall.FindClose(handle)

	// Combine high and low parts into int64
	size := int64(findData.FileSizeHigh)<<32 + int64(findData.FileSizeLow)
	return size, nil
}