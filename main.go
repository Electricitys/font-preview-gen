package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/skrashevich/go-webp" // Pure Go WebP Encoder (No CGO errors)
	"github.com/tdewolff/font"       // For parsing WOFF2 files
	xfont "golang.org/x/image/font"  // Aliased to avoid collision
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: font-gen <font_path_or_fontsource_name_or_url> <text> <output_webp_path>")
		os.Exit(1)
	}

	fontSource := os.Args[1]
	text := os.Args[2]
	outputPath := os.Args[3]

	var fontBytes []byte
	var err error

	// 1. Core routing pipeline based on string matching characteristics
	if strings.HasPrefix(fontSource, "http://") || strings.HasPrefix(fontSource, "https://") {
		// Vector A: Fetch from direct url strings
		fontBytes, err = fetchFontFromURL(fontSource)
	} else if strings.HasSuffix(fontSource, ".ttf") || strings.HasSuffix(fontSource, ".otf") {
		// Vector B: Read from your application's disk/local file uploaded directories
		fontBytes, err = os.ReadFile(fontSource)
	} else {
		// Vector C: Fallback behavior defaults to parsing from Fontsource slug catalog names
		fontBytes, err = fetchFontFromFontsource(fontSource)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading font: %v\n", err)
		os.Exit(1)
	}

	// 2. Generate the WebP Preview Image
	webpData, err := renderWebP(fontBytes, text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering image: %v\n", err)
		os.Exit(1)
	}

	// 3. Write preview binary to disk destination
	err = os.WriteFile(outputPath, webpData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
		os.Exit(1)
	}
}

// Downloads any arbitrary direct or web font link (handles .woff2, .ttf, or .otf extensions seamlessly)
func fetchFontFromURL(targetURL string) ([]byte, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed fetching remote URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("remote URL returned non-200 status code: %d", resp.StatusCode)
	}

	rawBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading remote response body: %w", err)
	}

	// Check if data is WOFF2 by verifying its standard 4-byte magic signature ("wOF2")
	isWoff2 := len(rawBytes) >= 4 && string(rawBytes[:4]) == "wOF2"

	if strings.Contains(strings.ToLower(targetURL), ".woff2") || isWoff2 {
		ttfBytes, err := font.ParseWOFF2(rawBytes)
		if err != nil {
			return nil, fmt.Errorf("failed converting remote web-optimized woff2 format: %w", err)
		}
		return ttfBytes, nil
	}

	return rawBytes, nil
}

func fetchFontFromFontsource(fontName string) ([]byte, error) {
	cleanName := fontName
	if idx := strings.Index(cleanName, ":"); idx != -1 {
		cleanName = cleanName[:idx]
	}
	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(cleanName), " ", "-"))

	var urls []string
	if strings.Contains(strings.ToLower(fontName), "variable") || slug == "inter" {
		urls = []string{
			fmt.Sprintf("https://cdn.jsdelivr.net/fontsource/fonts/%s:vf@latest/latin-wght-normal.woff2", slug),
			fmt.Sprintf("https://cdn.jsdelivr.net/fontsource/fonts/%s@latest/latin-400-normal.woff2", slug),
		}
	} else {
		urls = []string{
			fmt.Sprintf("https://cdn.jsdelivr.net/fontsource/fonts/%s@latest/latin-400-normal.woff2", slug),
			fmt.Sprintf("https://cdn.jsdelivr.net/fontsource/fonts/%s:vf@latest/latin-wght-normal.woff2", slug),
		}
	}

	client := &http.Client{Timeout: 15 * time.Second}
	var woff2Bytes []byte
	var fetched bool

	for _, url := range urls {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			woff2Bytes, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				fetched = true
				break
			}
		}
		if resp != nil {
			resp.Body.Close()
		}
	}

	if !fetched {
		return nil, fmt.Errorf("could not retrieve fontsource asset payload for slug: %s", slug)
	}

	ttfBytes, err := font.ParseWOFF2(woff2Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed converting web-optimized woff2 format: %w", err)
	}

	return ttfBytes, nil
}

func renderWebP(fontBytes []byte, text string) ([]byte, error) {
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
		Lossy:   true,
		Quality: 85,
	})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
