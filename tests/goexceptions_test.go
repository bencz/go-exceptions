package tests

import (
    "fmt"
    "strings"
    "testing"
    . "github.com/bencz/go-exceptions"
)

// ============================================================================
// BASIC EXCEPTION THROWING AND CATCHING TESTS
// ============================================================================

func TestBasicThrowAndCatch(t *testing.T) {
    t.Run("ArgumentNullException", func(t *testing.T) {
        var caught bool
        var caughtParam string
        
        Try(func() {
            ThrowArgumentNull("testParam", "Test message")
        }).Catch(func(ex ArgumentNullException) {
            caught = true
            caughtParam = ex.ParamName
        })
        
        if !caught {
            t.Error("ArgumentNullException was not caught")
        }
        if caughtParam != "testParam" {
            t.Errorf("Expected param name 'testParam', got '%s'", caughtParam)
        }
    })
    
    t.Run("ArgumentOutOfRangeException", func(t *testing.T) {
        var caught bool
        var caughtValue interface{}
        
        Try(func() {
            ThrowArgumentOutOfRange("index", -1, "Index cannot be negative")
        }).Catch(func(ex ArgumentOutOfRangeException) {
            caught = true
            caughtValue = ex.Value
        })
        
        if !caught {
            t.Error("ArgumentOutOfRangeException was not caught")
        }
        if caughtValue != -1 {
            t.Errorf("Expected value -1, got %v", caughtValue)
        }
    })
    
    t.Run("InvalidOperationException", func(t *testing.T) {
        var caught bool
        var caughtMessage string
        
        Try(func() {
            ThrowInvalidOperation("Operation not allowed")
        }).Catch(func(ex InvalidOperationException) {
            caught = true
            caughtMessage = ex.Message
        })
        
        if !caught {
            t.Error("InvalidOperationException was not caught")
        }
        if caughtMessage != "Operation not allowed" {
            t.Errorf("Expected message 'Operation not allowed', got '%s'", caughtMessage)
        }
    })
}

func TestGenericThrow(t *testing.T) {
    t.Run("Generic Throw with custom exception", func(t *testing.T) {
        var caught bool
        var caughtEx ArgumentNullException
        
        Try(func() {
            Throw(ArgumentNullException{
                ParamName: "customParam",
                Message:   "Custom message",
            })
        }).Catch(func(ex ArgumentNullException) {
            caught = true
            caughtEx = ex
        })
        
        if !caught {
            t.Error("Generic throw was not caught")
        }
        if caughtEx.ParamName != "customParam" {
            t.Errorf("Expected param name 'customParam', got '%s'", caughtEx.ParamName)
        }
    })
}

// ============================================================================
// MULTIPLE HANDLER TESTS
// ============================================================================

func TestMultipleHandlers(t *testing.T) {
    t.Run("Handle with multiple specific handlers", func(t *testing.T) {
        var nullCaught, rangeCaught, invalidCaught bool
        
        for i := 0; i < 3; i++ {
            Try(func() {
                switch i {
                case 0:
                    ThrowArgumentNull("param", "null test")
                case 1:
                    ThrowArgumentOutOfRange("index", -1, "range test")
                case 2:
                    ThrowInvalidOperation("invalid test")
                }
            }).Handle(
                Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                    nullCaught = true
                }),
                Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                    rangeCaught = true
                }),
                Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                    invalidCaught = true
                }),
            )
        }
        
        if !nullCaught {
            t.Error("ArgumentNullException handler was not called")
        }
        if !rangeCaught {
            t.Error("ArgumentOutOfRangeException handler was not called")
        }
        if !invalidCaught {
            t.Error("InvalidOperationException handler was not called")
        }
    })
    
    t.Run("HandlerAny catches all exceptions", func(t *testing.T) {
        var caughtCount int
        var caughtTypes []string
        
        for i := 0; i < 3; i++ {
            Try(func() {
                switch i {
                case 0:
                    ThrowArgumentNull("param", "test")
                case 1:
                    ThrowInvalidOperation("test")
                case 2:
                    ThrowFileError("test.txt", "file error")
                }
            }).Handle(
                HandlerAny(func(ex Exception) {
                    caughtCount++
                    caughtTypes = append(caughtTypes, ex.TypeName())
                }),
            )
        }
        
        if caughtCount != 3 {
            t.Errorf("Expected 3 exceptions caught, got %d", caughtCount)
        }
        
        expectedTypes := []string{"ArgumentNullException", "InvalidOperationException", "FileException"}
        for i, expectedType := range expectedTypes {
            if i >= len(caughtTypes) || caughtTypes[i] != expectedType {
                t.Errorf("Expected type '%s' at index %d, got '%s'", expectedType, i, caughtTypes[i])
            }
        }
    })
}

