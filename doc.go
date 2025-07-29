// Package goexceptions provides a comprehensive try/catch/finally exception handling system for Go.
//
// This package implements a robust exception handling mechanism that mimics traditional
// exception handling found in other languages like Java, C#, and Python. It uses Go's
// panic/recover mechanism combined with generics for type-safe exception handling.
//
// Features:
//   - Type-safe exception handling with Go generics
//   - Multiple exception handling syntaxes (Catch function, Builder pattern, Handler interface)
//   - Built-in exception types (ArgumentNullException, ArgumentOutOfRangeException, etc.)
//   - Support for custom exception types
//   - Nested exception support with inner exceptions
//   - Stack trace capture and formatting
//   - Performance optimization with type caching
//   - Helper functions for common validations
//
// Basic Usage:
//
//	import . "github.com/bencz/go-exceptions"
//
//	// Simple try/catch
//	Try(func() {
//	    ThrowArgumentNull("param", "Parameter cannot be null")
//	}).Catch(func(ex ArgumentNullException) {
//	    fmt.Printf("Caught: %s\n", ex.ParamName)
//	})
//
//	// Multiple exception types
//	Try(func() {
//	    // Some operation that might throw
//	}).Handle(
//	    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
//	        fmt.Printf("Null parameter: %s\n", ex.ParamName)
//	    }),
//	    Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
//	        fmt.Printf("Invalid operation: %s\n", ex.Message)
//	    }),
//	)
//
//	// Generic exception handler
//	Try(func() {
//	    // Some operation
//	}).Any(func(ex Exception) {
//	    fmt.Printf("Caught any exception: %s\n", ex.Error())
//	})
//
// Custom Exception Types:
//
// You can create custom exception types by implementing the ExceptionType interface:
//
//	type MyCustomException struct {
//	    Code    int
//	    Message string
//	}
//
//	func (e MyCustomException) Error() string {
//	    return fmt.Sprintf("MyCustomException[%d]: %s", e.Code, e.Message)
//	}
//
//	func (e MyCustomException) TypeName() string {
//	    return "MyCustomException"
//	}
//
// For more examples and advanced usage, see the cmd/demo directory.
package goexceptions
