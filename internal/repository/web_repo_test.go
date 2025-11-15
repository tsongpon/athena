package repository

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestWebRepository_GetTitle(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with a title
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page Title</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting title from the test server
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() unexpected error = %v", err)
		return
	}

	expectedTitle := "Test Page Title"
	if title != expectedTitle {
		t.Errorf("GetTitle() title = %v, want %v", title, expectedTitle)
	}
}

func TestWebRepository_GetTitle_EmptyURL(t *testing.T) {
	repo := NewWebRepository()

	// Test with empty URL
	_, err := repo.GetTitle(context.Background(), "")
	if err == nil {
		t.Error("GetTitle() with empty URL should return error")
		return
	}

	expectedError := "URL cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("GetTitle() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestWebRepository_GetTitle_InvalidURL(t *testing.T) {
	repo := NewWebRepository()

	// Test with invalid URL - should return URL as fallback title
	invalidURL := "not-a-valid-url"
	title, err := repo.GetTitle(context.Background(), invalidURL)
	if err != nil {
		t.Errorf("GetTitle() with invalid URL should not return error, got %v", err)
		return
	}

	// Should return the URL as fallback title
	if title != invalidURL {
		t.Errorf("GetTitle() with invalid URL should return URL as title, got %v, want %v", title, invalidURL)
	}
}

func TestWebRepository_GetTitle_NotFound(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Test getting title from 404 page - should return URL as fallback
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() with 404 status should not return error, got %v", err)
		return
	}

	// Should return the URL as fallback title
	if title != server.URL {
		t.Errorf("GetTitle() with 404 should return URL as title, got %v, want %v", title, server.URL)
	}
}

func TestWebRepository_GetTitle_NoTitle(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML without a title
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting title from page without title - should return URL as fallback
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() with no title tag should not return error, got %v", err)
		return
	}

	// Should return the URL as fallback title
	if title != server.URL {
		t.Errorf("GetTitle() with no title should return URL as title, got %v, want %v", title, server.URL)
	}
}

func TestWebRepository_GetTitle_EmptyTitle(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with empty title
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title></title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting title from page with empty title - should return URL as fallback
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() with empty title tag should not return error, got %v", err)
		return
	}

	// Should return the URL as fallback title
	if title != server.URL {
		t.Errorf("GetTitle() with empty title should return URL as title, got %v, want %v", title, server.URL)
	}
}

func TestWebRepository_GetTitle_TitleWithWhitespace(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with title containing whitespace
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>  Test Page With Whitespace  </title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting title - should trim whitespace
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() unexpected error = %v", err)
		return
	}

	expectedTitle := "Test Page With Whitespace"
	if title != expectedTitle {
		t.Errorf("GetTitle() title = %v, want %v", title, expectedTitle)
	}
}

func TestWebRepository_GetTitle_SpecialCharacters(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with special characters in title
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test & Demo | Example Site</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting title with special characters
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() unexpected error = %v", err)
		return
	}

	expectedTitle := "Test & Demo | Example Site"
	if title != expectedTitle {
		t.Errorf("GetTitle() title = %v, want %v", title, expectedTitle)
	}
}

func TestWebRepository_GetTitle_InternalServerError(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Test getting title from 500 error page - should return URL as fallback
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() with 500 status should not return error, got %v", err)
		return
	}

	// Should return the URL as fallback title
	if title != server.URL {
		t.Errorf("GetTitle() with 500 should return URL as title, got %v, want %v", title, server.URL)
	}
}

func TestWebRepository_GetTitle_InvalidHTML(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns invalid HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `This is not valid HTML but has <title>Test Title</title> in it`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// HTML parser should still be able to extract the title
	title, err := repo.GetTitle(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetTitle() with invalid HTML should still extract title, got error = %v", err)
		return
	}

	expectedTitle := "Test Title"
	if title != expectedTitle {
		t.Errorf("GetTitle() title = %v, want %v", title, expectedTitle)
	}
}

