package main

import (
	"fmt"
	. "github.com/bencz/go-exceptions"
)

// Quick example showing the exception system in action
func main() {
	fmt.Println("Quick Exception System Example")
	fmt.Println("==============================")

	// Example 1: Specific exception type
	fmt.Println("\n1. Specific exception handling:")
	Try(func() {
		ThrowArgumentNull("username", "Username cannot be null")
	}).Handle(
		Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
			fmt.Printf("   Caught ArgumentNull: %s\n", ex.ParamName)
		}),
	)

	// Example 2: Generic exception handler
	fmt.Println("\n2. Generic exception handling:")
	Try(func() {
		ThrowInvalidOperation("Operation not allowed")
	}).Any(func(ex Exception) {
		fmt.Printf("   Caught any exception: %s (Type: %s)\n", ex.Error(), ex.TypeName())
	})

	// Example 3: Multiple exception types
	fmt.Println("\n3. Multiple exception types:")
	Try(func() {
		age := -5
		ThrowArgumentOutOfRange("age", age, "Age must be positive")
	}).Handle(
		Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
			fmt.Printf("   This won't be called\n")
		}),
		Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
			fmt.Printf("   Range error: %s = %v\n", ex.ParamName, ex.Value)
		}),
	).Any(func(ex Exception) {
		fmt.Printf("   Fallback handler: %s\n", ex.Error())
	})

	// Example 4: With Finally
	fmt.Println("\n4. With Finally block:")
	Try(func() {
		ThrowInvalidOperation("Something failed")
	}).Handle(
		Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
			fmt.Printf("   Handled: %s\n", ex.Message)
		}),
	).Finally(func() {
		fmt.Println("   Cleanup always runs")
	})

	fmt.Println("\nDone!")
}
