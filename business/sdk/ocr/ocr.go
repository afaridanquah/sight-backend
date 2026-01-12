package ocr

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

type Rawo struct {
	Content []byte
}

var (
	ErrOcrFileByteRequired = errors.New("filename is required")
)

func NewByByteOcr(b []byte) (Rawo, error) {
	if len(b) == 0 {
		return Rawo{}, ErrOcrFileByteRequired
	}
	client := gosseract.NewClient()

	defer func() {
		_ = client.Close()
	}()

	if err := client.SetImageFromBytes(b); err != nil {
		return Rawo{}, err
	}

	text, err := client.Text()
	if err != nil {
		return Rawo{}, err
	}

	c, err := json.Marshal(text)
	if err != nil {
		return Rawo{}, err
	}
	fmt.Printf("%v", text)
	return Rawo{
		Content: c,
	}, nil
}
