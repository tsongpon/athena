package logger

import (
	"os"
	"testing"

	"go.uber.org/zap"
)

func TestInitialize(t *testing.T) {
	// Test development mode (default)
	err := Initialize()
	if err != nil {
		t.Errorf("Initialize() unexpected error = %v", err)
	}

	if globalLogger == nil {
		t.Error("Initialize() should set globalLogger")
	}
}

func TestInitialize_Production(t *testing.T) {
	// Test production mode
	os.Setenv("APP_ENV", "production")
	defer os.Unsetenv("APP_ENV")

	err := Initialize()
	if err != nil {
		t.Errorf("Initialize() with production env unexpected error = %v", err)
	}

	if globalLogger == nil {
		t.Error("Initialize() should set globalLogger in production mode")
	}

	// Reset for other tests
	globalLogger = nil
}

func TestInitialize_CustomLogLevel(t *testing.T) {
	// Test custom log level
	os.Setenv("LOG_LEVEL", "warn")
	defer os.Unsetenv("LOG_LEVEL")

	err := Initialize()
	if err != nil {
		t.Errorf("Initialize() with custom log level unexpected error = %v", err)
	}

	if globalLogger == nil {
		t.Error("Initialize() should set globalLogger with custom log level")
	}

	// Reset for other tests
	globalLogger = nil
}

func TestGet(t *testing.T) {
	// Reset globalLogger
	globalLogger = nil

	// Get should return a logger even if Initialize wasn't called
	logger := Get()
	if logger == nil {
		t.Error("Get() should return a logger")
	}

	// Reset for other tests
	globalLogger = nil
}

func TestGet_AfterInitialize(t *testing.T) {
	err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	logger := Get()
	if logger == nil {
		t.Error("Get() should return the initialized logger")
	}

	if logger != globalLogger {
		t.Error("Get() should return the same instance as globalLogger")
	}

	// Reset for other tests
	globalLogger = nil
}

func TestLoggingFunctions(t *testing.T) {
	// Initialize logger for testing
	err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Test logging functions (should not panic)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Logging functions should not panic: %v", r)
		}
	}()

	Info("test info message", zap.String("key", "value"))
	Debug("test debug message", zap.Int("number", 42))
	Warn("test warn message", zap.Bool("flag", true))
	Error("test error message", zap.String("error", "test error"))

	// Reset for other tests
	globalLogger = nil
}

func TestWith(t *testing.T) {
	err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	childLogger := With(zap.String("service", "test"))
	if childLogger == nil {
		t.Error("With() should return a child logger")
	}

	// Child logger should be different from parent
	if childLogger == globalLogger {
		t.Error("With() should return a new logger instance")
	}

	// Reset for other tests
	globalLogger = nil
}

func TestSync(t *testing.T) {
	err := Initialize()
	if err != nil {
		t.Fatalf("Initialize() failed: %v", err)
	}

	// Sync should not return error for stdout/stderr
	err = Sync()
	if err != nil {
		// Note: Sync() may return error on some systems when syncing stdout/stderr
		// This is expected behavior and not a real error
		t.Logf("Sync() returned error (may be expected): %v", err)
	}

	// Reset for other tests
	globalLogger = nil
}

func TestSync_WithoutInitialize(t *testing.T) {
	globalLogger = nil

	// Sync should not panic or error if logger is nil
	err := Sync()
	if err != nil {
		t.Errorf("Sync() without initialized logger should not error, got: %v", err)
	}
}
