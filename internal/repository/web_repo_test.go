package repository

import (
	"net/http"
	"net/http/httptest"
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
	title, err := repo.GetTitle(server.URL)
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
	_, err := repo.GetTitle("")
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
	title, err := repo.GetTitle(invalidURL)
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
	title, err := repo.GetTitle(server.URL)
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
	title, err := repo.GetTitle(server.URL)
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
	title, err := repo.GetTitle(server.URL)
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
	title, err := repo.GetTitle(server.URL)
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
	title, err := repo.GetTitle(server.URL)
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
	title, err := repo.GetTitle(server.URL)
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
	title, err := repo.GetTitle(server.URL)
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
