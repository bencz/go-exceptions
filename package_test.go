package goexceptions

import (
    "errors"
    "reflect"
    "strings"
    "testing"
)

// ============================================================================
// PACKAGE-LEVEL TESTS (Internal Testing)
// ============================================================================

func TestPackageInternals(t *testing.T) {
    t.Run("Type cache functionality", func(t *testing.T) {
        // Test that type cache is working internally
        // This test has access to package internals
        
        // Clear cache first
        typeCache = make(map[reflect.Type]bool)
        
        // Test caching behavior
        for i := 0; i < 10; i++ {
            Try(func() {
                ThrowArgumentNull("param", "test")
            }).Handle(
                Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                    // Handler
                }),
            )
        }
        
        // Verify cache has entries
        if len(typeCache) == 0 {
            t.Error("Type cache should have entries after exception handling")
        }
        
        t.Logf("Type cache has %d entries", len(typeCache))
    })
    
    t.Run("Exception wrapper creation", func(t *testing.T) {
        // Test internal exception wrapper functionality
        ex := ArgumentNullException{
            ParamName: "test",
            Message:   "test message",
        }
        
        wrapper := Exception{
            Type:       ex,
            StackTrace: []string{"test stack trace"},
            Inner:      nil,
        }
        
        if wrapper.Error() != ex.Error() {
            t.Error("Wrapper should delegate Error() to underlying exception")
        }
        
        if wrapper.TypeName() != ex.TypeName() {
            t.Error("Wrapper should delegate TypeName() to underlying exception")
        }
        
        if len(wrapper.StackTrace) == 0 {
            t.Error("Wrapper should have stack trace")
        }
    })
}

// ============================================================================
// BENCHMARK TESTS FOR PACKAGE PERFORMANCE
// ============================================================================

func BenchmarkTypeCache(b *testing.B) {
    // Clear cache
    typeCache = make(map[reflect.Type]bool)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        Try(func() {
            ThrowArgumentNull("param", "test")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                // Handle
            }),
        )
    }
}

func BenchmarkWithoutCache(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Clear cache each time to simulate no caching
        typeCache = make(map[reflect.Type]bool)
        
        Try(func() {
            ThrowArgumentNull("param", "test")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                // Handle
            }),
        )
    }
}

// ============================================================================
// COMPREHENSIVE COVERAGE TESTS
// ============================================================================

func TestAllExceptionTypesCoverage(t *testing.T) {
    t.Run("ArgumentOutOfRangeException complete coverage", func(t *testing.T) {
        ex := ArgumentOutOfRangeException{
            ParamName: "index",
            Value:     -1,
            Message:   "Index out of range",
        }
        
        if ex.Error() == "" {
            t.Error("Error() should return non-empty string")
        }
        if ex.TypeName() != "ArgumentOutOfRangeException" {
            t.Error("TypeName() should return correct type")
        }
    })
    
    t.Run("InvalidOperationException complete coverage", func(t *testing.T) {
        ex := InvalidOperationException{
            Message: "Invalid operation",
        }
        
        if ex.Error() == "" {
            t.Error("Error() should return non-empty string")
        }
        if ex.TypeName() != "InvalidOperationException" {
            t.Error("TypeName() should return correct type")
        }
    })
    
    t.Run("FileException complete coverage", func(t *testing.T) {
        ex := FileException{
            Filename: "test.txt",
            Message:  "File error",
            Cause:    errors.New("underlying error"),
        }
        
        if ex.Error() == "" {
            t.Error("Error() should return non-empty string")
        }
        if ex.TypeName() != "FileException" {
            t.Error("TypeName() should return correct type")
        }
        
        // Test without cause
        ex2 := FileException{
            Filename: "test2.txt",
            Message:  "File error without cause",
        }
        if ex2.Error() == "" {
            t.Error("Error() should work without cause")
        }
    })
    
    t.Run("NetworkException complete coverage", func(t *testing.T) {
        ex := NetworkException{
            URL:     "https://example.com",
            Message: "Network error",
            Cause:   errors.New("connection failed"),
        }
        
        if ex.Error() == "" {
            t.Error("Error() should return non-empty string")
        }
        if ex.TypeName() != "NetworkException" {
            t.Error("TypeName() should return correct type")
        }
        
        // Test without cause
        ex2 := NetworkException{
            URL:     "https://example2.com",
            Message: "Network error without cause",
        }
        if ex2.Error() == "" {
            t.Error("Error() should work without cause")
        }
    })
}