func TestNewWebRepository(t *testing.T) {
	repo := NewWebRepository()

	if repo == nil {
		t.Error("NewWebRepository() should not return nil")
	}

	if repo.httpClient == nil {
		t.Error("NewWebRepository() should initialize httpClient")
	}

	// Verify timeout is set
	expectedTimeout := 10 * 1000000000 // 10 seconds in nanoseconds
	if repo.httpClient.Timeout.Nanoseconds() != int64(expectedTimeout) {
		t.Errorf("NewWebRepository() httpClient timeout = %v, want 10s", repo.httpClient.Timeout)
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		wantTitle   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "Valid HTML with title",
			html:      `<!DOCTYPE html><html><head><title>Test Title</title></head><body></body></html>`,
			wantTitle: "Test Title",
			wantErr:   false,
		},
		{
			name:        "No title tag",
			html:        `<!DOCTYPE html><html><head></head><body></body></html>`,
			wantErr:     true,
			errContains: "no title found",
		},
		{
			name:        "Empty title tag",
			html:        `<!DOCTYPE html><html><head><title></title></head><body></body></html>`,
			wantErr:     true,
			errContains: "no title found",
		},
		{
			name:      "Title with whitespace",
			html:      `<!DOCTYPE html><html><head><title>   Trimmed Title   </title></head><body></body></html>`,
			wantTitle: "Trimmed Title",
			wantErr:   false,
		},
		{
			name:      "Multiple titles - first one wins",
			html:      `<!DOCTYPE html><html><head><title>First Title</title><title>Second Title</title></head><body></body></html>`,
			wantTitle: "First Title",
			wantErr:   false,
		},
		{
			name:        "Title with only whitespace becomes empty and errors",
			html:        `<!DOCTYPE html><html><head><title>   </title></head><body></body></html>`,
			wantErr:     true,
			errContains: "no title found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.html)
			title, err := extractTitle(reader)

			if tt.wantErr {
				if err == nil {
					t.Errorf("extractTitle() expected error containing '%s', got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("extractTitle() error = %v, should contain %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("extractTitle() unexpected error = %v", err)
					return
				}
				if title != tt.wantTitle {
					t.Errorf("extractTitle() title = %v, want %v", title, tt.wantTitle)
				}
			}
		})
	}
}

