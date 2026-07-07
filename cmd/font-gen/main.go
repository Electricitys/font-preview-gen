// cmd/font-gen/main.go
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"undangon/font-preview-gen/internal/downloader"
	"undangon/font-preview-gen/internal/renderer"
)

func main() {
	// 1. Bind flags for both short (-o) and long (--output) formats
	var outputPath string
	flag.StringVar(&outputPath, "o", "", "Path to save the output WebP file (bypasses stdout mapping)")
	flag.StringVar(&outputPath, "output", "", "Path to save the output WebP file (bypasses stdout mapping)")

	// Customize usage errors so they always route to Stderr
	flag.CommandLine.SetOutput(os.Stderr)
	flag.Parse()

	// 2. Extract remaining positional variables after parsing flags
	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: font-gen [options] <font_path_or_url_or_slug> <text>")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fontSource := args[0]
	text := args[1]

	var fontBytes []byte
	var err error

	// 3. Font routing pipeline
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

	// 4. Generate the WebP Preview Image
	webpData, err := renderer.RenderWebP(fontBytes, text)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error rendering image: %v\n", err)
		os.Exit(1)
	}

	// 5. Intelligent Output Delivery
	if outputPath != "" {
		// Scenario A: Flag is provided -> Save directly as a file to disk
		err = os.WriteFile(outputPath, webpData, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving file directly to %s: %v\n", outputPath, err)
			os.Exit(1)
		}
	} else {
		// Scenario B: No flag -> Stream raw data buffer to stdout (Supports TypeScript wrappers & shell operators)
		_, err = os.Stdout.Write(webpData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error streaming binary array to stdout: %v\n", err)
			os.Exit(1)
		}
	}
}
