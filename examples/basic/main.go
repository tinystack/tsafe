package main

import (
	"fmt"
	"time"

	"github.com/tinystack/tsafe"
)

func main() {
	fmt.Println("=== TSafe Basic Usage Examples ===")

	// Example 1: Basic safe goroutine execution
	fmt.Println("\n1. Basic safe goroutine with panic:")
	tsafe.Go(func() {
		fmt.Println("   Goroutine starting...")
		panic("This panic will be caught safely!")
	})

	// Give goroutine time to execute
	time.Sleep(100 * time.Millisecond)

	// Example 2: Safe goroutine with custom recovery
	fmt.Println("\n2. Safe goroutine with custom recovery:")
	tsafe.GoWithRecover(func() {
		fmt.Println("   Another goroutine starting...")
		panic("Custom recovery example")
	}, func(err any) {
		fmt.Printf("   Custom recovery caught: %v\n", err)
	})

	// Give goroutine time to execute
	time.Sleep(100 * time.Millisecond)

	// Example 3: Normal execution (no panic)
	fmt.Println("\n3. Normal execution without panic:")
	tsafe.Go(func() {
		fmt.Println("   This goroutine completes normally")
		fmt.Println("   No panic here!")
	})

	// Give goroutine time to execute
	time.Sleep(100 * time.Millisecond)

	// Example 4: Multiple concurrent safe goroutines
	fmt.Println("\n4. Multiple concurrent safe goroutines:")
	for i := 0; i < 5; i++ {
		id := i
		tsafe.Go(func() {
			if id%2 == 0 {
				fmt.Printf("   Goroutine %d: Normal execution\n", id)
			} else {
				fmt.Printf("   Goroutine %d: About to panic...\n", id)
				panic(fmt.Sprintf("Panic from goroutine %d", id))
			}
		})
	}

	// Give all goroutines time to execute
	time.Sleep(200 * time.Millisecond)

	fmt.Println("\n=== All examples completed successfully ===")
}
