package ipcamera

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

func createTextBanner(text string, width int, height int, fnt *sfnt.Font, fontSize float64, dpi float64, textColor color.Color, backgroundColor color.Color) (image.Image, error) {
	face, err := opentype.NewFace(fnt, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	defer func(face font.Face) {
		_ = face.Close()
	}(face)
	metrics := face.Metrics()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), image.NewUniform(backgroundColor), image.Point{}, draw.Src)

	bounds, _ := font.BoundString(face, text)
	boundsWidth := bounds.Max.X - bounds.Min.X
	boundsHeight := bounds.Max.Y - bounds.Min.Y

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(textColor),
		Face: face,
	}

	x := (fixed.I(width) - boundsWidth) / 2
	y := ((fixed.I(height) - boundsHeight) / 2) + metrics.Ascent
	d.Dot = fixed.Point26_6{
		X: x,
		Y: y,
	}

	d.DrawString(text)

	return img, nil
}

func getDisconnectedImage(fnt *sfnt.Font, width int, height int) (image.Image, error) {
	var err error
	if fnt == nil {
		fnt, err = sfnt.Parse(goregular.TTF)
		if err != nil {
			return nil, err
		}
	}
	img, err := createTextBanner("No signal", width, height, fnt, 48, 72, color.White, color.Black)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func getDisconnectedFrame(fnt *sfnt.Font, width int, height int) ([]byte, error) {
	img, err := getDisconnectedImage(fnt, width, height)
	if err != nil {
		return nil, err
	}
	imageBuffer := new(bytes.Buffer)
	err = jpeg.Encode(imageBuffer, img, nil)
	if err != nil {
		return nil, err
	}
	return imageBuffer.Bytes(), nil
}
