package repository

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/anthropic"
	"github.com/tmc/langchaingo/llms/googleai"
	"github.com/tmc/langchaingo/llms/openai"
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

// GetContentSummary fetches the HTML content from the given URL and uses LangChain with Anthropic Claude
// to generate a summary within 500 characters
func (r *WebRepository) GetContentSummary(url string) (string, error) {
	if url == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}

	// Fetch the URL
	resp, err := r.httpClient.Get(url)
	if err != nil {
		logger.Debug("failed to fetch content from URL", zap.String("url", url), zap.Error(err))
		return "", nil
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		logger.Debug("failed to fetch content from URL", zap.String("url", url), zap.Int("status_code", resp.StatusCode))
		return "", nil
	}

	// Read the HTML body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("failed to read response body", zap.String("url", url), zap.Error(err))
		return "", nil
	}

	// Extract text content from HTML
	textContent := extractTextContent(string(bodyBytes))
	if textContent == "" {
		logger.Debug("no text content found", zap.String("url", url))
		return "", nil
	}

	// Limit the text content to avoid token limits (first 4000 characters)
	if len(textContent) > 4000 {
		textContent = textContent[:4000]
	}

	llmModelName := os.Getenv("LLM_MODEL")
	var llmModel llms.Model
	switch llmModelName {
	case "anthropic":
		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			logger.Debug("ANTHROPIC_API_KEY not set, skipping summarization")
			return "", nil
		}
		llmModel, err = anthropic.New(anthropic.WithToken(apiKey))
		if err != nil {
			logger.Debug("failed to create Anthropic client", zap.Error(err))
			return "", nil
		}
	case "openai":
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey == "" {
			logger.Debug("OPENAI_API_KEY not set, skipping summarization")
			return "", nil
		}
		llmModel, err = openai.New(openai.WithToken(apiKey))
		if err != nil {
			logger.Debug("failed to create OpenAI client", zap.Error(err))
			return "", nil
		}
	case "gemini":
		ctx := context.Background()
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			logger.Debug("GEMINI_API_KEY not set, skipping summarization")
			return "", nil
		}
		llmModel, err = googleai.New(ctx, googleai.WithAPIKey(apiKey))
		if err != nil {
			logger.Debug("failed to create Gemini client", zap.Error(err))
			return "", nil
		}
	default:
		logger.Debug("unsupported LLM model", zap.String("model", llmModelName))
		return "", nil
	}

	ctx := context.Background()
	prompt := fmt.Sprintf("Summarize the following website content in 1000 characters or less. Be concise and capture the main points:\n\n%s", textContent)

	summary, err := llms.GenerateFromSinglePrompt(ctx, llmModel, prompt)
	if err != nil {
		logger.Debug("failed to generate summary", zap.String("url", url), zap.Error(err))
		return "", nil
	}

	return strings.TrimSpace(summary), nil
}

// extractTextContent extracts readable text content from HTML
func extractTextContent(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return ""
	}

	var textBuilder strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		// Skip script and style tags
		if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style") {
			return
		}

		if n.Type == html.TextNode {
			text := strings.TrimSpace(n.Data)
			if text != "" {
				textBuilder.WriteString(text)
				textBuilder.WriteString(" ")
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)

	return strings.TrimSpace(textBuilder.String())
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
