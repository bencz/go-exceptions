package tests

import (
    "fmt"
    "strings"
    "testing"
    . "github.com/bencz/go-exceptions"
)

// ============================================================================
// INTEGRATION AND COMPLEX SCENARIO TESTS
// ============================================================================

func TestComplexExceptionScenarios(t *testing.T) {
    t.Run("Nested Try blocks with different exception types", func(t *testing.T) {
        var outerCaught bool
        var innerCaught bool
        var finallyExecuted bool
        
        Try(func() {
            Try(func() {
                ThrowArgumentNull("innerParam", "Inner exception")
            }).Handle(
                Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                    innerCaught = true
                    // Re-throw as different exception type
                    ThrowInvalidOperation("Wrapped: " + ex.Error())
                }),
            )
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                outerCaught = true
            }),
        ).Finally(func() {
            finallyExecuted = true
        })
        
        if !innerCaught {
            t.Error("Inner exception should be caught")
        }
        if !outerCaught {
            t.Error("Outer exception should be caught")
        }
        if !finallyExecuted {
            t.Error("Finally block should execute")
        }
    })
    
    t.Run("Exception chain with multiple inner exceptions", func(t *testing.T) {
        var caught bool
        var chainLength int
        var fullMessage string
        
        // Create a chain of exceptions
        var level1 *Exception
        Try(func() {
            ThrowArgumentNull("level1", "First level error")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                level1 = &full
            }),
        )
        
        var level2 *Exception
        Try(func() {
            ThrowWithInner(InvalidOperationException{
                Message: "Second level error",
            }, level1)
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                level2 = &full
            }),
        )
        
        Try(func() {
            ThrowWithInner(ArgumentOutOfRangeException{
                ParamName: "level3",
                Value:     -1,
                Message:   "Third level error",
            }, level2)
        }).Handle(
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                caught = true
                chainLength = len(full.GetAllExceptions())
                fullMessage = full.GetFullMessage()
            }),
        )
        
        if !caught {
            t.Error("Final exception should be caught")
        }
        if chainLength != 3 {
            t.Errorf("Expected chain length 3, got %d", chainLength)
        }
        if !strings.Contains(fullMessage, "level1") {
            t.Errorf("Full message should contain all levels, got: %s", fullMessage)
        }
    })
}

func TestConcurrentExceptionHandling(t *testing.T) {
    t.Run("Concurrent exception throwing and handling", func(t *testing.T) {
        const numGoroutines = 10
        results := make(chan bool, numGoroutines)
        
        for i := 0; i < numGoroutines; i++ {
            go func(id int) {
                var caught bool
                
                Try(func() {
                    if id%2 == 0 {
                        ThrowArgumentNull(fmt.Sprintf("param%d", id), "Concurrent test")
                    } else {
                        ThrowInvalidOperation(fmt.Sprintf("Operation %d failed", id))
                    }
                }).Handle(
                    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                        caught = true
                    }),
                    Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                        caught = true
                    }),
                )
                
                results <- caught
            }(i)
        }
        
        // Collect results
        for i := 0; i < numGoroutines; i++ {
            if !<-results {
                t.Errorf("Goroutine %d should have caught exception", i)
            }
        }
    })
}

func TestPerformanceScenarios(t *testing.T) {
    t.Run("High frequency exception handling", func(t *testing.T) {
        const iterations = 1000
        var successCount int
        
        for i := 0; i < iterations; i++ {
            var handled bool
            
            Try(func() {
                if i%3 == 0 {
                    ThrowArgumentNull("param", "test")
                } else if i%3 == 1 {
                    ThrowInvalidOperation("test operation")
                } else {
                    ThrowArgumentOutOfRange("index", i, "test range")
                }
            }).Handle(
                Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                    handled = true
                }),
                Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                    handled = true
                }),
                Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                    handled = true
                }),
            )
            
            if handled {
                successCount++
            }
        }
        
        if successCount != iterations {
            t.Errorf("Expected %d successful handles, got %d", iterations, successCount)
        }
    })
}

