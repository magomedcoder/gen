package domain

import "time"

type File struct {
	Id          int64
	Filename    string
	MimeType    string
	Size        int64
	StoragePath string
	CreatedAt   time.Time
}

func NewFile(filename, mimeType string, size int64, storagePath string) *File {
	return &File{
		Filename:    filename,
		MimeType:    mimeType,
		Size:        size,
		StoragePath: storagePath,
		CreatedAt:   time.Now(),
	}
}
