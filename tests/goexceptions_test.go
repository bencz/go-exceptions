package tests

import (
    "strings"
    "testing"
    . "github.com/bencz/go-exceptions"
)

// ============================================================================
// CORE FUNCTIONALITY TESTS
// ============================================================================

func TestBasicTryCatch(t *testing.T) {
    t.Run("Basic ArgumentNullException handling", func(t *testing.T) {
        var caught bool
        var exceptionMessage string
        
        Try(func() {
            ThrowArgumentNull("param", "Parameter cannot be null")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
                exceptionMessage = ex.Error()
            }),
        )
        
        if !caught {
            t.Error("Exception should have been caught")
        }
        if !strings.Contains(exceptionMessage, "param") {
            t.Errorf("Exception message should contain parameter name, got: %s", exceptionMessage)
        }
    })
    
    t.Run("Basic ArgumentOutOfRangeException handling", func(t *testing.T) {
        var caught bool
        var paramName string
        var value interface{}
        
        Try(func() {
            ThrowArgumentOutOfRange("index", -1, "Index cannot be negative")
        }).Handle(
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                caught = true
                paramName = ex.ParamName
                value = ex.Value
            }),
        )
        
        if !caught {
            t.Error("Exception should have been caught")
        }
        if paramName != "index" {
            t.Errorf("Expected parameter name 'index', got '%s'", paramName)
        }
        if value != -1 {
            t.Errorf("Expected value -1, got %v", value)
        }
    })
    
    t.Run("Basic InvalidOperationException handling", func(t *testing.T) {
        var caught bool
        var message string
        
        Try(func() {
            ThrowInvalidOperation("Operation not allowed")
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                caught = true
                message = ex.Message
            }),
        )
        
        if !caught {
            t.Error("Exception should have been caught")
        }
        if message != "Operation not allowed" {
            t.Errorf("Expected message 'Operation not allowed', got '%s'", message)
        }
    })
}

func TestMultipleHandlers(t *testing.T) {
    t.Run("Multiple handlers with correct type matching", func(t *testing.T) {
        var nullCaught bool
        var rangeCaught bool
        
        Try(func() {
            ThrowArgumentNull("param", "Parameter is null")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                nullCaught = true
            }),
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                rangeCaught = true
            }),
        )
        
        if !nullCaught {
            t.Error("ArgumentNullException should have been caught")
        }
        if rangeCaught {
            t.Error("ArgumentOutOfRangeException should not have been caught")
        }
    })
}

func TestHelperFunctions(t *testing.T) {
    t.Run("ThrowIf conditional throwing", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            ThrowIf(true, ArgumentNullException{
                ParamName: "test",
                Message:   "Condition was true",
            })
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
            }),
        )
        
        if !caught {
            t.Error("ThrowIf should throw when condition is true")
        }
    })
    
    t.Run("ThrowIf no throwing when condition false", func(t *testing.T) {
        var caught bool
        var executed bool
        
        Try(func() {
            ThrowIf(false, ArgumentNullException{
                ParamName: "test",
                Message:   "Should not throw",
            })
            executed = true
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
            }),
        )
        
        if caught {
            t.Error("ThrowIf should not throw when condition is false")
        }
        if !executed {
            t.Error("Code after ThrowIf should execute when condition is false")
        }
    })
    
    t.Run("ThrowIfNil with nil pointer", func(t *testing.T) {
        var caught bool
        var paramName string
        
        Try(func() {
            var ptr *string = nil
            ThrowIfNil("pointer", ptr)
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
                paramName = ex.ParamName
            }),
        )
        
        if !caught {
            t.Error("ThrowIfNil should throw for nil pointer")
        }
        if paramName != "pointer" {
            t.Errorf("Expected parameter name 'pointer', got '%s'", paramName)
        }
    })
    
    t.Run("ThrowIfNil with nil slice", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            var slice []int = nil
            ThrowIfNil("slice", slice)
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
            }),
        )
        
        if !caught {
            t.Error("ThrowIfNil should throw for nil slice")
        }
    })
    
    t.Run("ThrowIfNil with nil map", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            var m map[string]int = nil
            ThrowIfNil("map", m)
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
            }),
        )
        
        if !caught {
            t.Error("ThrowIfNil should throw for nil map")
        }
    })
}