func TestRealWorldScenarios(t *testing.T) {
    t.Run("Web API error handling simulation", func(t *testing.T) {
        type APIResponse struct {
            StatusCode int
            Message    string
        }
        
        var response APIResponse
        
        Try(func() {
            userID := ""
            if userID == "" {
                ThrowArgumentNull("userID", "User ID is required")
            }
            
            // Simulate user not found
            ThrowInvalidOperation("User not found")
            
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                response = APIResponse{StatusCode: 400, Message: "Bad Request: " + ex.ParamName}
            }),
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                response = APIResponse{StatusCode: 404, Message: "Not Found: " + ex.Message}
            }),
        ).Any(func(ex Exception) {
            response = APIResponse{StatusCode: 500, Message: "Internal Server Error"}
        })
        
        if response.StatusCode != 400 {
            t.Errorf("Expected status code 400, got %d", response.StatusCode)
        }
        if !strings.Contains(response.Message, "userID") {
            t.Errorf("Response should contain parameter name, got: %s", response.Message)
        }
    })
    
    t.Run("Database transaction simulation", func(t *testing.T) {
        var transactionRolledBack bool
        var errorLogged string
        
        Try(func() {
            // Simulate transaction start
            
            // Simulate validation error
            ThrowArgumentOutOfRange("amount", -100, "Amount cannot be negative")
            
        }).Handle(
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                transactionRolledBack = true
                errorLogged = fmt.Sprintf("Transaction failed: %s = %v", ex.ParamName, ex.Value)
            }),
        ).Finally(func() {
            // Cleanup resources
        })
        
        if !transactionRolledBack {
            t.Error("Transaction should be rolled back on validation error")
        }
        if !strings.Contains(errorLogged, "amount") {
            t.Errorf("Error log should contain parameter name, got: %s", errorLogged)
        }
    })
    
    t.Run("File processing with cleanup", func(t *testing.T) {
        var fileProcessed bool
        var cleanupExecuted bool
        var errorHandled bool
        
        Try(func() {
            // Simulate file processing
            filename := ""
            if filename == "" {
                ThrowFileError("", "Filename cannot be empty", fmt.Errorf("empty filename"))
            }
            
            fileProcessed = true
            
        }).Handle(
            Handler[FileException](func(ex FileException, full Exception) {
                errorHandled = true
            }),
        ).Finally(func() {
            cleanupExecuted = true
        })
        
        if fileProcessed {
            t.Error("File should not be processed when filename is empty")
        }
        if !errorHandled {
            t.Error("File exception should be handled")
        }
        if !cleanupExecuted {
            t.Error("Cleanup should always execute")
        }
    })
}

// ============================================================================
// EDGE CASE TESTS
// ============================================================================

func TestIntegrationEdgeCases(t *testing.T) {
    t.Run("Exception in finally block", func(t *testing.T) {
        var mainExceptionCaught bool
        var finallyExecuted bool
        
        defer func() {
            if r := recover(); r != nil {
                // Finally block exception should cause panic
                finallyExecuted = true
            }
        }()
        
        Try(func() {
            ThrowInvalidOperation("Main exception")
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                mainExceptionCaught = true
            }),
        ).Finally(func() {
            // This will cause a panic
            panic("Exception in finally block")
        })
        
        if !mainExceptionCaught {
            t.Error("Main exception should be caught")
        }
        if !finallyExecuted {
            t.Error("Finally block panic should be caught by defer")
        }
    })
    
    t.Run("Deep nesting of Try blocks", func(t *testing.T) {
        var level1Caught, level2Caught, level3Caught bool
        
        Try(func() {
            Try(func() {
                Try(func() {
                    ThrowArgumentNull("deep", "Deep nesting test")
                }).Handle(
                    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                        level3Caught = true
                        ThrowInvalidOperation("Level 3 to 2")
                    }),
                )
            }).Handle(
                Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                    level2Caught = true
                    ThrowArgumentOutOfRange("level2", 0, "Level 2 to 1")
                }),
            )
        }).Handle(
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                level1Caught = true
            }),
        )
        
        if !level1Caught || !level2Caught || !level3Caught {
            t.Error("All levels should be caught in deep nesting")
        }
    })
    
    t.Run("Mixed handler types", func(t *testing.T) {
        var specificCaught bool
        var anyCaught bool
        var finallyExecuted bool
        
        Try(func() {
            ThrowFileError("test.txt", "File error", nil)
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                t.Error("Should not catch ArgumentNullException")
            }),
            Handler[FileException](func(ex FileException, full Exception) {
                specificCaught = true
            }),
        ).Any(func(ex Exception) {
            anyCaught = true
        }).Finally(func() {
            finallyExecuted = true
        })
        
        if !specificCaught {
            t.Error("Specific FileException handler should catch")
        }
        if anyCaught {
            t.Error("Any handler should not catch when specific handler matches")
        }
        if !finallyExecuted {
            t.Error("Finally should always execute")
        }
    })
}

// ============================================================================
// STRESS TESTS
// ============================================================================

func TestStressScenarios(t *testing.T) {
    t.Run("Large exception chain", func(t *testing.T) {
        const chainLength = 50
        var currentEx *Exception
        
        // Build a long chain of exceptions
        for i := 0; i < chainLength; i++ {
            var tempEx *Exception
            
            Try(func() {
                ThrowWithInner(InvalidOperationException{
                    Message: fmt.Sprintf("Level %d", i),
                }, currentEx)
            }).Handle(
                Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                    tempEx = &full
                }),
            )
            
            currentEx = tempEx
        }
        
        // Verify the chain
        if currentEx == nil {
            t.Fatal("Final exception should not be nil")
        }
        
        allExceptions := currentEx.GetAllExceptions()
        if len(allExceptions) != chainLength {
            t.Errorf("Expected chain length %d, got %d", chainLength, len(allExceptions))
        }
        
        fullMessage := currentEx.GetFullMessage()
        if !strings.Contains(fullMessage, "Level 0") {
            t.Error("Full message should contain first level")
        }
        if !strings.Contains(fullMessage, fmt.Sprintf("Level %d", chainLength-1)) {
            t.Error("Full message should contain last level")
        }
    })
}