func TestAllHelperFunctionsCoverage(t *testing.T) {
    t.Run("ThrowArgumentOutOfRange", func(t *testing.T) {
        var caught bool
        var ex ArgumentOutOfRangeException
        
        Try(func() {
            ThrowArgumentOutOfRange("param", 100, "Value too high")
        }).Handle(
            Handler[ArgumentOutOfRangeException](func(e ArgumentOutOfRangeException, full Exception) {
                caught = true
                ex = e
            }),
        )
        
        if !caught {
            t.Error("ThrowArgumentOutOfRange should throw exception")
        }
        if ex.ParamName != "param" {
            t.Error("Parameter name should be preserved")
        }
        if ex.Value != 100 {
            t.Error("Value should be preserved")
        }
    })
    
    t.Run("ThrowInvalidOperation", func(t *testing.T) {
        var caught bool
        var ex InvalidOperationException
        
        Try(func() {
            ThrowInvalidOperation("Operation failed")
        }).Handle(
            Handler[InvalidOperationException](func(e InvalidOperationException, full Exception) {
                caught = true
                ex = e
            }),
        )
        
        if !caught {
            t.Error("ThrowInvalidOperation should throw exception")
        }
        if ex.Message != "Operation failed" {
            t.Error("Message should be preserved")
        }
    })
    
    t.Run("ThrowFileError", func(t *testing.T) {
        var caught bool
        var ex FileException
        
        Try(func() {
            ThrowFileError("config.txt", "File not found", errors.New("system error"))
        }).Handle(
            Handler[FileException](func(e FileException, full Exception) {
                caught = true
                ex = e
            }),
        )
        
        if !caught {
            t.Error("ThrowFileError should throw exception")
        }
        if ex.Filename != "config.txt" {
            t.Error("Filename should be preserved")
        }
        if ex.Message != "File not found" {
            t.Error("Message should be preserved")
        }
    })
    
    t.Run("ThrowNetworkError", func(t *testing.T) {
        var caught bool
        var ex NetworkException
        
        Try(func() {
            ThrowNetworkError("https://api.test.com", "Connection timeout", errors.New("timeout"))
        }).Handle(
            Handler[NetworkException](func(e NetworkException, full Exception) {
                caught = true
                ex = e
            }),
        )
        
        if !caught {
            t.Error("ThrowNetworkError should throw exception")
        }
        if ex.URL != "https://api.test.com" {
            t.Error("URL should be preserved")
        }
        if ex.Message != "Connection timeout" {
            t.Error("Message should be preserved")
        }
    })
    
    t.Run("ThrowIf with true condition", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            ThrowIf(true, InvalidOperationException{Message: "Condition was true"})
        }).Handle(
            Handler[InvalidOperationException](func(e InvalidOperationException, full Exception) {
                caught = true
            }),
        )
        
        if !caught {
            t.Error("ThrowIf should throw when condition is true")
        }
    })
    
    t.Run("ThrowIf with false condition", func(t *testing.T) {
        var caught bool
        var executed bool
        
        Try(func() {
            ThrowIf(false, InvalidOperationException{Message: "Should not throw"})
            executed = true
        }).Handle(
            Handler[InvalidOperationException](func(e InvalidOperationException, full Exception) {
                caught = true
            }),
        )
        
        if caught {
            t.Error("ThrowIf should not throw when condition is false")
        }
        if !executed {
            t.Error("Code should continue executing when condition is false")
        }
    })
    
    t.Run("ThrowIfNil with various nil types", func(t *testing.T) {
        // Test nil pointer
        var caught1 bool
        Try(func() {
            var ptr *int = nil
            ThrowIfNil("pointer", ptr)
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                caught1 = true
            }),
        )
        if !caught1 {
            t.Error("ThrowIfNil should throw for nil pointer")
        }
        
        // Test non-nil value (should not throw)
        var caught7 bool
        var executed bool
        Try(func() {
            value := "not nil"
            ThrowIfNil("value", value)
            executed = true
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                caught7 = true
            }),
        )
        if caught7 {
            t.Error("ThrowIfNil should not throw for non-nil value")
        }
        if !executed {
            t.Error("Code should continue executing for non-nil value")
        }
    })
    
    t.Run("ThrowWithInner", func(t *testing.T) {
        var innerEx *Exception
        var outerCaught bool
        var hasInner bool
        
        // Create inner exception
        Try(func() {
            ThrowArgumentNull("inner", "Inner exception")
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                innerEx = &full
            }),
        )
        
        // Throw outer exception with inner
        Try(func() {
            ThrowWithInner(InvalidOperationException{
                Message: "Outer exception",
            }, innerEx)
        }).Handle(
            Handler[InvalidOperationException](func(e InvalidOperationException, full Exception) {
                outerCaught = true
                hasInner = full.HasInnerException()
            }),
        )
        
        if !outerCaught {
            t.Error("Outer exception should be caught")
        }
        if !hasInner {
            t.Error("Exception should have inner exception")
        }
    })
}

