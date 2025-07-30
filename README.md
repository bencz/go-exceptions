# Go Exception System

[![Go Version](https://img.shields.io/badge/Go-1.24+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/bencz/go-exceptions?status.svg)](https://godoc.org/github.com/bencz/go-exceptions)
[![Coverage](https://img.shields.io/badge/Coverage-97%25-brightgreen.svg)](https://github.com/bencz/go-exceptions)

A robust try/catch/finally exception handling system for Go using generics and reflection. Provides clean, type-safe exception handling with multiple syntaxes, nested exceptions, custom types, and performance optimization.

## Key Features

- **Type-Safe**: Uses Go generics for compile-time type safety
- **Multiple Syntaxes**: Catch function, Builder pattern, Handler interface
- **Nested Exceptions**: Full support for exception chaining like .NET
- **Custom Types**: Easy creation of domain-specific exception types
- **Performance Optimized**: Type caching for high-demand applications
- **Finally Blocks**: Guaranteed cleanup code execution
- **Helper Functions**: ThrowIf, ThrowIfNil for common validations
- **Stack Traces**: Automatic capture and formatting
- **Comprehensive Tests**: 97% test coverage with benchmarks

## Installation

```bash
go get github.com/bencz/go-exceptions
```

## Quick Start

```go
package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"
)

func main() {
    Try(func() {
        ThrowArgumentNull("param", "Parameter cannot be null")
    }).Catch(func(ex ArgumentNullException) {
        fmt.Printf("Caught: %s\n", ex.ParamName)
    })
}
```

## Features

- **Type Safety**: Use of generics for typed exception catching
- **Multiple Approaches**: 3 different syntaxes for different needs
- **Optimized Performance**: Type cache to avoid repeated reflection
- **Nested Exceptions**: Full support for inner exceptions
- **Stack Trace**: Automatic stack trace capture
- **Quick Validations**: Helper functions like `ThrowIfNil`

## Available Exception Types

```go
ArgumentNullException        // Null parameters
ArgumentOutOfRangeException  // Values out of range
InvalidOperationException    // Invalid operations
FileException               // File errors
NetworkException            // Network errors
```

## Creating Custom Exception Types

You can easily create your own exception types by implementing the `ExceptionType` interface:

```go
// Custom exception for database operations
type DatabaseException struct {
    Query     string
    Message   string
    ErrorCode int
}

func (e DatabaseException) Error() string {
    return fmt.Sprintf("DatabaseException[%d]: %s (Query: %s)", e.ErrorCode, e.Message, e.Query)
}

func (e DatabaseException) TypeName() string {
    return "DatabaseException"
}

// Custom exception for business logic
type BusinessRuleException struct {
    Rule    string
    Value   interface{}
    Message string
}

func (e BusinessRuleException) Error() string {
    return fmt.Sprintf("BusinessRuleException: %s (Rule: %s, Value: %v)", e.Message, e.Rule, e.Value)
}

func (e BusinessRuleException) TypeName() string {
    return "BusinessRuleException"
}

// Helper functions for custom exceptions
func ThrowDatabaseError(query, message string, errorCode int) {
    Throw(DatabaseException{
        Query:     query,
        Message:   message,
        ErrorCode: errorCode,
    })
}

func ThrowBusinessRule(rule string, value interface{}, message string) {
    Throw(BusinessRuleException{
        Rule:    rule,
        Value:   value,
        Message: message,
    })
}
```

### Using Custom Exceptions

```go
// Throwing custom exceptions
Try(func() {
    ThrowDatabaseError("SELECT * FROM users", "Connection timeout", 1001)
}).Handle(
    Handler[DatabaseException](func(ex DatabaseException, full Exception) {
        fmt.Printf("Database error [%d]: %s\n", ex.ErrorCode, ex.Message)
        fmt.Printf("Failed query: %s\n", ex.Query)
    }),
    HandlerAny(func(ex Exception) {
        fmt.Printf("Unexpected error: %s\n", ex.Error())
    }),
)

// Business rule validation
Try(func() {
    userAge := 15
    if userAge < 18 {
        ThrowBusinessRule("MinimumAge", userAge, "User must be at least 18 years old")
    }
}).Handle(
    Handler[BusinessRuleException](func(ex BusinessRuleException, full Exception) {
        fmt.Printf("Business rule violated: %s\n", ex.Rule)
        fmt.Printf("Invalid value: %v\n", ex.Value)
    }),
)
```

## Available Syntaxes

### 1. Catch Function (Recommended)

```go
result := Try(func() {
    ThrowArgumentNull("username", "Required for login")
})

Catch[ArgumentNullException](result, func(ex ArgumentNullException, full Exception) {
    fmt.Printf("Null parameter: %s\n", ex.ParamName)
})

Catch[InvalidOperationException](result, func(ex InvalidOperationException, full Exception) {
    fmt.Printf("Invalid operation: %s\n", ex.Message)
})

result.Finally(func() {
    fmt.Println("Limpeza executada")
})
```

### 2. Builder Pattern

```go
Try(func() {
    ThrowNetworkError("https://api.com", "Connection failed", errors.New("timeout"))
}).When().
On[NetworkException](func(ex NetworkException, full Exception) {
    fmt.Printf("Network error: %s\n", ex.URL)
}).
On[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
    fmt.Printf("Null parameter: %s\n", ex.ParamName)
}).
Finally(func() {
    fmt.Println("Limpeza executada")
})
```

### 3. Handler Interface (More Flexible)

```go
Try(func() {
    ThrowFileError("config.txt", "File not found", os.ErrNotExist)
}).Handle(
    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
        fmt.Printf("Null parameter: %s\n", ex.ParamName)
    }),
    Handler[FileException](func(ex FileException, full Exception) {
        fmt.Printf("File error: %s\n", ex.Filename)
    }),
).Any(func(ex Exception) {
    fmt.Printf("Unhandled error: %s\n", ex.Error())
}).Finally(func() {
    fmt.Println("Cleanup executed")
})
```

### Generic Exception Handler

```go
// Catches ANY exception type
Try(func() {
    // Could throw any type of exception
    ThrowInvalidOperation("Something went wrong")
}).Any(func(ex Exception) {
    // This handler catches all exception types
    fmt.Printf("Caught: %s (Type: %s)\n", ex.Error(), ex.TypeName())
})

// Or use HandlerAny inside Handle() for mixed specific/generic handling
Try(func() {
    ThrowInvalidOperation("Something went wrong")
}).Handle(
    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
        fmt.Printf("Specific: %s\n", ex.ParamName)
    }),
    HandlerAny(func(ex Exception) {
        fmt.Printf("Generic: %s (Type: %s)\n", ex.Error(), ex.TypeName())
    }),
)
```

## Throw Functions

### Basic Throws
```go
ThrowArgumentNull("paramName", "message")
ThrowArgumentOutOfRange("paramName", value, "message")
ThrowInvalidOperation("message")
ThrowFileError("filename", "message", cause)
ThrowNetworkError("url", "message", cause)
```

### Conditional Throws
```go
ThrowIf(condition, ArgumentNullException{
    ParamName: "username",
    Message: "Username is required",
})

ThrowIfNil("paramName", value) // Automatically checks if nil
```

### Throws with Nested Exceptions
```go
ThrowWithInner(InvalidOperationException{
    Message: "Failed to initialize",
}, innerException)
```

## Nested Exceptions

```go
Try(func() {
    var innerException *Exception
    
    // Capture inner exception
    Try(func() {
        ThrowFileError("database.db", "Connection failed", errors.New("timeout"))
    }).Handle(
        Handler[FileException](func(ex FileException, full Exception) {
            innerException = &full
        }),
    )
    
    // Throw new exception with previous one as inner
    ThrowWithInner(InvalidOperationException{
        Message: "Failed to initialize application",
    }, innerException)
}).Handle(
    Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
        fmt.Printf("Main exception: %s\n", ex.Message)
        
        if full.HasInnerException() {
            fmt.Printf("Full message: %s\n", full.GetFullMessage())
            
            // Look for specific type in chain
            if fileEx := FindInnerException[FileException](&full); fileEx != nil {
                fmt.Printf("Problematic file: %s\n", fileEx.Filename)
            }
        }
    }),
)
```

## Helper Methods for Exceptions

```go
exception.HasInnerException()           // Check if has inner exception
exception.GetInnerException()           // Return inner exception
exception.GetFullMessage()              // Full message with chain
exception.GetAllExceptions()            // All exceptions in chain
FindInnerException[T](&exception)       // Find specific type in chain
```

## Performance

The system includes type cache to optimize performance in high-demand applications:

```go
// Automatic cache avoids repeated reflection
for i := 0; i < 1000; i++ {
    Try(func() {
        if i%2 == 0 {
            ThrowArgumentNull("test", "test message")
        } else {
            ThrowInvalidOperation("test operation")
        }
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            // Handler otimizado
        }),
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            // Handler otimizado
        }),
    )
}
```

## Quick Validations

```go
func ValidateUser(user *User, config map[string]string) {
    // Automatic nil validations
    ThrowIfNil("user", user)
    ThrowIfNil("config", config)
    
    // Conditional validations
    ThrowIf(user.Age < 18, ArgumentOutOfRangeException{
        ParamName: "age",
        Value: user.Age,
        Message: "Must be 18 or older",
    })
}
```

## Project Structure

```
go-exceptions/
├── goexceptions.go         # Main exception system (package)
├── package_test.go         # Package-level tests
├── doc.go                  # Package documentation
├── tests/
│   ├── goexceptions_test.go    # Core functionality tests
│   ├── exception_types_test.go # Exception type validation tests
│   └── integration_test.go     # Integration and complex scenario tests
├── cmd/
│   └── demo/
│       ├── basic_examples.go   # Basic examples
│       ├── advanced_examples.go # Advanced examples
│       ├── custom_exceptions.go # Custom exception examples
│       ├── quick_example.go    # Quick standalone example
│       ├── simple_demo.go      # Simple demo
│       └── test_panic.go       # Panic handling test
├── go.mod
└── README.md
```

## Running Examples

```bash
# Basic Examples
go run ./cmd/demo/quick_example.go     # Quick standalone example
go run ./cmd/demo/simple_demo.go       # Simple demo with core functionality
go run ./cmd/demo/basic_examples.go    # Basic usage patterns and syntax

# Advanced Examples
go run ./cmd/demo/advanced_examples.go # Advanced features and complex scenarios
go run ./cmd/demo/custom_exceptions.go # Custom exception types and usage
go run ./cmd/demo/test_panic.go        # Simple panic handling test
```

## Using as Package

```bash
# Install the package
go get github.com/bencz/go-exceptions

# In your Go code
go mod init your-project
go get github.com/bencz/go-exceptions
```

```go
package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"  // Dot import for convenience
    // OR
    // goex "github.com/bencz/go-exceptions"  // Aliased import
)

func main() {
    Try(func() {
        ThrowArgumentNull("param", "Parameter is required")
    }).Catch(func(ex ArgumentNullException) {
        fmt.Printf("Error: %s\n", ex.ParamName)
    })
}
```

## Testing

The package includes comprehensive tests to validate functionality and behavior:

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test files
go test -v ./tests/goexceptions_test.go        # Core functionality
go test -v ./tests/exception_types_test.go     # Exception types
go test -v ./tests/integration_test.go         # Integration tests
go test -v package_test.go                     # Package internals

# Run benchmarks
go test -bench=. -v

# Run tests with coverage
go test -cover ./...
```

### Test Categories

- **Core Functionality Tests** (`tests/goexceptions_test.go`): Basic throw/catch, multiple handlers, helper functions
- **Exception Type Tests** (`tests/exception_types_test.go`): Validation of all built-in exception types
- **Integration Tests** (`tests/integration_test.go`): Complex scenarios, real-world usage, performance under load
- **Package Tests** (`package_test.go`): Internal functionality, type cache, benchmarks

### Test Coverage

- Basic exception throwing and catching
- All exception types (ArgumentNull, ArgumentOutOfRange, InvalidOperation, File, Network)
- Multiple handler patterns (Handle, Catch, Any)
- Helper functions (ThrowIf, ThrowIfNil)
- Nested exceptions and exception chaining
- Finally blocks execution
- Custom exception types
- Concurrent exception handling
- Performance and caching
- Real-world integration scenarios
- Edge cases and error conditions

## Implemented Improvements

1. **reflect.TypeOf Fix**: Uses `reflect.TypeOf((*T)(nil)).Elem()` for interfaces
2. **Nested Exceptions**: Full support for inner exceptions like .NET
3. **ThrowIfNil**: Automatic validation of nil/null values
4. **Performance Cache**: Avoids repeated reflection in high-demand applications

## Compatibility

- **Go Version**: 1.24+
- **Dependencies**: Standard library only
- **Platforms**: All platforms supported by Go

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
