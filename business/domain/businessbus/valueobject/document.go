package valueobject

import (
	"fmt"
	"time"

	"github.com/segmentio/ksuid"
)

type Document struct {
	Type         string
	Name         string
	OriginalName string
	Status       string
	CreatedAt    time.Time
}

func NewDocument(t string, origName string, status string, createdAt time.Time, name *string) (Document, error) {
	if t == "" {
		return Document{}, fmt.Errorf("type is required")
	}
	if origName == "" {
		return Document{}, fmt.Errorf("original name is required")
	}
	if status == "" {
		return Document{}, fmt.Errorf("status name is required")
	}
	if createdAt.IsZero() {
		return Document{}, fmt.Errorf("created at is required")
	}

	doc := Document{
		Name:         ksuid.New().String(),
		Type:         t,
		OriginalName: origName,
		Status:       status,
		CreatedAt:    createdAt,
	}

	if name != nil {
		doc.Name = *name
	}

	return doc, nil
}
