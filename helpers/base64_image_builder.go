package helpers

import (
	"bytes"
	"fmt"
	"image"
	"strings"
)

type ImageBuilder struct {
	Base64Pict string
}

func NewImageBuilder(pict string) ImageBuilder {
	return ImageBuilder{Base64Pict: pict}
}

func (ib ImageBuilder) ReconstructBase64Image() (bfr bytes.Buffer, format string, err error) {
	dec, err := Base64Decoder(ib.Base64Pict)
	if err != nil {
		return
	}

	_, err = bfr.ReadFrom(dec)

	if err != nil {
		return
	}

	_, format, err = image.Decode(bytes.NewReader(bfr.Bytes()))
	return
}

func (ib ImageBuilder) IsImageTypeValid() (isValid bool) {
	imgType := ib.Base64Pict[(strings.Index(ib.Base64Pict, "/") + 1):(strings.Index(ib.Base64Pict, ";"))]

	fmt.Println(imgType)
	switch imgType {
	case "jpeg":
		isValid = true
	case "png":
		isValid = true
	default:
		isValid = false
	}
	return
}