func TestNestedExceptions(t *testing.T) {
    t.Run("Exception with inner exception", func(t *testing.T) {
        var caught bool
        var hasInner bool
        var innerMessage string
        
        // Create inner exception first
        var innerEx *Exception
        Try(func() {
            ThrowArgumentNull("innerParam", "Inner exception")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                innerEx = &full
            }),
        )
        
        // Now throw outer exception with inner
        Try(func() {
            ThrowWithInner(InvalidOperationException{
                Message: "Outer exception",
            }, innerEx)
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                caught = true
                hasInner = full.HasInnerException()
                if hasInner {
                    innerMessage = full.GetInnerException().Error()
                }
            }),
        )
        
        if !caught {
            t.Error("Outer exception should have been caught")
        }
        if !hasInner {
            t.Error("Exception should have inner exception")
        }
        if !strings.Contains(innerMessage, "innerParam") {
            t.Errorf("Inner exception should contain 'innerParam', got: %s", innerMessage)
        }
    })
}

func TestFinallyBlock(t *testing.T) {
    t.Run("Finally block executes after exception", func(t *testing.T) {
        var finallyExecuted bool
        var exceptionCaught bool
        
        Try(func() {
            ThrowInvalidOperation("Test exception")
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                exceptionCaught = true
            }),
        ).Finally(func() {
            finallyExecuted = true
        })
        
        if !exceptionCaught {
            t.Error("Exception should have been caught")
        }
        if !finallyExecuted {
            t.Error("Finally block should have executed")
        }
    })
    
    t.Run("Finally block executes without exception", func(t *testing.T) {
        var finallyExecuted bool
        var normalExecution bool
        
        Try(func() {
            normalExecution = true
        }).Finally(func() {
            finallyExecuted = true
        })
        
        if !normalExecution {
            t.Error("Normal execution should complete")
        }
        if !finallyExecuted {
            t.Error("Finally block should execute even without exception")
        }
    })
}

// Custom exception type for testing
type CustomException struct {
    Code    int
    Message string
}

func (e CustomException) Error() string {
    return e.Message
}

func (e CustomException) TypeName() string {
    return "CustomException"
}

func TestCustomExceptions(t *testing.T) {
    t.Run("Custom exception type", func(t *testing.T) {
        
        var caught bool
        var code int
        
        Try(func() {
            Throw(CustomException{
                Code:    404,
                Message: "Custom error occurred",
            })
        }).Handle(
            Handler[CustomException](func(ex CustomException, full Exception) {
                caught = true
                code = ex.Code
            }),
        )
        
        if !caught {
            t.Error("Custom exception should have been caught")
        }
        if code != 404 {
            t.Errorf("Expected code 404, got %d", code)
        }
    })
}

func TestGenericHandlers(t *testing.T) {
    t.Run("Generic Any handler catches any exception", func(t *testing.T) {
        var caught bool
        var exceptionType string
        
        Try(func() {
            ThrowInvalidOperation("Test exception")
        }).Any(func(ex Exception) {
            caught = true
            exceptionType = ex.TypeName()
        })
        
        if !caught {
            t.Error("Any handler should catch any exception")
        }
        if exceptionType != "InvalidOperationException" {
            t.Errorf("Expected type 'InvalidOperationException', got '%s'", exceptionType)
        }
    })
    
    t.Run("Specific handler takes precedence over Any", func(t *testing.T) {
        var specificCaught bool
        var anyCaught bool
        
        Try(func() {
            ThrowArgumentNull("param", "Test")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                specificCaught = true
            }),
        ).Any(func(ex Exception) {
            anyCaught = true
        })
        
        if !specificCaught {
            t.Error("Specific handler should catch the exception")
        }
        if anyCaught {
            t.Error("Any handler should not catch when specific handler matches")
        }
    })
}

func TestCoreEdgeCases(t *testing.T) {
    t.Run("No exception thrown", func(t *testing.T) {
        var caught bool
        var normalExecution bool
        
        Try(func() {
            normalExecution = true
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
            }),
        )
        
        if caught {
            t.Error("No exception should be caught when none is thrown")
        }
        if !normalExecution {
            t.Error("Normal execution should complete")
        }
    })
    
    t.Run("Multiple Try blocks", func(t *testing.T) {
        var firstCaught bool
        var secondCaught bool
        
        Try(func() {
            ThrowArgumentNull("first", "First exception")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                firstCaught = true
            }),
        )
        
        Try(func() {
            ThrowInvalidOperation("Second exception")
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                secondCaught = true
            }),
        )
        
        if !firstCaught {
            t.Error("First exception should be caught")
        }
        if !secondCaught {
            t.Error("Second exception should be caught")
        }
    })
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkBasicTryCatch(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Try(func() {
            ThrowArgumentNull("param", "test")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                // Handle exception
            }),
        )
    }
}

func BenchmarkMultipleHandlers(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Try(func() {
            ThrowArgumentNull("param", "test")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                // Handle
            }),
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                // Handle
            }),
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                // Handle
            }),
        )
    }
}
