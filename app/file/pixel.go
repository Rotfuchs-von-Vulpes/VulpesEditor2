package file

import (
	"image"
	"image/png"
	"io"

	"golang.org/x/image/draw"
)

func EnhancePixels(img image.Image, coeff int) (finalImg *image.RGBA) {
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()
	width := coeff * originalWidth
	height := coeff * originalHeight
	finalImg = image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(finalImg, finalImg.Rect, img, img.Bounds(), draw.Over, nil)
	return
}

func LoadImage(r io.Reader, coeff int32) (finalImg *image.RGBA, err error) {
	if coeff <= 0 {
		coeff = 1
	}
	img, err := png.Decode(r)
	if err != nil {
		return
	}
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()
	if originalWidth < 64 || originalHeight < 64 {
		finalImg = EnhancePixels(img, 10)
		return
	}
	if originalWidth < 128 || originalHeight < 128 {
		finalImg = EnhancePixels(img, 5)
		return
	}
	finalImg = image.NewRGBA(img.Bounds())
	draw.Draw(finalImg, finalImg.Bounds(), img, img.Bounds().Min, draw.Src)
	return
}

func LoadPixelImage(r io.Reader, coeff int) (finalImg *image.RGBA, err error) {
	if coeff <= 0 {
		coeff = 1
	}
	img, err := png.Decode(r)
	if err != nil {
		return
	}
	finalImg = EnhancePixels(img, coeff)
	return
}
