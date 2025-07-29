package main

import (
	"errors"
	"fmt"
	. "github.com/bencz/go-exceptions"
	"os"
)

// Basic examples demonstrating the exception system

func main() {
	fmt.Println("=== Approach 1: Catch Function (Recommended) ===")
	result := Try(func() {
		ThrowArgumentNull("username", "Required for login")
	})

	Catch[ArgumentNullException](result, func(ex ArgumentNullException, full Exception) {
		fmt.Printf("Caught ArgumentNull: %s\n", ex.ParamName)
	})

	Catch[InvalidOperationException](result, func(ex InvalidOperationException, full Exception) {
		fmt.Printf("Should not catch InvalidOperation\n")
	})

	fmt.Println("\n=== Approach 2: Builder Pattern ===")
	builder := Try(func() {
		ThrowNetworkError("https://api.com", "Connection failed", errors.New("timeout"))
	}).When()

	On[NetworkException](builder, func(ex NetworkException, full Exception) {
		fmt.Printf("Network error: %s\n", ex.URL)
	})

	On[ArgumentNullException](builder, func(ex ArgumentNullException, full Exception) {
		fmt.Printf("Should not catch ArgumentNull\n")
	})

	builder.Finally(func() {
		fmt.Println("Builder cleanup")
	})

	fmt.Println("\n=== Approach 3: Handler Interface (More Flexible) ===")
	Try(func() {
		ThrowFileError("config.txt", "File not found", os.ErrNotExist)
	}).Handle(
		Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
			fmt.Printf("Should not catch ArgumentNull\n")
		}),
		Handler[FileException](func(ex FileException, full Exception) {
			fmt.Printf("File error: %s\n", ex.Filename)
		}),
		Handler[NetworkException](func(ex NetworkException, full Exception) {
			fmt.Printf("Should not catch Network\n")
		}),
	).Finally(func() {
		fmt.Println("Handler cleanup")
	})

	fmt.Println("\n=== Practical Example: Multiple Catch Types ===")
	result2 := Try(func() {
		condition := 2 // Change to 1, 2 or 3 to test different exceptions
		switch condition {
		case 1:
			ThrowArgumentNull("param", "Parameter is required")
		case 2:
			ThrowInvalidOperation("Operation not allowed in current state")
		case 3:
			ThrowFileError("data.json", "Configuration file missing", nil)
		}
	})

	// Chain of catches - syntax closer to traditional try/catch
	Catch[ArgumentNullException](result2, func(ex ArgumentNullException, full Exception) {
		fmt.Printf("Missing parameter: %s\n", ex.ParamName)
	})

	Catch[InvalidOperationException](result2, func(ex InvalidOperationException, full Exception) {
		fmt.Printf("Invalid operation: %s\n", ex.Message)
	})

	Catch[FileException](result2, func(ex FileException, full Exception) {
		fmt.Printf("File problem: %s\n", ex.Filename)
	})

	result2.Finally(func() {
		fmt.Println("Multiple catch cleanup")
	})
}
