# TSafe - 防崩溃 Goroutine 启动器

[English](README.md) | [中文](README_CN.md)

TSafe 是一个轻量级 Go 语言库，让 goroutine 远离 panic 崩溃。安全启动 goroutine，自动恢复 panic 并支持自定义错误处理。

[![Go Report Card](https://goreportcard.com/badge/github.com/tinystack/tsafe)](https://goreportcard.com/report/github.com/tinystack/tsafe)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/tinystack/tsafe)](https://pkg.go.dev/mod/github.com/tinystack/tsafe)

## 特性

- **防崩溃 Goroutine**: 自动捕获和处理 goroutine 中的 panic
- **自定义错误处理**: 可插拔日志接口和自定义恢复函数
- **线程安全**: 所有操作都是线程安全的，适合并发使用
- **零依赖**: 无外部依赖（仅测试时使用 testify）
- **极低开销**: 简单 API，每个 goroutine 仅约 220ns 开销

## 安装

```bash
go get -u github.com/tinystack/tsafe
```

## 快速开始

### 基础安全 Goroutine 执行

```go
package main

import (
    "fmt"
    "github.com/tinystack/tsafe"
)

func main() {
    // 安全的goroutine，自动panic恢复
    tsafe.Go(func() {
        panic("这不会让你的程序崩溃！")
    })

    // 安全的goroutine，自定义panic处理
    tsafe.GoWithRecover(func() {
        panic("自定义处理")
    }, func(err any) {
        fmt.Printf("捕获到panic: %v\n", err)
    })
}
```

### 自定义日志器

```go
type MyLogger struct{}

func (l *MyLogger) Print(err, stack any) {
    log.Printf("自定义日志器 - 错误: %v\n堆栈: %s", err, stack)
}

func main() {
    // 设置自定义日志器
    tsafe.SetLogger(&MyLogger{})

    tsafe.Go(func() {
        panic("这将由MyLogger记录")
    })
}
```

## API 参考

### 核心功能

#### `Go(goroutine func())`

启动一个带自动 panic 恢复的 goroutine。panic 将使用配置的日志器记录。

```go
tsafe.Go(func() {
    // 你的代码
    panic("这将被安全地捕获")
})
```

#### `GoWithRecover(goroutine func(), customRecover func(err any))`

启动一个带自定义 panic 恢复处理的 goroutine。

```go
tsafe.GoWithRecover(func() {
    // 你的代码
    panic("自定义处理")
}, func(err any) {
    fmt.Printf("捕获到: %v\n", err)
})
```

#### `SetLogger(l Logger)`

为 panic 处理设置自定义日志器。日志器必须实现`Logger`接口。

```go
tsafe.SetLogger(&MyCustomLogger{})
```

### Logger 接口

```go
type Logger interface {
    Print(err, stack any) // 记录错误和堆栈跟踪
}
```

## 最佳实践

1. **用于即发即忘的 goroutine**: 非常适合不应该让应用程序崩溃的 goroutine
2. **自定义日志**: 实现自定义日志器以便更好地进行错误跟踪和监控
3. **资源清理**: 确保在 goroutine 函数中正确清理资源
4. **错误处理**: 当需要特定错误处理逻辑时使用`GoWithRecover`

## 示例

查看[examples](examples/)目录获取全面的使用示例：

- [基础用法](examples/basic/main.go) - 简单的 panic 恢复示例
- [自定义日志器实现](examples/logger/main.go) - 各种日志器实现

## 测试

运行测试套件：

```bash
go test -v ./...
```

运行性能测试：

```bash
go test -bench=. -benchmem
```

## 性能

TSafe 专为高性能设计，开销极低：

- Goroutine 创建开销：每个 goroutine 约 220ns
- 内存分配：每个 goroutine 24B

## 贡献

欢迎贡献！请阅读我们的[贡献指南](CONTRIBUTING.md)，了解如何提交拉取请求、报告问题和为项目做出贡献的详细信息。

## 许可证

该项目使用 MIT 许可证 - 详情请查看[LICENSE](LICENSE)文件。

## 变更日志

查看[CHANGELOG.md](CHANGELOG.md)获取详细的变更列表和版本历史。

## 支持

- **文档**: [pkg.go.dev](https://pkg.go.dev/github.com/tinystack/tsafe)
- **问题**: [GitHub Issues](https://github.com/tinystack/tsafe/issues)
- **讨论**: [GitHub Discussions](https://github.com/tinystack/tsafe/discussions)

---

由 TSafe 团队用 ❤️ 制作。如果你觉得这个项目有用，请考虑在 GitHub 上给它一个 ⭐！
