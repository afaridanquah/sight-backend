package valueobject

import (
	"bytes"
	"image"

	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/segment"
)

func NewThreshHoldImage(b []byte) ([]byte, error) {
	reader := bytes.NewReader(b)
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	grayscale := effect.Grayscale(img)
	th := segment.Threshold(grayscale, 120)

	return th.Pix, nil
}
