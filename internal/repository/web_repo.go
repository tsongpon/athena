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
