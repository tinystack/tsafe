# safe
Go goroutine safety Package

[![Go Report Card](https://goreportcard.com/badge/github.com/tinystack/tsafe)](https://goreportcard.com/report/github.com/    tinystack/tsafe)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.18-61CFDD.svg?style=flat-square)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/tinystack/tsafe)](https://pkg.go.dev/mod/github.com/tinystack/tsafe)

## 安装

go get -u github.com/tinystack/tsafe

## 示例

```go
import "github.com/tinystack/tsafe"

// 安全的启动goroutine
// 当发生 panic 时会捕捉相关报错信息并通过 SetLogger 设置的 logger 实例输出报错信息和堆栈信息
Go(func() {
    panic("test panic")
})

// 安全的启动goroutine
// 当发生 panic 时会执行对应的recover函数
GoWithRecover(func() {
    panic("test panic")
}, func(err any) {
    fmt.Println(err)
})
```

### API

- safe.SetLogger(l Logger)
- safe.Go(goroutine func())
- safe.GoWithRecover(goroutine func(), customRecover func(err any))