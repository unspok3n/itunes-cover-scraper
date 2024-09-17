package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unspok3n/itunes-cover-scraper/pkg/itunes"
)

const coverFilename = "cover.jpg"

func main() {
	var input string
	if len(os.Args) < 2 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter search query: ")
		input, _ = reader.ReadString('\n')
	} else {
		_, filename := filepath.Split(os.Args[1])
		input = strings.TrimSuffix(filename, filepath.Ext(filename))
	}

	fmt.Println("Searching...")
	query := itunes.PrepareQuery(input)
	search, err := itunes.Search(query)
	if err != nil {
		fmt.Println("Search:", err)
	} else {
		if search.ResultCount > 0 {
			rawUrl := search.Results[0].ArtworkUrl100
			originalUrl := itunes.OriginalUrl(rawUrl)

			var coverUrl string
			if originalUrl != "" {
				coverUrl = originalUrl
			} else {
				coverUrl = strings.Replace(rawUrl, "100x100bb.jpg", "3000x3000bb.jpg", 1)
			}

			fmt.Println("Downloading...")
			err := DownloadFile(coverUrl, coverFilename)
			if err != nil {
				fmt.Println("Error downloading cover:", err)
			}

			fmt.Println("Removing metadata...")
			data, err := os.ReadFile(coverFilename)
			if err != nil {
				fmt.Println("Error reading cover:", err)
			}

			filtered, err := StripExif(data)
			if err != nil {
				if !errors.Is(err, ErrExifMarkerNotFound) {
					fmt.Println("Error removing EXIF metadata:", err)
				}
			} else {
				if err := os.WriteFile(coverFilename, filtered, 0644); err != nil {
					fmt.Println("Error saving cover file:", err)
				}
			}

		} else {
			fmt.Println("No results found")
		}
	}

	fmt.Println("Press enter to exit")
	fmt.Scanln()
}
