package main

import (
    "errors"
    "fmt"
    "strings"
    "time"
    . "github.com/bencz/go-exceptions"
)

// Advanced examples demonstrating improved features

func AdvancedExamples() {
    fmt.Println("\n=== Advanced Example: Business Logic ===")
    Try(func() {
        // Simulate complex business logic
        ProcessUserRegistration("", "user@email.com", 17)
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("Validation error: %s is required\n", ex.ParamName)
        }),
        Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
            fmt.Printf("Range error: %s = %v is invalid\n", ex.ParamName, ex.Value)
        }),
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("Business rule violation: %s\n", ex.Message)
        }),
    ).Any(func(ex Exception) {
        fmt.Printf("Unexpected error: %s\n", ex.Error())
    }).Finally(func() {
        fmt.Println("Registration cleanup")
    })

    fmt.Println("\n=== Demonstration: Stack Trace ===")
    Try(func() {
        deepFunction()
    }).Handle(
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("Exception: %s\n", ex.Message)
            fmt.Printf("Stack trace (first 2 calls):\n")
            for i, trace := range full.StackTrace {
                if i < 2 {
                    fmt.Printf("   %s\n", trace)
                }
            }
        }),
    )
}

// Examples of implemented improvements
func ImprovementExamples() {
    fmt.Println("\n=== Examples of Improvements ===")
    
    // 1. ThrowIfNil Example
    fmt.Println("\n--- ThrowIfNil Example ---")
    Try(func() {
        var ptr *string = nil
        var slice []int = nil
        var m map[string]int = nil
        
        ThrowIfNil("ptr", ptr)
        ThrowIfNil("slice", slice)
        ThrowIfNil("map", m)
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("ThrowIfNil caught: %s\n", ex.ParamName)
        }),
    )
    
    // 2. Nested Exceptions Example
    fmt.Println("\n--- Nested Exceptions Example ---")
    Try(func() {
        // Simulate an operation that fails and is wrapped in another exception
        var innerException *Exception
        
        Try(func() {
            ThrowFileError("database.db", "Connection failed", errors.New("network timeout"))
        }).Handle(
            Handler[FileException](func(ex FileException, full Exception) {
                innerException = &full
            }),
        )
        
        // Throw a new exception with the previous one as inner
        ThrowWithInner(InvalidOperationException{
            Message: "Failed to initialize application",
        }, innerException)
    }).Handle(
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("Main exception: %s\n", ex.Message)
            
            if full.HasInnerException() {
                fmt.Printf("Full message: %s\n", full.GetFullMessage())
                
                // Look for FileException in the chain
                if fileEx := FindInnerException[FileException](&full); fileEx != nil {
                    fmt.Printf("Problematic file found: %s\n", fileEx.Filename)
                }
                
                fmt.Printf("Total exceptions in chain: %d\n", len(full.GetAllExceptions()))
            }
        }),
    )
    
    // 3. Performance Example with Cache
    fmt.Println("\n--- Performance Example (Type Cache) ---")
    start := time.Now()
    
    for i := 0; i < 1000; i++ {
        Try(func() {
            if i%2 == 0 {
                ThrowArgumentNull("test", "test message")
            } else {
                ThrowInvalidOperation("test operation")
            }
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                // Empty handler for performance test
            }),
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                // Empty handler for performance test
            }),
        )
    }
    
    elapsed := time.Since(start)
    fmt.Printf("1000 Try/Catch operations executed in: %v\n", elapsed)
    fmt.Printf("Type cache contains %d entries\n", len(typeCache))
    
    // 4. Complex Validation Example
    fmt.Println("\n--- Complex Validation Example ---")
    Try(func() {
        ValidateComplexObject(nil, "", -1, nil)
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("Null parameter: %s\n", ex.ParamName)
        }),
        Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
            fmt.Printf("Value out of range: %s = %v\n", ex.ParamName, ex.Value)
        }),
    ).Any(func(ex Exception) {
        fmt.Printf("Unhandled error: %s\n", ex.Error())
    })
}

// Helper function to demonstrate business logic
func ProcessUserRegistration(username, email string, age int) {
    ThrowIf(username == "", ArgumentNullException{
        ParamName: "username",
        Message:   "Username is required for registration",
    })
    
    ThrowIf(age < 18, ArgumentOutOfRangeException{
        ParamName: "age",
        Value:     age,
        Message:   "Must be 18 or older to register",
    })
    
    // Simulate other validations...
    if strings.Contains(email, "fake") {
        ThrowInvalidOperation("Email domain not allowed")
    }
}

// Function to demonstrate stack trace
func deepFunction() {
    anotherFunction()
}

func anotherFunction() {
    ThrowInvalidOperation("Error in deep function call")
}

// Function to demonstrate ThrowIfNil and complex validations
func ValidateComplexObject(obj interface{}, name string, age int, config map[string]string) {
    // Use ThrowIfNil for quick validations
    ThrowIfNil("obj", obj)
    ThrowIfNil("config", config)
    
    // Traditional validations
    ThrowIf(name == "", ArgumentNullException{
        ParamName: "name",
        Message:   "Name cannot be empty",
    })
    
    ThrowIf(age < 0, ArgumentOutOfRangeException{
        ParamName: "age",
        Value:     age,
        Message:   "Age must be non-negative",
    })
}
