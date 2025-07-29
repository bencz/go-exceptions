package main

import (
    "fmt"
    "strings"
    "time"
    . "github.com/bencz/go-exceptions"
)

func main() {
    fmt.Println("=== Advanced Exception Examples ===")
    
    // Run all advanced examples
    runBusinessLogicExample()
    runStackTraceExample()
    runNestedExceptionsExample()
    runPerformanceExample()
    runComplexValidationExample()
    runRetryPatternExample()
    runPanicHandlingExample()
    
    fmt.Println("\nAdvanced examples completed!")
}

// ============================================================================
// BUSINESS LOGIC EXAMPLE
// ============================================================================

func runBusinessLogicExample() {
    fmt.Println("\n1. Business Logic Validation:")
    Try(func() {
        // Simulate complex business logic
        processUserRegistration("", "user@email.com", 17)
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("   Validation error: %s is required\n", ex.ParamName)
        }),
        Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
            fmt.Printf("   Range error: %s = %v is invalid\n", ex.ParamName, ex.Value)
        }),
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("   Business rule violation: %s\n", ex.Message)
        }),
    ).Any(func(ex Exception) {
        fmt.Printf("   Unexpected error: %s\n", ex.Error())
    }).Finally(func() {
        fmt.Println("   Registration cleanup executed")
    })
}

// ============================================================================
// STACK TRACE EXAMPLE
// ============================================================================

func runStackTraceExample() {
    fmt.Println("\n2. Stack Trace Demonstration:")
    Try(func() {
        deepFunction()
    }).Handle(
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("   Exception: %s\n", ex.Message)
            fmt.Printf("   Stack trace (first 3 calls):\n")
            for i, trace := range full.StackTrace {
                if i < 3 {
                    fmt.Printf("     %s\n", trace)
                }
            }
        }),
    )
}

// ============================================================================
// NESTED EXCEPTIONS EXAMPLE
// ============================================================================

func runNestedExceptionsExample() {
    fmt.Println("\n3. Nested Exceptions:")
    Try(func() {
        // Simulate an operation that fails and is wrapped in another exception
        var innerException *Exception
        
        Try(func() {
            ThrowFileError("database.db", "Connection failed", fmt.Errorf("network timeout"))
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
            fmt.Printf("   Main exception: %s\n", ex.Message)
            
            if full.HasInnerException() {
                fmt.Printf("   Full message: %s\n", full.GetFullMessage())
                
                // Look for FileException in the chain
                if fileEx := FindInnerException[FileException](&full); fileEx != nil {
                    fmt.Printf("   Problematic file: %s\n", fileEx.Filename)
                }
                
                fmt.Printf("   Total exceptions in chain: %d\n", len(full.GetAllExceptions()))
            }
        }),
    )
}

// ============================================================================
// PERFORMANCE EXAMPLE
// ============================================================================

func runPerformanceExample() {
    fmt.Println("\n4. Performance Test (Type Cache):")
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
    fmt.Printf("   1000 Try/Catch operations executed in: %v\n", elapsed)
    fmt.Printf("   Performance optimized with type caching\n")
}

// ============================================================================
// COMPLEX VALIDATION EXAMPLE
// ============================================================================

func runComplexValidationExample() {
    fmt.Println("\n5. Complex Validation:")
    Try(func() {
        validateComplexObject(nil, "", -1, nil)
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            fmt.Printf("   Null parameter: %s\n", ex.ParamName)
        }),
        Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
            fmt.Printf("   Value out of range: %s = %v\n", ex.ParamName, ex.Value)
        }),
    ).Any(func(ex Exception) {
        fmt.Printf("   Unhandled error: %s\n", ex.Error())
    })
}

// ============================================================================
// RETRY PATTERN EXAMPLE
// ============================================================================

