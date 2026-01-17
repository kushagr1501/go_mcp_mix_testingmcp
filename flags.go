package main

import (
	"strings"
)

type FILTERTYPE int

const (
	FILTERALL FILTERTYPE = iota
	FILEPDFS
	FILTERImages
	FILTERDocuments
	FilterArchive
)

// struct ,a file configuration holds filter settings
type FilterConfig struct {
	FileType       FILTERTYPE
	IncludePattern string
	ExcludePattern string
	MinSizeMB      int64
	MaxSizeMB      int64
}

func getFileName(path string) string {
	str := strings.Split(path, "\\")
	return str[len(str)-1] //returns the filename with extension
}

func getExtension(filename string) string {
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex == -1 {
		return ""
	}
	return filename[dotIndex:]
}

//parse the extension to the FilterType
func parserFilterType(a string) FILTERTYPE {
	str := strings.ToLower(a)
	switch str {
	case "pdf":
		return FILEPDFS
	case "all":
		return FILTERALL
	case "image", "img", "jpg", "png", "gif", "jpeg":
		return FILTERImages
	case "docs", "doc", "docx", "txt":
		return FILTERDocuments
	case "archive", "zip", "tar", "rar", "7z":
		return FilterArchive
	default:
		return FILTERALL
	}

}

//what files to include during the scan 
func ShouldInclude(path string, config FilterConfig) bool {
	filename := getFileName(path)
	ext := strings.ToLower(getExtension(filename))

	switch config.FileType {

	case FILEPDFS:
		return ext == ".pdf"

	case FILTERImages:
		return ext == ".jpg" || ext == ".jpeg" ||
			ext == ".png" || ext == ".gif"

	case FILTERDocuments:
		return ext == ".pdf" || ext == ".doc" ||
			ext == ".docx" || ext == ".txt"

	case FilterArchive:
		return ext == ".zip" || ext == ".tar" ||
			ext == ".rar" || ext == ".7z"

	case FILTERALL:
		return true

	default:
		return true
	}
}
func ShouldExclude(path string, config FilterConfig) bool {
	if config.ExcludePattern == "" {
		return false
	}

	filename := strings.ToLower(getFileName(path))
	pattern := strings.ToLower(config.ExcludePattern)

	return strings.Contains(filename, pattern)
}

//optional thp ,allow files which is in the defined limit
func ShouldIncludeSize(path string, config FilterConfig, size int64) bool {
	if config.MinSizeMB == 0 && config.MaxSizeMB == 0 {
		return true
	}

	sizeMB := size / (1024 * 1024)

	if config.MinSizeMB > 0 && sizeMB < config.MinSizeMB {
		return false
	}

	if config.MaxSizeMB > 0 && sizeMB > config.MaxSizeMB {
		return false
	}

	return true
}

func (f FILTERTYPE) String() string {
	switch f {
	case FILTERALL:
		return "All Files"
	case FILEPDFS:
		return "PDFs"
	case FILTERImages:
		return "Images"
	case FILTERDocuments:
		return "Documents"
	case FilterArchive:
		return "Archives"
	default:
		return "Unknown"
	}
}
