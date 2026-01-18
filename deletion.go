package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"
)

type DeletionRecord struct {
	OrigionalFilePath string    `json:"origional_file_path"`
	FileName          string    `json:"filename"`
	FileSize          int64     `json:"filesize"`
	DeletedAt         time.Time `json:"deleted_at"`
	FileType          string    `json:"file_type`
	// RecycleBinFilePath string    `json:"recyclebin_file_path`
}

// deletionHistory :Manages the history of deleted files
type DeletionHistory struct {
	Records []DeletionRecord
}

// Windows API structures for recycle bin
type _SHFILEOPSTRUCT struct {
	Hwnd                  uintptr
	WFunc                 uint32
	PFrom                 *uint16
	PTto                  *uint16
	FFlags                uint16
	FAnyOperationsAborted bool
	HNameMappings         uintptr
	LpszProgressTitle     *uint16
}

// Windows API constants
const (
	FO_DELETE          = 0x0003
	FOF_ALLOWUNDO      = 0x0040
	FOF_NOCONFIRMATION = 0x0010
	FOF_SILENT         = 0x0004
)

// Windows API DLL imports
var (
	shell32             = syscall.NewLazyDLL("shell32.dll")
	procSHFileOperation = shell32.NewProc("SHFileOperationW")
)

//MovetoRecycleBin moves file to the recycle Bin

// MoveToRecycleBin deletes the given file by delegating the operation to the
// Windows Shell, causing the file to be moved to the Windows Recycle Bin
// instead of being permanently removed.
//
// This function uses the native Windows API (SHFileOperationW) because Go's
// standard library does not provide a way to interact with the Recycle Bin.
// Using the Shell API ensures Explorer-consistent behavior, including:
//   - Support for undo / restore operations
//   - Proper handling of OneDrive and cloud-placeholder files
//   - Correct metadata preservation required by the Recycle Bin
//
// The implementation appears complex because it must:
//   - Convert file paths to UTF-16 (required by Windows APIs)
//   - Use a Windows-defined struct with an exact memory layout
//   - Call into a system DLL via syscall and unsafe.Pointer
//
// This complexity is inherent to the Windows API boundary and is not business
// logic. The function should be treated as platform-specific glue code and
// generally not modified unless the underlying Windows API changes.
//
// Note: This function is Windows-only

func MoveToRecycleBin(filePath string) error {
	// SHFileOperationW requires double-null terminated string
	// UTF16PtrFromString adds ONE null, we need to manually create
	// the double-null terminated buffer
	
	// First, convert to UTF-16
	pathUTF16, err := syscall.UTF16FromString(filePath)
	if err != nil {
		return fmt.Errorf("failed to convert path: %v", err)
	}
	
	// UTF16FromString already adds one null terminator
	// Append another null for double-null termination
	pathUTF16 = append(pathUTF16, 0)

	// Set up file operation structure
	shFileOp := &_SHFILEOPSTRUCT{
		WFunc:  FO_DELETE,
		PFrom:  &pathUTF16[0],
		FFlags: FOF_ALLOWUNDO | FOF_NOCONFIRMATION | FOF_SILENT,
	}

	// Call Windows API
	ret, _, _ := procSHFileOperation.Call(uintptr(unsafe.Pointer(shFileOp)))
	if ret != 0 {
		return fmt.Errorf("SHFileOperation failed with code: %d", ret)
	}

	return nil
}

// The following function GetRecycleBinPath() was an early attempt to locate the Windows Recycle Bin
// by guessing common filesystem paths under the user profile.
// This approach is intentionally commented out because it is NOT reliable on
// modern Windows versions. The Recycle Bin does not have a stable or user-visible
// directory path; it is implemented as a hidden, per-drive, per-user (SID-based)
// system folder (e.g. "<Drive>:\\$Recycle.Bin\\<SID>") with internal metadata.
//
// Locating the Recycle Bin via filesystem inspection is therefore best-effort
// at best and incorrect at worst. It is not required for delete or undo
// functionality.
//
// Instead, this project correctly delegates delete operations to the Windows
// Shell API (SHFileOperationW) with FOF_ALLOWUNDO, allowing Windows itself to
// manage Recycle Bin storage, metadata, and restore behavior.
//
// The commented code is kept for documentation purposes only, to highlight
// a naive approach that was intentionally abandoned in favor of the correct
// OS-level solution.

// // GetRecycleBinPath attempts to find the recycle bin path
// func GetRecycleBinPath() (string, error) {
// 	userProfile := os.Getenv("USERPROFILE")  //C:\Users\<username>

// 	if userProfile == "" {
// 		return "", fmt.Errorf("USERPROFILE environment variable not found")
// 	}

// 	recyclePaths := []string{
// 		filepath.Join(userProfile, "AppData", "Local", "Microsoft", "Windows", "Recycle Bin"),
// 		filepath.Join(userProfile, "Recycle Bin"),
// 	}

