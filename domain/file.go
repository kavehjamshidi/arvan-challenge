package domain

import "io"

type File struct {
	Data   io.Reader
	Size   int64
	FileID string
	UserID string
}
