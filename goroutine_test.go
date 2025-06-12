package tsafe

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Mock logger for testing
type mockLogger struct {
	lastError any
	lastStack any
	callCount int
	mutex     sync.Mutex
}

func (m *mockLogger) Print(err, stack any) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.lastError = err
	m.lastStack = stack
	m.callCount++
}

func (m *mockLogger) getLastError() any {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.lastError
}

func (m *mockLogger) getCallCount() int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.callCount
}

func TestGo(t *testing.T) {
	t.Run("should not panic when goroutine panics", func(t *testing.T) {
		assert.NotPanics(t, func() {
			Go(func() {
				panic("test panic")
			})
		})

		// Give goroutine time to execute
		time.Sleep(10 * time.Millisecond)
	})

	t.Run("should handle nil function gracefully", func(t *testing.T) {
		assert.NotPanics(t, func() {
			GoWithRecover(nil, nil)
		})
	})

	t.Run("should use custom logger", func(t *testing.T) {
		mock := &mockLogger{}
		originalLogger := getLogger()
		SetLogger(mock)
		defer SetLogger(originalLogger)

		done := make(chan struct{})
		Go(func() {
			defer func() { done <- struct{}{} }()
			panic("test panic with custom logger")
		})

		select {
		case <-done:
			// Wait a bit more for logger to be called
			time.Sleep(10 * time.Millisecond)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("goroutine did not complete in time")
		}

		assert.Equal(t, "test panic with custom logger", mock.getLastError())
		assert.Equal(t, 1, mock.getCallCount())
	})
}

func TestGoWithRecover(t *testing.T) {
	t.Run("should execute custom recover function", func(t *testing.T) {
		var recoveredErr any
		waitChan := make(chan struct{})

		assert.NotPanics(t, func() {
			GoWithRecover(func() {
				panic("test panic")
			}, func(err any) {
				recoveredErr = err
				waitChan <- struct{}{}
			})
		})

		select {
		case <-waitChan:
			assert.Equal(t, "test panic", recoveredErr)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("custom recover function was not called")
		}
	})

	t.Run("should handle normal execution without panic", func(t *testing.T) {
		executed := false
		waitChan := make(chan struct{})

		GoWithRecover(func() {
			executed = true
			waitChan <- struct{}{}
		}, func(err any) {
			t.Errorf("recover function should not be called for normal execution")
		})

		select {
		case <-waitChan:
			assert.True(t, executed)
		case <-time.After(100 * time.Millisecond):
			t.Fatal("goroutine did not complete")
		}
	})

	t.Run("should handle nil recover function", func(t *testing.T) {
		assert.NotPanics(t, func() {
			GoWithRecover(func() {
				panic("test panic")
			}, nil)
		})
		time.Sleep(10 * time.Millisecond)
	})
}

func TestSetLogger(t *testing.T) {
	t.Run("should be thread-safe", func(t *testing.T) {
		originalLogger := getLogger()
		defer SetLogger(originalLogger)

		var wg sync.WaitGroup
		loggers := make([]*mockLogger, 100)

		// Create 100 goroutines setting different loggers
		for i := 0; i < 100; i++ {
			loggers[i] = &mockLogger{}
			wg.Add(1)
			go func(logger *mockLogger) {
				defer wg.Done()
				SetLogger(logger)
			}(loggers[i])
		}

		wg.Wait()

		// Verify that one of the loggers is set (thread-safe operation completed)
		currentLogger := getLogger()
		assert.NotEqual(t, originalLogger, currentLogger)
	})
}

// Benchmark tests
func BenchmarkGo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Go(func() {
			// Do minimal work
		})
	}
}

func BenchmarkGoWithRecover(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GoWithRecover(func() {
			// Do minimal work
		}, func(err any) {
			// Handle error
		})
	}
}
