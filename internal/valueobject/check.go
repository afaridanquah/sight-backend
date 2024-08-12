package valueobject

import (
	"time"

	"github.com/google/uuid"
)

const (
	_ Kind = iota
	StandardAML
	ExtensiveAML
	DocumentCheck
)

type Kind int8

type Check struct {
	results   any
	score     int
	kind      Kind
	createdBy uuid.UUID
	createdAt time.Time
}