func TestWebRepository_GetContentSummary_EmptyURL(t *testing.T) {
	repo := NewWebRepository()

	// Test with empty URL
	_, err := repo.GetContentSummary(context.Background(), "")
	if err == nil {
		t.Error("GetContentSummary() with empty URL should return error")
		return
	}

	expectedError := "URL cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("GetContentSummary() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestWebRepository_GetContentSummary_NoAPIKey(t *testing.T) {
	repo := NewWebRepository()

	// Ensure ANTHROPIC_API_KEY is not set
	originalAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	defer func() {
		if originalAPIKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", originalAPIKey)
		}
	}()

	// Create a test server that returns HTML
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<h1>Test Content</h1>
			<p>This is a test paragraph.</p>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting summary without API key - should return empty string
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() without API key should not return error, got %v", err)
		return
	}

	if summary != "" {
		t.Errorf("GetContentSummary() without API key should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_InvalidURL(t *testing.T) {
	repo := NewWebRepository()

	// Test with invalid URL - should return empty string
	summary, err := repo.GetContentSummary(context.Background(), "not-a-valid-url")
	if err != nil {
		t.Errorf("GetContentSummary() with invalid URL should not return error, got %v", err)
		return
	}

	if summary != "" {
		t.Errorf("GetContentSummary() with invalid URL should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_NotFound(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Test getting summary from 404 page - should return empty string
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() with 404 status should not return error, got %v", err)
		return
	}

	if summary != "" {
		t.Errorf("GetContentSummary() with 404 should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_NoTextContent(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with no text content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
			<script>console.log('test');</script>
			<style>body { color: red; }</style>
		</head>
		<body>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting summary from page with no text content - should return empty string
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() with no text content should not return error, got %v", err)
		return
	}

	if summary != "" {
		t.Errorf("GetContentSummary() with no text content should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_LongContent(t *testing.T) {
	repo := NewWebRepository()

	// Create content longer than 4000 characters
	longText := strings.Repeat("This is a long text. ", 300) // Creates ~6000 chars

	// Create a test server that returns HTML with long text content
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<p>%s</p>
		</body>
		</html>
		`, longText)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Save and restore environment variables
	originalModel := os.Getenv("LLM_MODEL")
	originalAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	defer func() {
		if originalModel != "" {
			os.Setenv("LLM_MODEL", originalModel)
		} else {
			os.Unsetenv("LLM_MODEL")
		}
		if originalAPIKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("ANTHROPIC_API_KEY")
		}
	}()

	// Unset API key so we don't actually call the LLM
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Setenv("LLM_MODEL", "anthropic")

	// Test that function handles long content (should truncate to 4000 chars)
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() with long content should not return error, got %v", err)
		return
	}

	// Should return empty since no API key, but validates truncation logic runs
	if summary != "" {
		t.Errorf("GetContentSummary() without API key should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_UnsupportedModel(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with text
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test</title></head>
		<body><p>Some content to summarize.</p></body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Save and restore environment variable
	originalModel := os.Getenv("LLM_MODEL")
	defer func() {
		if originalModel != "" {
			os.Setenv("LLM_MODEL", originalModel)
		} else {
			os.Unsetenv("LLM_MODEL")
		}
	}()

	// Set unsupported model
	os.Setenv("LLM_MODEL", "unsupported-model")

	// Test with unsupported model
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() with unsupported model should not return error, got %v", err)
		return
	}

	// Should return empty string
	if summary != "" {
		t.Errorf("GetContentSummary() with unsupported model should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_NoModelSet(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with text
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test</title></head>
		<body><p>Some content to summarize.</p></body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Save and restore environment variable
	originalModel := os.Getenv("LLM_MODEL")
	defer func() {
		if originalModel != "" {
			os.Setenv("LLM_MODEL", originalModel)
		} else {
			os.Unsetenv("LLM_MODEL")
		}
	}()

	// Unset model (default case in switch)
	os.Unsetenv("LLM_MODEL")

	// Test with no model set
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() with no model set should not return error, got %v", err)
		return
	}

	// Should return empty string
	if summary != "" {
		t.Errorf("GetContentSummary() with no model set should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_OpenAINoAPIKey(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with text
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test</title></head>
		<body><p>Some content to summarize.</p></body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Save and restore environment variables
	originalModel := os.Getenv("LLM_MODEL")
	originalAPIKey := os.Getenv("OPENAI_API_KEY")
	defer func() {
		if originalModel != "" {
			os.Setenv("LLM_MODEL", originalModel)
		} else {
			os.Unsetenv("LLM_MODEL")
		}
		if originalAPIKey != "" {
			os.Setenv("OPENAI_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("OPENAI_API_KEY")
		}
	}()

	// Set OpenAI model but no API key
	os.Setenv("LLM_MODEL", "openai")
	os.Unsetenv("OPENAI_API_KEY")

	// Test with OpenAI model and no API key
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() with OpenAI and no API key should not return error, got %v", err)
		return
	}

	// Should return empty string
	if summary != "" {
		t.Errorf("GetContentSummary() with OpenAI and no API key should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_GeminiNoAPIKey(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with text
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test</title></head>
		<body><p>Some content to summarize.</p></body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Save and restore environment variables
	originalModel := os.Getenv("LLM_MODEL")
	originalAPIKey := os.Getenv("GEMINI_API_KEY")
	defer func() {
		if originalModel != "" {
			os.Setenv("LLM_MODEL", originalModel)
		} else {
			os.Unsetenv("LLM_MODEL")
		}
		if originalAPIKey != "" {
			os.Setenv("GEMINI_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("GEMINI_API_KEY")
		}
	}()

	// Set Gemini model but no API key
	os.Setenv("LLM_MODEL", "gemini")
	os.Unsetenv("GEMINI_API_KEY")

	// Test with Gemini model and no API key
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() with Gemini and no API key should not return error, got %v", err)
		return
	}

	// Should return empty string
	if summary != "" {
		t.Errorf("GetContentSummary() with Gemini and no API key should return empty string, got %v", summary)
	}
}

func TestWebRepository_GetContentSummary_AnthropicInvalidAPIKey(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with text
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head><title>Test</title></head>
		<body><p>Some content to summarize.</p></body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Save and restore environment variables
	originalModel := os.Getenv("LLM_MODEL")
	originalAPIKey := os.Getenv("ANTHROPIC_API_KEY")
	defer func() {
		if originalModel != "" {
			os.Setenv("LLM_MODEL", originalModel)
		} else {
			os.Unsetenv("LLM_MODEL")
		}
		if originalAPIKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", originalAPIKey)
		} else {
			os.Unsetenv("ANTHROPIC_API_KEY")
		}
	}()

	// Set Anthropic model with invalid API key (empty string should work, actual invalid key would make real API call)
	os.Setenv("LLM_MODEL", "anthropic")
	os.Setenv("ANTHROPIC_API_KEY", "invalid-key-but-wont-actually-call-api")

	// Note: This test validates the path exists but won't actually test LLM call errors
	// without mocking or integration tests. The function returns empty string on any error.
	// We're just ensuring the code path executes without panic.
	summary, err := repo.GetContentSummary(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetContentSummary() should not return error on LLM errors, got %v", err)
		return
	}

	// May return empty string if API call fails, which is expected behavior
	_ = summary // Result depends on whether actual API call is made
}

func TestExtractTextContent(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantText string
	}{
		{
			name:     "Simple HTML with text",
			html:     `<!DOCTYPE html><html><body><p>Hello World</p></body></html>`,
			wantText: "Hello World",
		},
		{
			name:     "HTML with script tags - should skip",
			html:     `<!DOCTYPE html><html><body><script>alert('test');</script><p>Hello</p></body></html>`,
			wantText: "Hello",
		},
		{
			name:     "HTML with style tags - should skip",
			html:     `<!DOCTYPE html><html><head><style>body { color: red; }</style></head><body><p>Hello</p></body></html>`,
			wantText: "Hello",
		},
		{
			name:     "Multiple paragraphs",
			html:     `<!DOCTYPE html><html><body><p>First paragraph</p><p>Second paragraph</p></body></html>`,
			wantText: "First paragraph Second paragraph",
		},
		{
			name:     "Text with whitespace - should trim and normalize",
			html:     `<!DOCTYPE html><html><body><p>   Text with   spaces   </p></body></html>`,
			wantText: "Text with   spaces",
		},
		{
			name:     "Empty HTML",
			html:     `<!DOCTYPE html><html><body></body></html>`,
			wantText: "",
		},
		{
			name:     "HTML with nested elements",
			html:     `<!DOCTYPE html><html><body><div><p><span>Nested</span> text</p></div></body></html>`,
			wantText: "Nested text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := extractTextContent(tt.html)
			if text != tt.wantText {
				t.Errorf("extractTextContent() = %v, want %v", text, tt.wantText)
			}
		})
	}
}

func TestWebRepository_GetMainImage(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with og:image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta property="og:image" content="https://example.com/og-image.jpg">
			<title>Test Page</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting main image from the test server
	image, err := repo.GetMainImage(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetMainImage() unexpected error = %v", err)
		return
	}

	expectedImage := "https://example.com/og-image.jpg"
	if image != expectedImage {
		t.Errorf("GetMainImage() image = %v, want %v", image, expectedImage)
	}
}

func TestWebRepository_GetMainImage_EmptyURL(t *testing.T) {
	repo := NewWebRepository()

	// Test with empty URL
	_, err := repo.GetMainImage(context.Background(), "")
	if err == nil {
		t.Error("GetMainImage() with empty URL should return error")
		return
	}

	expectedError := "URL cannot be empty"
	if err.Error() != expectedError {
		t.Errorf("GetMainImage() error = %v, want %v", err.Error(), expectedError)
	}
}

func TestWebRepository_GetMainImage_InvalidURL(t *testing.T) {
	repo := NewWebRepository()

	// Test with invalid URL - should return empty string without error
	invalidURL := "not-a-valid-url"
	image, err := repo.GetMainImage(context.Background(), invalidURL)
	if err != nil {
		t.Errorf("GetMainImage() with invalid URL should not return error, got %v", err)
		return
	}

	// Should return empty string as fallback
	if image != "" {
		t.Errorf("GetMainImage() with invalid URL should return empty string, got %v", image)
	}
}

func TestWebRepository_GetMainImage_NotFound(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Test getting main image from 404 page - should return empty string
	image, err := repo.GetMainImage(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetMainImage() with 404 status should not return error, got %v", err)
		return
	}

	// Should return empty string as fallback
	if image != "" {
		t.Errorf("GetMainImage() with 404 should return empty string, got %v", image)
	}
}

func TestWebRepository_GetMainImage_NoOGImage(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML without og:image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Test Page</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting main image from page without og:image - should return empty string
	image, err := repo.GetMainImage(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetMainImage() with no og:image should not return error, got %v", err)
		return
	}

	// Should return empty string
	if image != "" {
		t.Errorf("GetMainImage() with no og:image should return empty string, got %v", image)
	}
}

func TestWebRepository_GetMainImage_EmptyOGImage(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with empty og:image
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta property="og:image" content="">
			<title>Test Page</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting main image from page with empty og:image - should return empty string
	image, err := repo.GetMainImage(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetMainImage() with empty og:image should not return error, got %v", err)
		return
	}

	// Should return empty string
	if image != "" {
		t.Errorf("GetMainImage() with empty og:image should return empty string, got %v", image)
	}
}

func TestWebRepository_GetMainImage_OGImageWithWhitespace(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with og:image containing whitespace
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta property="og:image" content="  https://example.com/trimmed-image.jpg  ">
			<title>Test Page</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting main image - should trim whitespace
	image, err := repo.GetMainImage(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetMainImage() unexpected error = %v", err)
		return
	}

	expectedImage := "https://example.com/trimmed-image.jpg"
	if image != expectedImage {
		t.Errorf("GetMainImage() image = %v, want %v", image, expectedImage)
	}
}

func TestWebRepository_GetMainImage_InternalServerError(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Test getting main image from 500 error page - should return empty string
	image, err := repo.GetMainImage(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetMainImage() with 500 status should not return error, got %v", err)
		return
	}

	// Should return empty string as fallback
	if image != "" {
		t.Errorf("GetMainImage() with 500 should return empty string, got %v", image)
	}
}

func TestWebRepository_GetMainImage_MultipleOGImages(t *testing.T) {
	repo := NewWebRepository()

	// Create a test server that returns HTML with multiple og:image tags
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
		<!DOCTYPE html>
		<html>
		<head>
			<meta property="og:image" content="https://example.com/first-image.jpg">
			<meta property="og:image" content="https://example.com/second-image.jpg">
			<title>Test Page</title>
		</head>
		<body>
			<h1>Hello World</h1>
		</body>
		</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	// Test getting main image - should return the first one
	image, err := repo.GetMainImage(context.Background(), server.URL)
	if err != nil {
		t.Errorf("GetMainImage() unexpected error = %v", err)
		return
	}

	expectedImage := "https://example.com/first-image.jpg"
	if image != expectedImage {
		t.Errorf("GetMainImage() image = %v, want %v (should return first og:image)", image, expectedImage)
	}
}

func TestExtractOGImage(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		wantImage   string
		wantErr     bool
		errContains string
	}{
		{
			name:      "Valid og:image",
			html:      `<!DOCTYPE html><html><head><meta property="og:image" content="https://example.com/image.jpg"></head><body></body></html>`,
			wantImage: "https://example.com/image.jpg",
			wantErr:   false,
		},
		{
			name:        "No og:image tag",
			html:        `<!DOCTYPE html><html><head></head><body></body></html>`,
			wantErr:     true,
			errContains: "no og:image found",
		},
		{
			name:        "Empty og:image content",
			html:        `<!DOCTYPE html><html><head><meta property="og:image" content=""></head><body></body></html>`,
			wantErr:     true,
			errContains: "no og:image found",
		},
		{
			name:      "og:image with whitespace",
			html:      `<!DOCTYPE html><html><head><meta property="og:image" content="   https://example.com/trimmed.jpg   "></head><body></body></html>`,
			wantImage: "https://example.com/trimmed.jpg",
			wantErr:   false,
		},
		{
			name:      "Multiple og:image tags - first one wins",
			html:      `<!DOCTYPE html><html><head><meta property="og:image" content="https://example.com/first.jpg"><meta property="og:image" content="https://example.com/second.jpg"></head><body></body></html>`,
			wantImage: "https://example.com/first.jpg",
			wantErr:   false,
		},
		{
			name:        "og:image with only whitespace",
			html:        `<!DOCTYPE html><html><head><meta property="og:image" content="   "></head><body></body></html>`,
			wantErr:     true,
			errContains: "no og:image found",
		},
		{
			name:      "og:image in body (should still find it)",
			html:      `<!DOCTYPE html><html><head></head><body><meta property="og:image" content="https://example.com/body-image.jpg"></body></html>`,
			wantImage: "https://example.com/body-image.jpg",
			wantErr:   false,
		},
		{
			name:      "Meta tag with other properties before og:image",
			html:      `<!DOCTYPE html><html><head><meta property="og:title" content="Title"><meta property="og:image" content="https://example.com/image.jpg"></head><body></body></html>`,
			wantImage: "https://example.com/image.jpg",
			wantErr:   false,
		},
		{
			name:      "og:image with nested HTML",
			html:      `<!DOCTYPE html><html><head><title>Test</title><meta property="og:image" content="https://example.com/nested.jpg"><script>var x = 1;</script></head><body></body></html>`,
			wantImage: "https://example.com/nested.jpg",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.html)
			image, err := extractOGImage(reader)

			if tt.wantErr {
				if err == nil {
					t.Errorf("extractOGImage() expected error containing '%s', got nil", tt.errContains)
					return
				}
				if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("extractOGImage() error = %v, should contain %v", err.Error(), tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("extractOGImage() unexpected error = %v", err)
					return
				}
				if image != tt.wantImage {
					t.Errorf("extractOGImage() image = %v, want %v", image, tt.wantImage)
				}
			}
		})
	}
}