func runRetryPatternExample() {
    fmt.Println("\n6. Retry Pattern with Exceptions:")
    
    maxRetries := 3
    for attempt := 1; attempt <= maxRetries; attempt++ {
        var success bool
        
        Try(func() {
            // Simulate operation that fails first 2 times
            if attempt < 3 {
                ThrowNetworkError("Connection failed", 500)
            }
            success = true
            fmt.Printf("   Operation succeeded on attempt %d\n", attempt)
        }).Handle(
            Handler[NetworkException](func(ex NetworkException, full Exception) {
                fmt.Printf("   Attempt %d failed: %s (Status: %d)\n", attempt, ex.Message, ex.StatusCode)
                if attempt == maxRetries {
                    fmt.Printf("   Max retries reached, operation failed\n")
                }
            }),
        )
        
        if success {
            break
        }
    }
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

func processUserRegistration(username, email string, age int) {
    ThrowIf(username == "", ArgumentNullException{
        ParamName: "username",
        Message:   "Username is required for registration",
    })
    
    ThrowIf(age < 18, ArgumentOutOfRangeException{
        ParamName: "age",
        Value:     age,
        Message:   "Must be 18 or older to register",
    })
    
    // Simulate other validations
    if strings.Contains(email, "fake") {
        ThrowInvalidOperation("Email domain not allowed")
    }
}

func deepFunction() {
    anotherFunction()
}

func anotherFunction() {
    ThrowInvalidOperation("Error in deep function call")
}

func validateComplexObject(obj interface{}, name string, age int, config map[string]string) {
    ThrowIfNil("obj", obj)
    ThrowIfNil("config", config)
    
    ThrowIf(name == "", ArgumentNullException{
        ParamName: "name",
        Message:   "Name cannot be empty",
    })
    
    ThrowIf(age < 0, ArgumentOutOfRangeException{
        ParamName: "age",
        Value:     age,
        Message:   "Age cannot be negative",
    })
}

// ============================================================================
// PANIC HANDLING EXAMPLE
// ============================================================================

func runPanicHandlingExample() {
    fmt.Println("\n7. Panic Handling Test:")
    
    // Test 1: Function A with try/catch calling Function B that panics
    fmt.Println("\n   Test 1: Try/Catch capturing native Go panic:")
    Try(func() {
        advancedFunctionA()
    }).Handle(
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("     Caught InvalidOperation: %s\n", ex.Message)
        }),
    ).Any(func(ex Exception) {
        fmt.Printf("     Caught unexpected exception: %s (Type: %s)\n", ex.Error(), ex.TypeName())
    })
    
    // Test 2: Function with try/catch calling function that panics with string
    fmt.Println("\n   Test 2: Try/Catch capturing string panic:")
    Try(func() {
        functionWithStringPanic()
    }).Any(func(ex Exception) {
        fmt.Printf("     Caught panic as exception: %s\n", ex.Error())
    })
    
    // Test 3: Function with try/catch calling function that panics with custom error
    fmt.Println("\n   Test 3: Try/Catch capturing custom error panic:")
    Try(func() {
        functionWithErrorPanic()
    }).Any(func(ex Exception) {
        fmt.Printf("     Caught error panic: %s\n", ex.Error())
    })
    
    // Test 4: Nested try/catch with panic
    fmt.Println("\n   Test 4: Nested try/catch with panic:")
    Try(func() {
        Try(func() {
            functionWithNilPointerPanic()
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                fmt.Printf("     Inner catch: %s\n", ex.ParamName)
            }),
        ).Any(func(ex Exception) {
            fmt.Printf("     Inner generic catch: %s\n", ex.Error())
            // Re-throw as our exception
            ThrowInvalidOperation("Wrapped panic: " + ex.Error())
        })
    }).Handle(
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            fmt.Printf("     Outer catch: %s\n", ex.Message)
        }),
    )
}

// Function A that calls Function B (renamed to avoid conflicts)
func advancedFunctionA() {
    fmt.Printf("     Function A: calling Function B...\n")
    advancedFunctionB()
    fmt.Printf("     Function A: this should not print\n")
}

// Function B that panics (renamed to avoid conflicts)
func advancedFunctionB() {
    fmt.Printf("     Function B: about to panic...\n")
    panic("Native Go panic from Function B!")
}

// Function that panics with string
func functionWithStringPanic() {
    fmt.Printf("     About to panic with string...\n")
    panic("This is a string panic!")
}

// Function that panics with error
func functionWithErrorPanic() {
    fmt.Printf("     About to panic with error...\n")
    panic(fmt.Errorf("this is an error panic: %s", "custom error message"))
}

// Function that causes nil pointer panic
func functionWithNilPointerPanic() {
    fmt.Printf("     About to cause nil pointer panic...\n")
    var ptr *string
    fmt.Printf("     Value: %s\n", *ptr) // This will panic
}