func TestBuilderPatternCoverage(t *testing.T) {
    t.Run("When().On() pattern", func(t *testing.T) {
        var caught bool
        
        builder := Try(func() {
            ThrowArgumentNull("param", "Test message")
        }).When()
        
        On(builder, func(ex ArgumentNullException, full Exception) {
            caught = true
        }).End()
        
        if !caught {
            t.Error("When().On() should catch exception")
        }
    })
    
    t.Run("When().Any() pattern", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            ThrowInvalidOperation("Test operation")
        }).When().Any(func(ex Exception) {
            caught = true
        }).End()
        
        if !caught {
            t.Error("When().Any() should catch any exception")
        }
    })
    
    t.Run("When().Finally() pattern", func(t *testing.T) {
        var caught bool
        var finallyExecuted bool
        
        builder := Try(func() {
            ThrowArgumentNull("param", "Test")
        }).When()
        
        On(builder, func(ex ArgumentNullException, full Exception) {
            caught = true
        }).Finally(func() {
            finallyExecuted = true
        })
        
        if !caught {
            t.Error("Exception should be caught")
        }
        if !finallyExecuted {
            t.Error("Finally should be executed")
        }
    })
}

func TestCatchFunctionCoverage(t *testing.T) {
    t.Run("Catch function usage", func(t *testing.T) {
        var caught bool
        
        result := Try(func() {
            ThrowArgumentNull("param", "Test")
        })
        
        Catch(result, func(ex ArgumentNullException, full Exception) {
            caught = true
        })
        
        if !caught {
            t.Error("Catch function should handle exception")
        }
    })
}

