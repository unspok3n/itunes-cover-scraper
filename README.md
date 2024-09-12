Simple script that uses [iTunes Search API](https://developer.apple.com/library/archive/documentation/AudioVideo/Conceptual/iTuneSearchAPI/index.html) to find and automatically download full resolution music album covers in JPEG by either search query or an input file.

![Screenshot](/screenshot.png?raw=true "Screenshot")

Build
---
Just run:

    go build ./cmd/itunes-cover-scraper

Usage
---
Either open and type the search query or provide an input file, and the filename will be used as a search query.