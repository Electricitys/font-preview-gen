package downloader

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tdewolff/font" // For parsing WOFF2 files
)

// FetchFontFromURL downloads any arbitrary direct or web font link (handles .woff2, .ttf, or .otf extensions seamlessly)
func FetchFontFromURL(targetURL string) ([]byte, error) {
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

func FetchFontFromFontsource(fontName string) ([]byte, error) {
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
