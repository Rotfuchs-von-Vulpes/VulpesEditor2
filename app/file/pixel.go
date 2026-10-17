package file

import (
	"image"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

func LoadPixelImage(r io.Reader, coeff int32) (finalImg *image.RGBA, err error) {
	if coeff <= 0 {
		coeff = 1
	}
	img, err := png.Decode(r)
	if err != nil {
		return
	}
	width := coeff * int32(img.Bounds().Dx())
	height := coeff * int32(img.Bounds().Dy())
	finalImg = image.NewRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.NearestNeighbor.Scale(finalImg, finalImg.Rect, img, img.Bounds(), draw.Over, nil)
	return
}
