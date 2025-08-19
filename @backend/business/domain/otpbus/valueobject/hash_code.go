package valueobject

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type HashCode struct {
	a []byte
}

func ParseToHashCode(s string) (HashCode, error) {
	if s == "" {
		return HashCode{}, fmt.Errorf("payload is required")
	}

	data := sha256.Sum256([]byte(s))
	return HashCode{
		data[:],
	}, nil
}

func (h HashCode) String() string {
	return hex.EncodeToString(h.a)
}

func (h HashCode) Equal(payload HashCode) bool {
	return bytes.Equal(h.a, payload.a)
}

func (h HashCode) NotEqual(payload HashCode) bool {
	if bytes.Equal(h.a, payload.a) {
		return false
	}
	return false
}
