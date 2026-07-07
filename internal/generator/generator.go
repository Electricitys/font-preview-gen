package generator

import (
	"os"
	"strings"

	"undangon/font-preview-gen/internal/downloader"
	"undangon/font-preview-gen/internal/renderer"
)

// Generate returns a WebP preview image for the given font source and text.
// fontSource can be a URL, a local file path ending with .ttf/.otf, or a
// fontsource name that will be resolved via downloader.FetchFontFromFontsource.
func Generate(fontSource, text string) ([]byte, error) {
	var fontBytes []byte
	var err error

	if strings.HasPrefix(fontSource, "http://") || strings.HasPrefix(fontSource, "https://") {
		fontBytes, err = downloader.FetchFontFromURL(fontSource)
	} else if strings.HasSuffix(fontSource, ".ttf") || strings.HasSuffix(fontSource, ".otf") {
		fontBytes, err = os.ReadFile(fontSource)
	} else {
		fontBytes, err = downloader.FetchFontFromFontsource(fontSource)
	}
	if err != nil {
		return nil, err
	}
	return renderer.RenderWebP(fontBytes, text)
}
