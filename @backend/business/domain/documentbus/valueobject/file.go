package valueobject

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/segmentio/ksuid"
)

type File struct {
	//Generated filename to ensure uniques
	Name string
	// Original filename by the uploader
	FileName string
	Size     int64
	Data     multipart.File
}

var allowedExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".pdf":  true,
}

const maxFileSize = 2 * 1024 * 1024 // 2 MB

func NewFile(n string, size int64, data multipart.File) (File, error) {
	ext := filepath.Ext(n)
	if !allowedExtensions[ext] {
		return File{}, fmt.Errorf("extension not allowed %s", ext)
	}

	if size > maxFileSize {
		return File{}, fmt.Errorf("file too large")
	}

	id := ksuid.New().String()
	name := strings.Join([]string{id, ext}, "")

	return File{
		FileName: name,
		Name:     name,
		Size:     size,
		Data:     data,
	}, nil
}
