package main

import (
	"time"
)

type FileInfo struct {
	Path        string
	SizeBytes   int64
	CreatedAt   time.Time
	ModifiedAt  time.Time
	AccessedAt  time.Time
	IsFile      bool
	IsDirectory bool
	MimeType    string
}

//btw MimeType tells what's the extension of a file ,whether it's .pdf,.txt etc
