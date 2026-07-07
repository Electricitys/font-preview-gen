package renderer

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"

	"github.com/skrashevich/go-webp" // Pure Go WebP Encoder
	xfont "golang.org/x/image/font"  // Aliased to avoid collision
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func RenderWebP(fontBytes []byte, text string) ([]byte, error) {
	parsedFont, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    26,
		DPI:     72,
		Hinting: xfont.HintingFull,
	})
	if err != nil {
		return nil, err
	}
	defer face.Close()

	width, height := 450, 60
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.Transparent), image.Point{}, draw.Src)

	textColor := image.NewUniform(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	metrics := face.Metrics()
	fontHeight := metrics.Ascent + metrics.Descent
	baselineY := fixed.I(height)/2 + fontHeight/2 - metrics.Descent

	drawer := &xfont.Drawer{
		Dst:  dst,
		Src:  textColor,
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(15), Y: baselineY},
	}
	drawer.DrawString(text)

	var buf bytes.Buffer
	err = webp.Encode(&buf, dst, &webp.Options{
		Lossy:   true, // Exact matching field
		Quality: 85,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
