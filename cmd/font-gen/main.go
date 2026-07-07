package main

import (
	"fmt"
	"os"
	"strings"

	"undangon/font-preview-gen/internal/downloader"
	"undangon/font-preview-gen/internal/renderer"
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
		fontBytes, err = downloader.FetchFontFromURL(fontSource)
	} else if strings.HasSuffix(fontSource, ".ttf") || strings.HasSuffix(fontSource, ".otf") {
		fontBytes, err = os.ReadFile(fontSource)
	} else {
		fontBytes, err = downloader.FetchFontFromFontsource(fontSource)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading font: %v\n", err)
		os.Exit(1)
	}

	// 2. Generate the WebP Preview Image
	webpData, err := renderer.RenderWebP(fontBytes, text)
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
