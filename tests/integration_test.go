package tests

import (
    "fmt"
    "strings"
    "testing"
    . "github.com/bencz/go-exceptions"
)

// ============================================================================
// INTEGRATION TESTS - COMPLEX SCENARIOS
// ============================================================================

func TestComplexExceptionScenarios(t *testing.T) {
    t.Run("Nested function calls with exceptions", func(t *testing.T) {
        var caught bool
        var stackTrace string
        
        Try(func() {
            level1Function()
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                caught = true
                stackTrace = full.StackTrace()
            }),
        )
        
        if !caught {
            t.Error("Exception from nested function should be caught")
        }
        if !strings.Contains(stackTrace, "level1Function") {
            t.Error("Stack trace should contain function names")
        }
    })
    
    t.Run("Multiple exception types in business logic", func(t *testing.T) {
        testCases := []struct {
            name     string
            username string
            email    string
            age      int
            expected string
        }{
            {"Empty username", "", "test@example.com", 25, "ArgumentNullException"},
            {"Invalid age", "john", "test@example.com", 15, "ArgumentOutOfRangeException"},
            {"Invalid email", "john", "invalid-email", 25, "InvalidOperationException"},
            {"Valid user", "john", "john@example.com", 25, ""},
        }
        
        for _, tc := range testCases {
            t.Run(tc.name, func(t *testing.T) {
                var caughtType string
                
                Try(func() {
                    validateUserRegistration(tc.username, tc.email, tc.age)
                }).Handle(
                    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                        caughtType = "ArgumentNullException"
                    }),
                    Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                        caughtType = "ArgumentOutOfRangeException"
                    }),
                    Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                        caughtType = "InvalidOperationException"
                    }),
                )
                
                if tc.expected != caughtType {
                    t.Errorf("Expected %s, got %s", tc.expected, caughtType)
                }
            })
        }
    })
    
    t.Run("Exception chaining and propagation", func(t *testing.T) {
        var outerCaught bool
        var innerFound bool
        var chainLength int
        
        Try(func() {
            simulateServiceCall()
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                outerCaught = true
                
                // Check for inner exceptions
                if full.HasInnerException() {
                    allExceptions := full.GetAllExceptions()
                    chainLength = len(allExceptions)
                    
                    // Look for specific inner exception
                    if dbEx := FindInnerException[NetworkException](&full); dbEx != nil {
                        innerFound = true
                    }
                }
            }),
        )
        
        if !outerCaught {
            t.Error("Outer exception should be caught")
        }
        if !innerFound {
            t.Error("Inner NetworkException should be found in chain")
        }
        if chainLength < 2 {
            t.Errorf("Exception chain should have at least 2 exceptions, got %d", chainLength)
        }
    })
}

func TestConcurrentExceptionHandling(t *testing.T) {
    t.Run("Concurrent exception throwing", func(t *testing.T) {
        const numGoroutines = 10
        results := make(chan bool, numGoroutines)
        
        for i := 0; i < numGoroutines; i++ {
            go func(id int) {
                var caught bool
                
                Try(func() {
                    if id%2 == 0 {
                        ThrowArgumentNull(fmt.Sprintf("param%d", id), "Null parameter")
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
                t.Errorf("Goroutine %d failed to catch exception", i)
            }
        }
    })
}

func TestPerformanceUnderLoad(t *testing.T) {
    t.Run("High frequency exception handling", func(t *testing.T) {
        const iterations = 1000
        var successCount int
        
        for i := 0; i < iterations; i++ {
            var handled bool
            
            Try(func() {
                switch i % 3 {
                case 0:
                    ThrowArgumentNull("param", "test")
                case 1:
                    ThrowInvalidOperation("test")
                case 2:
                    ThrowArgumentOutOfRange("index", i, "test")
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

// ============================================================================
// HELPER FUNCTIONS FOR TESTS
// ============================================================================

func level1Function() {
    level2Function()
}

func level2Function() {
    level3Function()
}

func level3Function() {
    ThrowInvalidOperation("Error in deep function call")
}

func validateUserRegistration(username, email string, age int) {
    ThrowIf(username == "", ArgumentNullException{
        ParamName: "username",
        Message:   "Username is required",
    })
    
    ThrowIf(age < 18, ArgumentOutOfRangeException{
        ParamName: "age",
        Value:     age,
        Message:   "Must be 18 or older",
    })
    
    if !strings.Contains(email, "@") {
        ThrowInvalidOperation("Invalid email format")
    }
}

func simulateServiceCall() {
    var networkEx *Exception
    
    // Simulate network error
    Try(func() {
        ThrowNetworkError("Connection timeout", 408)
    }).Handle(
        Handler[NetworkException](func(ex NetworkException, full Exception) {
            networkEx = &full
        }),
    )
    
    // Simulate database error with network error as inner
    var dbEx *Exception
    Try(func() {
        ThrowWithInner(FileException{
            FileName: "database.db",
            Message:  "Database connection failed",
        }, networkEx)
    }).Handle(
        Handler[FileException](func(ex FileException, full Exception) {
            dbEx = &full
        }),
    )
    
    // Throw service error with database error as inner
    ThrowWithInner(InvalidOperationException{
        Message: "Service unavailable",
    }, dbEx)
}

// ============================================================================
// REAL-WORLD SCENARIO TESTS
// ============================================================================

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
            HandlerAny(func(ex Exception) {
                response = APIResponse{StatusCode: 500, Message: "Internal Server Error"}
            }),
        )
        
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
                ThrowFileError("", "Filename cannot be empty")
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

func TestEdgeCases(t *testing.T) {
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
            t.Error("Finally block should execute and cause panic")
        }
    })
    
    t.Run("Very deep exception nesting", func(t *testing.T) {
        var caught bool
        var nestingDepth int
        
        Try(func() {
            createDeeplyNestedExceptions(5)
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                caught = true
                allExceptions := full.GetAllExceptions()
                nestingDepth = len(allExceptions)
            }),
        )
        
        if !caught {
            t.Error("Deeply nested exception should be caught")
        }
        if nestingDepth < 5 {
            t.Errorf("Expected at least 5 nested exceptions, got %d", nestingDepth)
        }
    })
}

func createDeeplyNestedExceptions(depth int) {
    if depth <= 0 {
        ThrowInvalidOperation("Base exception")
        return
    }
    
    var innerEx *Exception
    
    Try(func() {
        createDeeplyNestedExceptions(depth - 1)
    }).Handle(
        HandlerAny(func(ex Exception) {
            innerEx = &ex
        }),
    )
    
    ThrowWithInner(InvalidOperationException{
        Message: fmt.Sprintf("Level %d exception", depth),
    }, innerEx)
}