// 	for _, path := range recyclePaths {
// 		if _, err := os.Stat(path); err == nil {
// 			return path, nil
// 		}
// 	}

// 	return "", fmt.Errorf("recycle bin path not found")
// }


// DeleteFile safely moves a file to recycle bin and records it
func DeleteFile(fileInfo FileInfo, history *DeletionHistory) error {
	// Check if file exists
	if _, err := os.Stat(fileInfo.Path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", fileInfo.Path)
	}
	

	// Move to recycle bin
	if err := MoveToRecycleBin(fileInfo.Path); err != nil {
		return fmt.Errorf("failed to move file to recycle bin: %v", err)
	}

	// Create deletion record
	record := DeletionRecord{
		OrigionalFilePath: fileInfo.Path,
		// RecycleBinPath: "", // Will try to find this
		FileName:  getFileName(fileInfo.Path),
		FileSize:  fileInfo.SizeBytes,
		DeletedAt: time.Now(),
		FileType:  getFileType(fileInfo.Path),
	}

	// // Try to find recycle bin path
	// if recyclePath, err := GetRecycleBinPath(); err == nil {
	// 	record.RecycleBinPath = recyclePath
	// }

	// Add to history
	history.Records = append(history.Records, record)

	return nil
}



// SaveHistory saves deletion history to JSON file
func SaveHistory(history *DeletionHistory, filePath string) error {
	data, err := json.MarshalIndent(history, "", "  ")  //convert to prettyuJson
	if err != nil {
		return fmt.Errorf("failed to marshal history: %v", err)
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write history file: %v", err)
	}

	return nil
}

// LoadHistory: loads deletion history from JSON file
func LoadHistory(filePath string) (*DeletionHistory, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty history if file doesn't exist
			return &DeletionHistory{Records: []DeletionRecord{}}, nil
		}
		return nil, fmt.Errorf("failed to read history file: %v", err)
	}

	var history DeletionHistory
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("failed to unmarshal history: %v", err)
	}

	return &history, nil
}


// GetHistoryFilePath returns the path to the history file
func GetHistoryFilePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".go-filesystem-deletion-history.json")
}


// ShowHistory displays the deletion history
func ShowHistory(history *DeletionHistory) {
	if len(history.Records) == 0 {
		fmt.Println(ColorGreen + "No files deleted yet." + ColorReset)
		return
	}

	fmt.Printf(ColorCyan + "Deletion History (%d files)" + ColorReset + "\n", len(history.Records))
	fmt.Println(ColorCyan + strings.Repeat("─", 80) + ColorReset)

	for i, record := range history.Records {
		fmt.Printf("%s%d.%s %s\n", ColorYellow, i+1, ColorReset, record.FileName)
		fmt.Printf("   Size: %s\n", formatFileSize(record.FileSize))
		fmt.Printf("   Type: %s\n", record.FileType)
		fmt.Printf("   Deleted: %s\n", record.DeletedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Original: %s\n", record.OrigionalFilePath)
		// if record.RecycleBinPath != "" {
		// 	fmt.Printf("   Recycle Bin: %s\n", record.RecycleBinPath)
		// }
		fmt.Println()
	}
}


// UndoLastDeletion attempts to restore the last deleted file
func UndoLastDeletion(history *DeletionHistory) error {
	if len(history.Records) == 0 {
		return fmt.Errorf("no files to undo")
	}

	lastRecord := history.Records[len(history.Records)-1]
	
	fmt.Printf(ColorYellow+"Attempting to restore: %s"+ColorReset+"\n", lastRecord.FileName)
	
	// Note: Actual restore from recycle bin is complex and requires
	// Windows shell APIs. For now, we'll show where the file was moved
	fmt.Printf("File was moved to recycle bin from: %s\n", lastRecord.OrigionalFilePath)
	fmt.Printf("Check your recycle bin to restore it manually.\n") //for now 
	
	// Remove that from history(make sure that file is removed manually)
	history.Records = history.Records[:len(history.Records)-1]
	
	return nil
}



// ConfirmDeletion asks user for confirmation before deleting
func ConfirmDeletion(fileInfo FileInfo) bool {
	fmt.Printf("\n" + ColorRed + "Delete this file?" + ColorReset + "\n")
	fmt.Printf("Name: %s\n", getFileName(fileInfo.Path))
	fmt.Printf("Size: %s\n", formatFileSize(fileInfo.SizeBytes))
	fmt.Printf("Type: %s\n", getFileType(fileInfo.Path))
	fmt.Printf("Path: %s\n", fileInfo.Path)
	fmt.Printf("\n" + ColorYellow + "Are you sure? (y/N): " + ColorReset)

	var response string
	fmt.Scanln(&response)

	return response == "y" || response == "Y"
}
