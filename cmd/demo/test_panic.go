package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"
)

func main() {
    fmt.Println("=== Panic Handling Test ===")
    
    // Test: Function A with try/catch calling Function B that panics
    fmt.Println("\nTesting: Function A (with try/catch) → Function B (panic)")
    
    Try(func() {
        fmt.Println("1. Starting Function A...")
        functionA()
        fmt.Println("6. This line should NOT be reached")
    }).Any(func(ex Exception) {
        fmt.Printf("5. Panic caught as exception: %s\n", ex.Error())
        fmt.Printf("   Exception type: %s\n", ex.TypeName())
    }).Finally(func() {
        fmt.Println("7. Finally block executed (cleanup)")
    })
    
    fmt.Println("\n=== Test completed ===")
    fmt.Println("Expected flow: 1 → 2 → 3 → 4 → 5 → 7")
    fmt.Println("Line 6 should NOT appear (execution stopped by panic)")
}

// Function A: Has no try/catch, just calls Function B
func functionA() {
    fmt.Println("2. Inside Function A")
    fmt.Println("3. Function A calling Function B...")
    
    functionB() // This will panic
    
    fmt.Println("   This line in Function A should NOT be reached")
}

// Function B: Generates a native Go panic
func functionB() {
    fmt.Println("4. Inside Function B - about to panic!")
    
    // Native Go panic
    panic("BOOM! Native Go panic from Function B")
    
    fmt.Println("   This line in Function B should NOT be reached")
}
