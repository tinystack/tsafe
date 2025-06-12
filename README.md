# TSafe - Panic-Safe Goroutine Launcher

[English](README.md) | [中文](README_CN.md)

TSafe is a lightweight Go library that makes goroutines panic-proof. Start goroutines safely with automatic panic recovery and customizable error handling.

[![Go Report Card](https://goreportcard.com/badge/github.com/tinystack/tsafe)](https://goreportcard.com/report/github.com/tinystack/tsafe)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/tinystack/tsafe)](https://pkg.go.dev/mod/github.com/tinystack/tsafe)

## Features

- **Panic-Proof Goroutines**: Automatically catch and handle panics in goroutines
- **Custom Error Handling**: Pluggable logger interface and custom recovery functions
- **Thread-Safe**: All operations are thread-safe and suitable for concurrent use
- **Zero Dependencies**: No external dependencies (only testify for testing)
- **Minimal Overhead**: Simple API with ~220ns per goroutine overhead

## Installation

```bash
go get -u github.com/tinystack/tsafe
```

## Quick Start

### Basic Safe Goroutine Execution

```go
package main

import (
    "fmt"
    "github.com/tinystack/tsafe"
)

func main() {
    // Safe goroutine with automatic panic recovery
    tsafe.Go(func() {
        panic("This won't crash your program!")
    })

    // Safe goroutine with custom panic handling
    tsafe.GoWithRecover(func() {
        panic("Custom handling")
    }, func(err any) {
        fmt.Printf("Caught panic: %v\n", err)
    })
}
```

### Custom Logger

```go
type MyLogger struct{}

func (l *MyLogger) Print(err, stack any) {
    log.Printf("Custom Logger - Error: %v\nStack: %s", err, stack)
}

func main() {
    // Set custom logger
    tsafe.SetLogger(&MyLogger{})

    tsafe.Go(func() {
        panic("This will be logged by MyLogger")
    })
}
```

## API Reference

### Core Functions

#### `Go(goroutine func())`

Starts a goroutine with automatic panic recovery. Panics are logged using the configured logger.

```go
tsafe.Go(func() {
    // Your code here
    panic("This will be caught safely")
})
```

#### `GoWithRecover(goroutine func(), customRecover func(err any))`

Starts a goroutine with custom panic recovery handling.

```go
tsafe.GoWithRecover(func() {
    // Your code here
    panic("Custom handling")
}, func(err any) {
    fmt.Printf("Caught: %v\n", err)
})
```

#### `SetLogger(l Logger)`

Sets a custom logger for panic handling. The logger must implement the `Logger` interface.

```go
tsafe.SetLogger(&MyCustomLogger{})
```

### Logger Interface

```go
type Logger interface {
    Print(err, stack any) // Print logs an error and its stack trace
}
```

## Best Practices

1. **Use for fire-and-forget goroutines**: Perfect for goroutines that shouldn't crash your application
2. **Custom logging**: Implement custom loggers for better error tracking and monitoring
3. **Resource cleanup**: Ensure proper resource cleanup in your goroutine functions
4. **Error handling**: Use `GoWithRecover` when you need specific error handling logic

## Examples

Check out the [examples](examples/) directory for comprehensive usage examples:

- [Basic Usage](examples/basic/main.go) - Simple panic recovery examples
- [Custom Logger Implementation](examples/logger/main.go) - Various logger implementations

## Testing

Run the test suite:

```bash
go test -v ./...
```

Run benchmarks:

```bash
go test -bench=. -benchmem
```

## Performance

TSafe is designed for high performance with minimal overhead:

- Goroutine creation overhead: ~220ns per goroutine
- Memory allocation: 24B per goroutine

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) for details on how to submit pull requests, report issues, and contribute to the project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed list of changes and version history.

## Support

- **Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/tinystack/tsafe)
- **Issues**: [GitHub Issues](https://github.com/tinystack/tsafe/issues)
- **Discussions**: [GitHub Discussions](https://github.com/tinystack/tsafe/discussions)

---

Made with ❤️ by the TSafe team. If you find this project useful, please consider giving it a ⭐ on GitHub!