func TestExceptionMethodsCoverage(t *testing.T) {
    t.Run("HasException and GetException", func(t *testing.T) {
        result := Try(func() {
            ThrowInvalidOperation("Test exception")
        })
        
        if !result.HasException() {
            t.Error("HasException should return true when exception occurred")
        }
        
        ex := result.GetException()
        if ex == nil {
            t.Error("GetException should return exception")
        }
        if ex.TypeName() != "InvalidOperationException" {
            t.Error("Exception should have correct type")
        }
    })
    
    t.Run("Rethrow", func(t *testing.T) {
        var caught bool
        
        defer func() {
            if r := recover(); r != nil {
                caught = true
            }
        }()
        
        result := Try(func() {
            ThrowInvalidOperation("Test exception")
        })
        
        result.Rethrow()
        
        if !caught {
            t.Error("Rethrow should cause panic when exception not handled")
        }
    })
    
    t.Run("Finally method", func(t *testing.T) {
        var finallyExecuted bool
        
        Try(func() {
            ThrowInvalidOperation("Test")
        }).Handle(
            Handler[InvalidOperationException](func(e InvalidOperationException, full Exception) {
                // Handle exception
            }),
        ).Finally(func() {
            finallyExecuted = true
        })
        
        if !finallyExecuted {
            t.Error("Finally method should execute")
        }
    })
    
    t.Run("Any method", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            ThrowInvalidOperation("Test")
        }).Any(func(ex Exception) {
            caught = true
        })
        
        if !caught {
            t.Error("Any method should catch exception")
        }
    })
}

func TestInnerExceptionMethodsCoverage(t *testing.T) {
    t.Run("Complete inner exception functionality", func(t *testing.T) {
        var level1Ex *Exception
        var level2Ex *Exception
        var level3Ex *Exception
        
        // Create level 1 exception
        Try(func() {
            ThrowArgumentNull("level1", "Level 1 error")
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                level1Ex = &full
            }),
        )
        
        // Create level 2 exception with level 1 as inner
        Try(func() {
            ThrowWithInner(InvalidOperationException{
                Message: "Level 2 error",
            }, level1Ex)
        }).Handle(
            Handler[InvalidOperationException](func(e InvalidOperationException, full Exception) {
                level2Ex = &full
            }),
        )
        
        // Create level 3 exception with level 2 as inner
        Try(func() {
            ThrowWithInner(ArgumentOutOfRangeException{
                ParamName: "level3",
                Value:     -1,
                Message:   "Level 3 error",
            }, level2Ex)
        }).Handle(
            Handler[ArgumentOutOfRangeException](func(e ArgumentOutOfRangeException, full Exception) {
                level3Ex = &full
            }),
        )
        
        // Test HasInnerException
        if !level3Ex.HasInnerException() {
            t.Error("Level 3 should have inner exception")
        }
        
        // Test GetInnerException
        inner := level3Ex.GetInnerException()
        if inner == nil {
            t.Error("GetInnerException should return inner exception")
        }
        if inner.TypeName() != "InvalidOperationException" {
            t.Error("Inner exception should be InvalidOperationException")
        }
        
        // Test GetFullMessage
        fullMessage := level3Ex.GetFullMessage()
        if !strings.Contains(fullMessage, "Level 1") {
            t.Error("Full message should contain all levels")
        }
        if !strings.Contains(fullMessage, "Level 2") {
            t.Error("Full message should contain all levels")
        }
        if !strings.Contains(fullMessage, "Level 3") {
            t.Error("Full message should contain all levels")
        }
        
        // Test GetAllExceptions
        allExceptions := level3Ex.GetAllExceptions()
        if len(allExceptions) != 3 {
            t.Errorf("Should have 3 exceptions in chain, got %d", len(allExceptions))
        }
        
        // Test FindInnerException
        foundArgNull := FindInnerException[ArgumentNullException](level3Ex)
        if foundArgNull == nil {
            t.Error("Should find ArgumentNullException in chain")
        }
        if foundArgNull.ParamName != "level1" {
            t.Error("Found exception should have correct parameter name")
        }
        
        foundInvalid := FindInnerException[InvalidOperationException](level3Ex)
        if foundInvalid == nil {
            t.Error("Should find InvalidOperationException in chain")
        }
        if foundInvalid.Message != "Level 2 error" {
            t.Error("Found exception should have correct message")
        }
        
        // Test FindInnerException for non-existent type
        foundFile := FindInnerException[FileException](level3Ex)
        if foundFile != nil {
            t.Error("Should not find FileException in chain")
        }
    })
}

