/*
Package goexceptions provides a comprehensive, production-ready exception handling system for Go.

This package brings structured exception handling to Go using generics and panic/recover,
offering multiple syntaxes and advanced features for robust error management.

# Overview

GoExceptions transforms Go's panic/recover mechanism into a structured, type-safe exception
handling system similar to Java, C#, and Python. It's designed for production use with
high performance, comprehensive testing (97.2% coverage), and clean APIs.

# Key Features

- Type-safe exception handling with Go 1.18+ generics
- Multiple handling syntaxes: Handler interface, Builder pattern, and Catch function
- Built-in exception types with rich context
- Custom exception type support
- Nested exceptions with inner exception chains
- Automatic stack trace capture and formatting
- Performance-optimized with reflection caching
- Helper functions for common validations
- Finally blocks for guaranteed cleanup
- Native Go panic interception and conversion
- Thread-safe operations
- Zero external dependencies

# Quick Start

	import . "github.com/bencz/go-exceptions"

	// Basic exception handling
	Try(func() {
	    ThrowArgumentNull("username", "Username cannot be empty")
	}).Handle(
	    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
	        log.Printf("Invalid parameter '%s': %s", ex.ParamName, ex.Message)
	    }),
	)

# Multiple Exception Types

	Try(func() {
	    validateUser(user)
	    processPayment(amount)
	}).Handle(
	    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
	        http.Error(w, "Missing required field: "+ex.ParamName, 400)
	    }),
	    Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
	        http.Error(w, "Invalid value for "+ex.ParamName, 400)
	    }),
	    Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
	        http.Error(w, "Operation failed: "+ex.Message, 500)
	    }),
	)

# Builder Pattern Syntax

	Try(func() {
	    connectToDatabase()
	}).When().On(func(ex NetworkException, full Exception) {
	    log.Printf("Database connection failed: %s", ex.URL)
	    useBackupConnection()
	}).Finally(func() {
	    cleanupResources()
	})

# Generic Exception Handling

	Try(func() {
	    riskyOperation()
	}).Any(func(ex Exception) {
	    log.Printf("Unexpected error [%s]: %s", ex.TypeName(), ex.Error())
	    notifyAdministrators(ex)
	})

# Finally Blocks

	Try(func() {
	    file := openFile("data.txt")
	    processFile(file)
	}).Handle(
	    Handler[FileException](func(ex FileException, full Exception) {
	        log.Printf("File error: %s", ex.Filename)
	    }),
	).Finally(func() {
	    // Always executed, even if exception occurs
	    cleanupTempFiles()
	})

# Built-in Exception Types

- ArgumentNullException - For null/nil parameter validation
- ArgumentOutOfRangeException - For parameter range validation
- InvalidOperationException - For invalid state operations
- FileException - For file system operations
- NetworkException - For network-related errors
- Exception - Base exception type

# Helper Functions

	// Validation helpers
	ThrowIfNil("config", config)
	ThrowIf(len(users) == 0, InvalidOperationException{Message: "No users found"})

	// Specific exception helpers
	ThrowArgumentNull("email", "Email address is required")
	ThrowArgumentOutOfRange("age", age, "Age must be between 0 and 150")
	ThrowInvalidOperation("Cannot delete active user")
	ThrowFileError("config.json", "Configuration file not found", err)
	ThrowNetworkError("https://api.example.com", "Connection timeout", err)

# Custom Exception Types

Create custom exceptions by implementing the ExceptionType interface:

	type DatabaseException struct {
	    Query     string
	    Database  string
	    Message   string
	    ErrorCode int
	}

	func (e DatabaseException) Error() string {
	    return fmt.Sprintf("Database error [%d] on %s: %s (Query: %s)",
	        e.ErrorCode, e.Database, e.Message, e.Query)
	}

	func (e DatabaseException) TypeName() string {
	    return "DatabaseException"
	}

	// Usage
	Try(func() {
	    Throw(DatabaseException{
	        Query:     "SELECT * FROM users",
	        Database:  "production",
	        Message:   "Connection lost",
	        ErrorCode: 2006,
	    })
	}).Handle(
	    Handler[DatabaseException](func(ex DatabaseException, full Exception) {
	        log.Printf("DB Error %d: %s", ex.ErrorCode, ex.Message)
	        reconnectToDatabase()
	    }),
	)

# Nested Exceptions

Build exception chains with inner exceptions:

	Try(func() {
	    // This will create a chain: ServiceException -> DatabaseException -> NetworkException
	    ThrowWithInner(ServiceException{Message: "User service failed"}, dbException)
	}).Handle(
	    Handler[ServiceException](func(ex ServiceException, full Exception) {
	        // Access the full exception chain
	        allExceptions := full.GetAllExceptions()
	        fullMessage := full.GetFullMessage()

	        // Find specific exception in chain
	        if networkEx := FindInnerException[NetworkException](&full); networkEx != nil {
	            log.Printf("Root cause was network issue: %s", networkEx.URL)
	        }
	    }),
	)

# Native Panic Interception

Automatically converts Go panics to exceptions:

	Try(func() {
	    // This panic will be caught and converted
	    panic("Something went wrong!")
	}).Any(func(ex Exception) {
	    log.Printf("Caught panic as exception: %s", ex.Error())
	})

# Performance

Optimized for production use:
- Reflection caching for type operations
- Minimal allocation overhead
- Efficient stack trace capture
- Benchmarked and tested

# Testing

Comprehensive test suite with 97.2% code coverage:
- Unit tests for all functionality
- Integration tests for complex scenarios
- Benchmark tests for performance validation
- Edge case and error condition testing

# Thread Safety

All operations are thread-safe and can be used in concurrent environments.

# Installation

	go get github.com/bencz/go-exceptions

# Examples

See the cmd/demo directory for comprehensive examples:
- Basic usage patterns
- Advanced exception handling
- Custom exception types
- Real-world scenarios
- Performance benchmarks

# License

MIT License - see LICENSE file for details.

# Contributing

Contributions welcome! Please see CONTRIBUTING.md for guidelines.

# Support

For issues, questions, or feature requests, please use GitHub Issues.
*/
package goexceptions
