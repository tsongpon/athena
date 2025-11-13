package repository

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/tsongpon/athena/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/net/html"
)

type WebRepository struct {
	httpClient *http.Client
}

func NewWebRepository() *WebRepository {
	return &WebRepository{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetTitle fetches the HTML content from the given URL and extracts the title
func (r *WebRepository) GetTitle(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Fetch the URL
	resp, err := r.httpClient.Get(url)
	if err != nil {
		logger.Debug("failed to fetch title from URL", zap.String("url", url), zap.Error(err))
		return url, nil
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		logger.Debug("failed to fetch title from URL", zap.String("url", url), zap.Int("status_code", resp.StatusCode))
		return url, nil
	}

	// Parse the HTML and extract the title
	title, err := extractTitle(resp.Body)
	if err != nil {
		logger.Debug("failed to fetch title from URL", zap.String("url", url), zap.Error(err))
		return url, nil
	}

	return title, nil
}

// GetMainImage fetches the HTML content from the given URL and extracts the OpenGraph image URL
func (r *WebRepository) GetMainImage(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Fetch the URL
	resp, err := r.httpClient.Get(url)
	if err != nil {
		logger.Debug("failed to fetch main image from URL", zap.String("url", url), zap.Error(err))
		return "", nil
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		logger.Debug("failed to fetch main image from URL", zap.String("url", url), zap.Int("status_code", resp.StatusCode))
		return "", nil
	}

	// Parse the HTML and extract the OG image
	imageURL, err := extractOGImage(resp.Body)
	if err != nil {
		logger.Debug("failed to extract OG image from URL", zap.String("url", url), zap.Error(err))
		return "", nil
	}

	return imageURL, nil
}

// extractOGImage parses HTML and extracts the content of the og:image meta tag
func extractOGImage(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	var imageURL string
	var found bool
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if found {
			return
		}
		if n.Type == html.ElementNode && n.Data == "meta" {
			var property, content string
			for _, attr := range n.Attr {
				if attr.Key == "property" && attr.Val == "og:image" {
					property = attr.Val
				}
				if attr.Key == "content" {
					content = attr.Val
				}
			}
			if property == "og:image" && content != "" {
				imageURL = strings.TrimSpace(content)
				found = true
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found {
				return
			}
			traverse(c)
		}
	}
	traverse(doc)

	if imageURL == "" {
		return "", fmt.Errorf("no og:image found")
	}

	return imageURL, nil
}

// extractTitle parses HTML and extracts the content of the <title> tag
func extractTitle(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	var title string
	var found bool
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if found {
			return
		}
		if n.Type == html.ElementNode && n.Data == "title" {
			if n.FirstChild != nil {
				title = strings.TrimSpace(n.FirstChild.Data)
				found = true
			}
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if found {
				return
			}
			traverse(c)
		}
	}
	traverse(doc)

	if title == "" {
		return "", fmt.Errorf("no title found")
	}

	return title, nil
}
