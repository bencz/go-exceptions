package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"
)

// Simple demonstration of the exception system
func SimpleDemo() {
    fmt.Println("=== Simple Exception Demo ===")
    
    // Example 1: Basic Try/Catch
    fmt.Println("\n1. Basic ArgumentNull exception:")
    Try(func() {
        ThrowArgumentNull("username", "Username is required")
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("   Caught: %s - %s\n", ex.ParamName, ex.Message)
        }),
    )
    
    // Example 2: Multiple exception types
    fmt.Println("\n2. Multiple exception types:")
    for i := 0; i < 3; i++ {
        Try(func() {
            if i%2 == 0 {
                ThrowArgumentNull("param", "Parameter cannot be null")
            } else {
                ThrowInvalidOperation("Operation not allowed")
            }
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                fmt.Printf("   ArgumentNull: %s\n", ex.ParamName)
            }),
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                fmt.Printf("   InvalidOperation: %s\n", ex.Message)
            }),
        )
    }
    
    // Example 3: With Finally block
    fmt.Println("\n3. With Finally block:")
    Try(func() {
        ThrowInvalidOperation("Something went wrong")
    }).Handle(
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("   Exception handled: %s\n", ex.Message)
        }),
    ).Finally(func() {
        fmt.Println("   Cleanup executed")
    })
    
    // Example 4: Range validation
    fmt.Println("\n4. Range validation:")
    Try(func() {
        age := -5
        ThrowArgumentOutOfRange("age", age, "Age must be positive")
    }).Handle(
        Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
            fmt.Printf("   Range error: %s = %v (%s)\n", ex.ParamName, ex.Value, ex.Message)
        }),
    )
    
    // Example 5: Null checking with ThrowIfNil
    fmt.Println("\n5. Null checking:")
    Try(func() {
        var data *string = nil
        ThrowIfNil("data", data)
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("   Null check failed: %s\n", ex.ParamName)
        }),
    )
    
    // Example 6: Generic Exception handler using Any()
    fmt.Println("\n6. Generic Exception handler with Any():")
    Try(func() {
        ThrowInvalidOperation("Invalid operation")
    }).Any(func(ex Exception) {
        // This catches ANY exception type
        fmt.Printf("   Caught with Any(): %s (Type: %s)\n", ex.Error(), ex.TypeName())
    })
    
    // Example 7: Generic Exception handler inside Handle()
    fmt.Println("\n7. Generic Exception handler inside Handle():")
    Try(func() {
        ThrowArgumentNull("data", "Data cannot be null")
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("   Specific handler: %s\n", ex.ParamName)
        }),
        HandlerAny(func(ex Exception) {
            // This also catches ANY exception type, but inside Handle()
            fmt.Printf("   Generic handler: %s (Type: %s)\n", ex.Error(), ex.TypeName())
        }),
    )
    
    // Example 8: Mixed handlers - specific and generic
    fmt.Println("\n8. Mixed handlers (specific + generic):")
    for i := 0; i < 3; i++ {
        Try(func() {
            switch i {
            case 0:
                ThrowArgumentNull("param", "Parameter is null")
            case 1:
                ThrowInvalidOperation("Invalid operation")
            case 2:
                ThrowArgumentOutOfRange("index", -1, "Index out of range")
            }
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                fmt.Printf("   Specific ArgumentNull: %s\n", ex.ParamName)
            }),
            HandlerAny(func(ex Exception) {
                fmt.Printf("   Generic fallback: %s (Type: %s)\n", ex.Error(), ex.TypeName())
            }),
        )
    }
    
    fmt.Println("\nSimple demo completed!")
}
