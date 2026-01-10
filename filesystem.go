package main

type FILEINFO struct{
	Path string
	Size string 
	IsDirectory string
	MimeType string         
	ModifedTime string
	LastAccessedTime string
}

//btw MimeType tells what's the extension of a file ,whether it's .pdf,.txt etc