package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/tinystack/tsafe"
)

// CustomLogger implements the tsafe.Logger interface with custom formatting
type CustomLogger struct {
	prefix string
}

func (cl *CustomLogger) Print(err, stack any) {
	log.Printf("%s ERROR: %v\n%s STACK:\n%s",
		cl.prefix, err, cl.prefix, stack)
}

// FileLogger logs errors to a file
type FileLogger struct {
	file *os.File
}

func NewFileLogger(filename string) (*FileLogger, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &FileLogger{file: file}, nil
}

func (fl *FileLogger) Print(err, stack any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(fl.file, "[%s] ERROR: %v\nSTACK:\n%s\n---\n",
		timestamp, err, stack)
}

func (fl *FileLogger) Close() error {
	return fl.file.Close()
}

// StructuredLogger provides structured logging
type StructuredLogger struct {
	component string
}

func (sl *StructuredLogger) Print(err, stack any) {
	// Extract just the panic location from stack trace
	stackStr := fmt.Sprintf("%s", stack)
	lines := strings.Split(stackStr, "\n")
	location := "unknown"
	if len(lines) > 3 {
		location = strings.TrimSpace(lines[3])
	}

	log.Printf(`{"level":"ERROR","component":"%s","error":"%v","location":"%s","timestamp":"%s"}`,
		sl.component, err, location, time.Now().Format(time.RFC3339))
}

// MetricsLogger keeps track of error statistics
type MetricsLogger struct {
	errorCount    int
	panicTypes    map[string]int
	lastErrorTime time.Time
}

func NewMetricsLogger() *MetricsLogger {
	return &MetricsLogger{
		panicTypes: make(map[string]int),
	}
}

func (ml *MetricsLogger) Print(err, stack any) {
	ml.errorCount++
	ml.lastErrorTime = time.Now()

	errorType := fmt.Sprintf("%T", err)
	ml.panicTypes[errorType]++

	log.Printf("METRICS ERROR #%d (type: %s): %v",
		ml.errorCount, errorType, err)
}

func (ml *MetricsLogger) GetStats() (int, map[string]int, time.Time) {
	return ml.errorCount, ml.panicTypes, ml.lastErrorTime
}

func main() {
	fmt.Println("=== TSafe Custom Logger Examples ===")

	// Example 1: Custom formatted logger
	fmt.Println("\n1. Custom formatted logger:")
	customLogger := &CustomLogger{prefix: "[CUSTOM]"}
	tsafe.SetLogger(customLogger)

	tsafe.Go(func() {
		panic("This error will be logged with custom formatting")
	})
	time.Sleep(100 * time.Millisecond)

	// Example 2: File logger
	fmt.Println("\n2. File logger:")
	fileLogger, err := NewFileLogger("tsafe_errors.log")
	if err != nil {
		fmt.Printf("   Failed to create file logger: %v\n", err)
	} else {
		tsafe.SetLogger(fileLogger)

		tsafe.Go(func() {
			panic("This error will be logged to a file")
		})
		time.Sleep(100 * time.Millisecond)

		fileLogger.Close()
		fmt.Println("   Error logged to tsafe_errors.log")
	}

	// Example 3: Structured logger
	fmt.Println("\n3. Structured JSON logger:")
	structuredLogger := &StructuredLogger{component: "tsafe-example"}
	tsafe.SetLogger(structuredLogger)

	tsafe.Go(func() {
		panic("Structured logging example")
	})
	time.Sleep(100 * time.Millisecond)

	// Example 4: Metrics logger
	fmt.Println("\n4. Metrics tracking logger:")
	metricsLogger := NewMetricsLogger()
	tsafe.SetLogger(metricsLogger)

	// Generate different types of panics
	errorTypes := []interface{}{
		"string panic",
		fmt.Errorf("error panic"),
		123,
		[]string{"slice", "panic"},
	}

	for _, errorType := range errorTypes {
		errorValue := errorType
		tsafe.Go(func() {
			panic(errorValue)
		})
		time.Sleep(50 * time.Millisecond)
	}

	// Show metrics
	count, types, lastError := metricsLogger.GetStats()
	fmt.Printf("   Total errors: %d\n", count)
	fmt.Printf("   Error types: %v\n", types)
	fmt.Printf("   Last error at: %s\n", lastError.Format("15:04:05"))

	// Example 5: Logger switching
	fmt.Println("\n5. Dynamic logger switching:")

	// Start with one logger
	logger1 := &CustomLogger{prefix: "[LOGGER-1]"}
	tsafe.SetLogger(logger1)

	tsafe.Go(func() {
		panic("Error with logger 1")
	})
	time.Sleep(50 * time.Millisecond)

	// Switch to another logger
	logger2 := &CustomLogger{prefix: "[LOGGER-2]"}
	tsafe.SetLogger(logger2)

	tsafe.Go(func() {
		panic("Error with logger 2")
	})
	time.Sleep(50 * time.Millisecond)

	// Example 6: Logger performance comparison
	fmt.Println("\n6. Logger performance comparison:")

	loggers := []struct {
		name   string
		logger tsafe.Logger
	}{
		{"Simple", &CustomLogger{prefix: "[SIMPLE]"}},
		{"Structured", &StructuredLogger{component: "perf-test"}},
		{"Metrics", NewMetricsLogger()},
	}

	for _, loggerInfo := range loggers {
		fmt.Printf("   Testing %s logger:\n", loggerInfo.name)
		tsafe.SetLogger(loggerInfo.logger)

		start := time.Now()
		for i := 0; i < 5; i++ {
			tsafe.Go(func() {
				panic("performance test panic")
			})
		}
		time.Sleep(100 * time.Millisecond)
		duration := time.Since(start)

		fmt.Printf("   %s logger: %v for 5 panics\n",
			loggerInfo.name, duration.Truncate(time.Millisecond))
	}

	// Example 7: Multiple concurrent loggers (demonstrating thread safety)
	fmt.Println("\n7. Concurrent logger switching test:")

	for i := 0; i < 10; i++ {
		go func(id int) {
			logger := &CustomLogger{prefix: fmt.Sprintf("[CONCURRENT-%d]", id)}
			tsafe.SetLogger(logger)

			tsafe.Go(func() {
				panic(fmt.Sprintf("Concurrent panic from goroutine %d", id))
			})
		}(i)
	}

	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== All logger examples completed ===")

	// Clean up - remove log file
	os.Remove("tsafe_errors.log")
}