func TestHandlerAnyCoverage(t *testing.T) {
    t.Run("HandlerAny function", func(t *testing.T) {
        var caught bool
        var exceptionType string
        
        Try(func() {
            ThrowInvalidOperation("Test exception")
        }).Handle(
            HandlerAny(func(ex Exception) {
                caught = true
                exceptionType = ex.TypeName()
            }),
        )
        
        if !caught {
            t.Error("HandlerAny should catch any exception")
        }
        if exceptionType != "InvalidOperationException" {
            t.Error("Should capture correct exception type")
        }
    })
}

func TestEdgeCasesCoverage(t *testing.T) {
    t.Run("ThrowIfNil with more nil types", func(t *testing.T) {
        // Test nil slice
        var caught1 bool
        Try(func() {
            var slice []string = nil
            ThrowIfNil("slice", slice)
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                caught1 = true
            }),
        )
        if !caught1 {
            t.Error("ThrowIfNil should throw for nil slice")
        }
        
        // Test nil map
        var caught2 bool
        Try(func() {
            var m map[string]int = nil
            ThrowIfNil("map", m)
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                caught2 = true
            }),
        )
        if !caught2 {
            t.Error("ThrowIfNil should throw for nil map")
        }
        
        // Test nil channel
        var caught3 bool
        Try(func() {
            var ch chan int = nil
            ThrowIfNil("channel", ch)
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                caught3 = true
            }),
        )
        if !caught3 {
            t.Error("ThrowIfNil should throw for nil channel")
        }
        
        // Test nil function
        var caught4 bool
        Try(func() {
            var fn func() = nil
            ThrowIfNil("function", fn)
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                caught4 = true
            }),
        )
        if !caught4 {
            t.Error("ThrowIfNil should throw for nil function")
        }
        
        // Test nil interface
        var caught5 bool
        Try(func() {
            var iface interface{} = nil
            ThrowIfNil("interface", iface)
        }).Handle(
            Handler[ArgumentNullException](func(e ArgumentNullException, full Exception) {
                caught5 = true
            }),
        )
        if !caught5 {
            t.Error("ThrowIfNil should throw for nil interface")
        }
    })
    
    t.Run("GetException when no exception", func(t *testing.T) {
        result := Try(func() {
            // No exception thrown
        })
        
        if result.HasException() {
            t.Error("Should not have exception")
        }
        
        ex := result.GetException()
        if ex != nil {
            t.Error("GetException should return nil when no exception")
        }
    })
    
    t.Run("Catch with no exception", func(t *testing.T) {
        var caught bool
        
        result := Try(func() {
            // No exception thrown
        })
        
        Catch(result, func(ex ArgumentNullException, full Exception) {
            caught = true
        })
        
        if caught {
            t.Error("Catch should not execute when no exception")
        }
    })
    
    t.Run("On with no exception", func(t *testing.T) {
        var caught bool
        
        builder := Try(func() {
            // No exception thrown
        }).When()
        
        On(builder, func(ex ArgumentNullException, full Exception) {
            caught = true
        }).End()
        
        if caught {
            t.Error("On should not execute when no exception")
        }
    })
    
    t.Run("Try with different panic types", func(t *testing.T) {
        // Test panic with error
        var caught1 bool
        Try(func() {
            panic(errors.New("error panic"))
        }).Any(func(ex Exception) {
            caught1 = true
        })
        if !caught1 {
            t.Error("Should catch error panic")
        }
        
        // Test panic with integer
        var caught2 bool
        Try(func() {
            panic(42)
        }).Any(func(ex Exception) {
            caught2 = true
        })
        if !caught2 {
            t.Error("Should catch integer panic")
        }
        
        // Test panic with struct
        var caught3 bool
        Try(func() {
            panic(struct{ msg string }{msg: "struct panic"})
        }).Any(func(ex Exception) {
            caught3 = true
        })
        if !caught3 {
            t.Error("Should catch struct panic")
        }
    })
}
