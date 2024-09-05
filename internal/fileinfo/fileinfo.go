package fileinfo

import (
	"os"
)

type FileInfo struct {
	Path string
	Info os.FileInfo
}

const (
	BytesInKB = 1024
	BytesInMB = 1024 * 1024
	BytesInGB = 1024 * 1024 * 1024
)