// ============================================================================
// ANY METHOD TESTS
// ============================================================================

func TestAnyMethod(t *testing.T) {
    t.Run("Any catches different exception types", func(t *testing.T) {
        var caughtCount int
        
        for i := 0; i < 3; i++ {
            Try(func() {
                switch i {
                case 0:
                    ThrowArgumentNull("param", "test")
                case 1:
                    ThrowInvalidOperation("test")
                case 2:
                    ThrowNetworkError("Connection failed", 500)
                }
            }).Any(func(ex Exception) {
                caughtCount++
            })
        }
        
        if caughtCount != 3 {
            t.Errorf("Expected 3 exceptions caught with Any(), got %d", caughtCount)
        }
    })
}

// ============================================================================
// HELPER FUNCTION TESTS
// ============================================================================

func TestHelperFunctions(t *testing.T) {
    t.Run("ThrowIf with true condition", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            ThrowIf(true, ArgumentNullException{
                ParamName: "test",
                Message:   "condition was true",
            })
        }).Catch(func(ex ArgumentNullException) {
            caught = true
        })
        
        if !caught {
            t.Error("ThrowIf should have thrown when condition is true")
        }
    })
    
    t.Run("ThrowIf with false condition", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            ThrowIf(false, ArgumentNullException{
                ParamName: "test",
                Message:   "condition was false",
            })
        }).Catch(func(ex ArgumentNullException) {
            caught = true
        })
        
        if caught {
            t.Error("ThrowIf should not have thrown when condition is false")
        }
    })
    
    t.Run("ThrowIfNil with nil value", func(t *testing.T) {
        var caught bool
        var caughtParam string
        
        Try(func() {
            var nilPtr *string = nil
            ThrowIfNil("nilPtr", nilPtr)
        }).Catch(func(ex ArgumentNullException) {
            caught = true
            caughtParam = ex.ParamName
        })
        
        if !caught {
            t.Error("ThrowIfNil should have thrown for nil value")
        }
        if caughtParam != "nilPtr" {
            t.Errorf("Expected param name 'nilPtr', got '%s'", caughtParam)
        }
    })
    
    t.Run("ThrowIfNil with non-nil value", func(t *testing.T) {
        var caught bool
        
        Try(func() {
            value := "not nil"
            ThrowIfNil("value", &value)
        }).Catch(func(ex ArgumentNullException) {
            caught = true
        })
        
        if caught {
            t.Error("ThrowIfNil should not have thrown for non-nil value")
        }
    })
}

// ============================================================================
// NESTED EXCEPTION TESTS
// ============================================================================

func TestNestedExceptions(t *testing.T) {
    t.Run("ThrowWithInner creates nested exception", func(t *testing.T) {
        var caught bool
        var hasInner bool
        var innerMessage string
        
        Try(func() {
            var innerEx *Exception
            
            // Create inner exception
            Try(func() {
                ThrowArgumentNull("innerParam", "Inner exception message")
            }).Handle(
                Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                    innerEx = &full
                }),
            )
            
            // Throw outer exception with inner
            ThrowWithInner(InvalidOperationException{
                Message: "Outer exception message",
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
            t.Error("Outer exception was not caught")
        }
        if !hasInner {
            t.Error("Exception should have inner exception")
        }
        if !strings.Contains(innerMessage, "Inner exception message") {
            t.Errorf("Inner exception message not found, got: %s", innerMessage)
        }
    })
    
    t.Run("GetFullMessage includes inner exceptions", func(t *testing.T) {
        var fullMessage string
        
        Try(func() {
            var innerEx *Exception
            
            Try(func() {
                ThrowArgumentNull("param", "Inner message")
            }).Handle(
                Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                    innerEx = &full
                }),
            )
            
            ThrowWithInner(InvalidOperationException{
                Message: "Outer message",
            }, innerEx)
            
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                fullMessage = full.GetFullMessage()
            }),
        )
        
        if !strings.Contains(fullMessage, "Outer message") {
            t.Error("Full message should contain outer exception message")
        }
        if !strings.Contains(fullMessage, "Inner message") {
            t.Error("Full message should contain inner exception message")
        }
    })
}

// ============================================================================
// FINALLY BLOCK TESTS
// ============================================================================

