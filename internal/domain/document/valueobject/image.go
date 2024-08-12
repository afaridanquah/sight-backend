package valueobject

import (
	"errors"
	"time"
)

type Image struct {
	filename         string
	originalFileName string
	side             DocumentSide
	size             int
	createdAt        time.Time
	updatedAt        time.Time
}

var (
	ErrFileNameIsRequired         = errors.New("filename is required")
	ErrOriginalFileNameIsRequired = errors.New("original filename is required")
	ErrSizeIsRequired             = errors.New("original file size is required")
)

func NewImage(filename string, originalName string, size int, ca time.Time, ua time.Time, side DocumentSide) (Image, error) {
	if filename == "" {
		return Image{}, ErrFileNameIsRequired
	}
	if originalName == "" {
		return Image{}, ErrOriginalFileNameIsRequired
	}
	if size == 0 {
		return Image{}, ErrSizeIsRequired
	}

	return Image{
		filename:         filename,
		originalFileName: originalName,
		side:             side,
		size:             size,
		createdAt:        ca,
		updatedAt:        ua,
	}, nil
}
