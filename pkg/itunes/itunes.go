package itunes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type Result struct {
	ArtworkUrl100 string `json:"artworkUrl100"`
}

type SearchResponse struct {
	ResultCount int      `json:"resultCount"`
	Results     []Result `json:"results"`
}

func OriginalUrl(url string) string {
	if strings.HasSuffix(url, ".jpg/100x100bb.jpg") {
		parts := strings.Split(url, "/100x100bb.jpg")
		if len(parts) > 1 {
			baseUrl := parts[0]
			newUrl := strings.Replace(baseUrl, "is1-ssl.mzstatic.com/image/thumb", "a1.mzstatic.com/r40", 1)
			return newUrl
		}
	}
	return ""
}

func PrepareQuery(query string) string {
	query = strings.TrimSpace(query)
	query = strings.ToLower(query)
	re := regexp.MustCompile(`^\d+[-.]?\s*`)
	query = re.ReplaceAllString(query, "")
	query = strings.ReplaceAll(query, "-", " ")
	query = strings.ReplaceAll(query, "_", " ")
	attributesToRemove := []string{
		"(radio)", "(radio edit)", "(extended)", "(extended mix)", "(original)", "(original mix)",
		"(instrumental)", "(instrumental mix)", "(radio version)", "(extended version)", "(pro mix)",
	}
	for _, t := range attributesToRemove {
		query = strings.ReplaceAll(query, t, "")
	}
	for strings.Contains(query, "  ") {
		query = strings.ReplaceAll(query, "  ", " ")
	}
	return query
}

var (
	ErrInvalidJsonp = errors.New("invalid jsonp")
)

func Search(query string) (*SearchResponse, error) {
	endpoint := fmt.Sprintf(
		"https://itunes.apple.com/search?callback=jQuery20309382196007363668_1726162953706&entity=song,album&media=music&entity=song&term=%s&country=nz&limit=1",
		url.QueryEscape(query),
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("fetch error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %w", err)
	}

	bodyStr := string(body)
	start := strings.Index(bodyStr, "(")
	end := strings.LastIndex(bodyStr, ")")
	if start == -1 || end == -1 || start >= end {
		fmt.Println(bodyStr)
		return nil, ErrInvalidJsonp
	}
	jsonData := bodyStr[start+1 : end]

	var apiResponse SearchResponse
	err = json.Unmarshal([]byte(jsonData), &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("json parse error: %w", err)
	}

	return &apiResponse, nil
}
