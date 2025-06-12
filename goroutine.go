// Package tsafe provides utilities for safe goroutine execution with panic recovery
package tsafe

import (
	"log"
	"runtime/debug"
	"sync"
)

// Logger defines the interface for custom error logging
// Users can implement this interface to customize how goroutine errors are logged
type Logger interface {
	// Print logs an error and its stack trace
	Print(err, stack any)
}

// defaultLoggerImpl is the default implementation of Logger interface
type defaultLoggerImpl struct{}

// Print implements the Logger interface for defaultLoggerImpl
// It logs errors using the standard log package
func (l *defaultLoggerImpl) Print(err, stack any) {
	log.Printf("Error in goroutine: %s\nStack trace: %s\n", err, stack)
}

// Thread-safe global logger management
var (
	defaultLogger Logger = &defaultLoggerImpl{}
	loggerMutex   sync.RWMutex
)

// SetLogger sets a custom logger for goroutine error handling
// This function is thread-safe
func SetLogger(l Logger) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	defaultLogger = l
}

// getLogger returns the current logger in a thread-safe manner
func getLogger() Logger {
	loggerMutex.RLock()
	defer loggerMutex.RUnlock()
	return defaultLogger
}

// Go starts a goroutine with automatic panic recovery
// When a panic occurs, it will be logged using the configured logger
// This is the most convenient way to start a safe goroutine
func Go(goroutine func()) {
	GoWithRecover(goroutine, func(err any) {
		if logger := getLogger(); logger != nil {
			logger.Print(err, debug.Stack())
		}
	})
}

// GoWithRecover starts a goroutine with custom panic recovery handling
// Parameters:
//   - goroutine: the function to execute in the goroutine
//   - customRecover: the function to handle panic recovery (called if panic occurs)
//
// This provides more control over error handling compared to Go()
func GoWithRecover(goroutine func(), customRecover func(err any)) {
	if goroutine == nil {
		return // Avoid creating goroutine for nil function
	}

	go func() {
		defer func() {
			if err := recover(); err != nil && customRecover != nil {
				customRecover(err)
			}
		}()
		goroutine()
	}()
}