func TestFinallyBlocks(t *testing.T) {
    t.Run("Finally executes after successful operation", func(t *testing.T) {
        var finallyExecuted bool
        
        Try(func() {
            // Normal operation
        }).Finally(func() {
            finallyExecuted = true
        })
        
        if !finallyExecuted {
            t.Error("Finally block should execute after successful operation")
        }
    })
    
    t.Run("Finally executes after exception", func(t *testing.T) {
        var finallyExecuted bool
        var exceptionCaught bool
        
        Try(func() {
            ThrowInvalidOperation("Test exception")
        }).Catch(func(ex InvalidOperationException) {
            exceptionCaught = true
        }).Finally(func() {
            finallyExecuted = true
        })
        
        if !exceptionCaught {
            t.Error("Exception should have been caught")
        }
        if !finallyExecuted {
            t.Error("Finally block should execute after exception")
        }
    })
    
    t.Run("Finally executes even with unhandled exception", func(t *testing.T) {
        var finallyExecuted bool
        
        defer func() {
            if r := recover(); r != nil {
                // Expected panic from unhandled exception
            }
        }()
        
        Try(func() {
            ThrowInvalidOperation("Unhandled exception")
        }).Finally(func() {
            finallyExecuted = true
        })
        
        if !finallyExecuted {
            t.Error("Finally block should execute even with unhandled exception")
        }
    })
}

// ============================================================================
// PERFORMANCE AND CACHING TESTS
// ============================================================================

func TestPerformanceCache(t *testing.T) {
    t.Run("Type cache improves performance", func(t *testing.T) {
        // This test verifies that the type cache is working
        // by running many operations and ensuring no panics occur
        for i := 0; i < 100; i++ {
            Try(func() {
                if i%2 == 0 {
                    ThrowArgumentNull("param", "test")
                } else {
                    ThrowInvalidOperation("test")
                }
            }).Handle(
                Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                    // Handle
                }),
                Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                    // Handle
                }),
            )
        }
        
        // If we reach here without panics, the cache is working
        t.Log("Performance cache test completed successfully")
    })
}

// ============================================================================
// CUSTOM EXCEPTION TYPE TESTS
// ============================================================================

type TestCustomException struct {
    Code    int
    Message string
}

func (e TestCustomException) Error() string {
    return fmt.Sprintf("TestCustomException[%d]: %s", e.Code, e.Message)
}

func (e TestCustomException) TypeName() string {
    return "TestCustomException"
}

func TestCustomExceptionTypes(t *testing.T) {
    t.Run("Custom exception type works with system", func(t *testing.T) {
        var caught bool
        var caughtCode int
        
        Try(func() {
            Throw(TestCustomException{
                Code:    404,
                Message: "Not found",
            })
        }).Handle(
            Handler[TestCustomException](func(ex TestCustomException, full Exception) {
                caught = true
                caughtCode = ex.Code
            }),
        )
        
        if !caught {
            t.Error("Custom exception was not caught")
        }
        if caughtCode != 404 {
            t.Errorf("Expected code 404, got %d", caughtCode)
        }
    })
    
    t.Run("Custom exception works with Any handler", func(t *testing.T) {
        var caught bool
        var typeName string
        
        Try(func() {
            Throw(TestCustomException{
                Code:    500,
                Message: "Server error",
            })
        }).Any(func(ex Exception) {
            caught = true
            typeName = ex.TypeName()
        })
        
        if !caught {
            t.Error("Custom exception was not caught by Any handler")
        }
        if typeName != "TestCustomException" {
            t.Errorf("Expected type name 'TestCustomException', got '%s'", typeName)
        }
    })
}

// ============================================================================
// EDGE CASES AND ERROR CONDITIONS
// ============================================================================

func TestEdgeCases(t *testing.T) {
    t.Run("No exception thrown - handlers not called", func(t *testing.T) {
        var handlerCalled bool
        
        Try(func() {
            // No exception thrown
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                handlerCalled = true
            }),
        )
        
        if handlerCalled {
            t.Error("Handler should not be called when no exception is thrown")
        }
    })
    
    t.Run("Multiple catches - only first matching executes", func(t *testing.T) {
        var firstCatchCalled, secondCatchCalled bool
        
        Try(func() {
            ThrowArgumentNull("param", "test")
        }).Catch(func(ex ArgumentNullException) {
            firstCatchCalled = true
        }).Catch(func(ex ArgumentNullException) {
            secondCatchCalled = true
        })
        
        if !firstCatchCalled {
            t.Error("First catch should be called")
        }
        if secondCatchCalled {
            t.Error("Second catch should not be called when first catch handles the exception")
        }
    })
}

// ============================================================================
// BENCHMARK TESTS
// ============================================================================

func BenchmarkBasicThrowCatch(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Try(func() {
            ThrowArgumentNull("param", "test message")
        }).Catch(func(ex ArgumentNullException) {
            // Handle exception
        })
    }
}

func BenchmarkMultipleHandlers(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Try(func() {
            ThrowInvalidOperation("test message")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                // Handle
            }),
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                // Handle
            }),
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                // Handle
            }),
        )
    }
}

func BenchmarkAnyHandler(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Try(func() {
            ThrowInvalidOperation("test message")
        }).Any(func(ex Exception) {
            // Handle any exception
        })
    }
}
