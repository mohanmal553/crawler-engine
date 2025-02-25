package modules

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/net/html"
)

// resolveURL resolves a relative URL against a base URL and returns the absolute URL
func resolveURL(ref, base string) string {
	parsedBase, err := url.Parse(base)
	if err != nil {
		return ""
	}

	parsedURL, err := url.Parse(ref)
	if err != nil {
		return ""
	}

	resolvedURL := parsedBase.ResolveReference(parsedURL)
	return resolvedURL.String()
}

// extractHref Extracts the href attribute from an anchor (<a>) tag
func extractHref(tag html.Token) string {
	for _, attr := range tag.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

// ExtractURL Extracts all the URL from the input HTML Code
func ExtractURL(body io.Reader, baseURL string) []string {
	tokenizer := html.NewTokenizer(body)
	links := make([]string, 0)

	for token := tokenizer.Next(); token != html.ErrorToken; token = tokenizer.Next() {
		tag := tokenizer.Token()

		if tag.Data != "a" {
			continue
		}

		url := extractHref(tag)
		if url == "" {
			continue
		}

		absoluteURL := resolveURL(url, baseURL)
		if absoluteURL != "" {
			links = append(links, absoluteURL)
		}
	}

	return links
}

// FetchData fetchs the Webpage Data from the input URL
//
// Throws an error if response status code is anything other than 200
func FetchData(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("HTTP Connection Error (ErrorCode: " + strconv.Itoa(resp.StatusCode) + ")")
	}
	return resp.Body, nil
}
