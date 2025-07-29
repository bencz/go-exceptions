package tests

import (
    "strings"
    "testing"
    . "github.com/bencz/go-exceptions"
)

// ============================================================================
// EXCEPTION TYPE VALIDATION TESTS
// ============================================================================

func TestArgumentNullException(t *testing.T) {
    t.Run("ArgumentNullException properties", func(t *testing.T) {
        ex := ArgumentNullException{
            ParamName: "testParam",
            Message:   "Test message",
        }
        
        if ex.ParamName != "testParam" {
            t.Errorf("Expected ParamName 'testParam', got '%s'", ex.ParamName)
        }
        
        if ex.Message != "Test message" {
            t.Errorf("Expected Message 'Test message', got '%s'", ex.Message)
        }
        
        if ex.TypeName() != "ArgumentNullException" {
            t.Errorf("Expected TypeName 'ArgumentNullException', got '%s'", ex.TypeName())
        }
        
        errorMsg := ex.Error()
        if !strings.Contains(errorMsg, "testParam") {
            t.Errorf("Error message should contain param name, got: %s", errorMsg)
        }
        if !strings.Contains(errorMsg, "Test message") {
            t.Errorf("Error message should contain message, got: %s", errorMsg)
        }
    })
}

func TestArgumentOutOfRangeException(t *testing.T) {
    t.Run("ArgumentOutOfRangeException properties", func(t *testing.T) {
        ex := ArgumentOutOfRangeException{
            ParamName: "index",
            Value:     -1,
            Message:   "Index out of range",
        }
        
        if ex.ParamName != "index" {
            t.Errorf("Expected ParamName 'index', got '%s'", ex.ParamName)
        }
        
        if ex.Value != -1 {
            t.Errorf("Expected Value -1, got %v", ex.Value)
        }
        
        if ex.Message != "Index out of range" {
            t.Errorf("Expected Message 'Index out of range', got '%s'", ex.Message)
        }
        
        if ex.TypeName() != "ArgumentOutOfRangeException" {
            t.Errorf("Expected TypeName 'ArgumentOutOfRangeException', got '%s'", ex.TypeName())
        }
        
        errorMsg := ex.Error()
        if !strings.Contains(errorMsg, "index") {
            t.Errorf("Error message should contain param name, got: %s", errorMsg)
        }
        if !strings.Contains(errorMsg, "-1") {
            t.Errorf("Error message should contain value, got: %s", errorMsg)
        }
    })
}

func TestInvalidOperationException(t *testing.T) {
    t.Run("InvalidOperationException properties", func(t *testing.T) {
        ex := InvalidOperationException{
            Message: "Operation not allowed",
        }
        
        if ex.Message != "Operation not allowed" {
            t.Errorf("Expected Message 'Operation not allowed', got '%s'", ex.Message)
        }
        
        if ex.TypeName() != "InvalidOperationException" {
            t.Errorf("Expected TypeName 'InvalidOperationException', got '%s'", ex.TypeName())
        }
        
        errorMsg := ex.Error()
        if !strings.Contains(errorMsg, "Operation not allowed") {
            t.Errorf("Error message should contain message, got: %s", errorMsg)
        }
    })
}

func TestFileException(t *testing.T) {
    t.Run("FileException properties", func(t *testing.T) {
        ex := FileException{
            FileName: "test.txt",
            Message:  "File not found",
        }
        
        if ex.FileName != "test.txt" {
            t.Errorf("Expected FileName 'test.txt', got '%s'", ex.FileName)
        }
        
        if ex.Message != "File not found" {
            t.Errorf("Expected Message 'File not found', got '%s'", ex.Message)
        }
        
        if ex.TypeName() != "FileException" {
            t.Errorf("Expected TypeName 'FileException', got '%s'", ex.TypeName())
        }
        
        errorMsg := ex.Error()
        if !strings.Contains(errorMsg, "test.txt") {
            t.Errorf("Error message should contain filename, got: %s", errorMsg)
        }
        if !strings.Contains(errorMsg, "File not found") {
            t.Errorf("Error message should contain message, got: %s", errorMsg)
        }
    })
}

func TestNetworkException(t *testing.T) {
    t.Run("NetworkException properties", func(t *testing.T) {
        ex := NetworkException{
            Message:    "Connection timeout",
            StatusCode: 408,
        }
        
        if ex.Message != "Connection timeout" {
            t.Errorf("Expected Message 'Connection timeout', got '%s'", ex.Message)
        }
        
        if ex.StatusCode != 408 {
            t.Errorf("Expected StatusCode 408, got %d", ex.StatusCode)
        }
        
        if ex.TypeName() != "NetworkException" {
            t.Errorf("Expected TypeName 'NetworkException', got '%s'", ex.TypeName())
        }
        
        errorMsg := ex.Error()
        if !strings.Contains(errorMsg, "Connection timeout") {
            t.Errorf("Error message should contain message, got: %s", errorMsg)
        }
        if !strings.Contains(errorMsg, "408") {
            t.Errorf("Error message should contain status code, got: %s", errorMsg)
        }
    })
}

// ============================================================================
// HELPER FUNCTION VALIDATION TESTS
// ============================================================================

func TestThrowHelperFunctions(t *testing.T) {
    t.Run("ThrowArgumentNull creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx ArgumentNullException
        
        Try(func() {
            ThrowArgumentNull("param", "Parameter is null")
        }).Catch(func(ex ArgumentNullException) {
            caught = true
            caughtEx = ex
        })
        
        if !caught {
            t.Error("ThrowArgumentNull should throw ArgumentNullException")
        }
        if caughtEx.ParamName != "param" {
            t.Errorf("Expected ParamName 'param', got '%s'", caughtEx.ParamName)
        }
        if caughtEx.Message != "Parameter is null" {
            t.Errorf("Expected Message 'Parameter is null', got '%s'", caughtEx.Message)
        }
    })
    
    t.Run("ThrowArgumentOutOfRange creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx ArgumentOutOfRangeException
        
        Try(func() {
            ThrowArgumentOutOfRange("index", -5, "Index cannot be negative")
        }).Catch(func(ex ArgumentOutOfRangeException) {
            caught = true
            caughtEx = ex
        })
        
        if !caught {
            t.Error("ThrowArgumentOutOfRange should throw ArgumentOutOfRangeException")
        }
        if caughtEx.ParamName != "index" {
            t.Errorf("Expected ParamName 'index', got '%s'", caughtEx.ParamName)
        }
        if caughtEx.Value != -5 {
            t.Errorf("Expected Value -5, got %v", caughtEx.Value)
        }
        if caughtEx.Message != "Index cannot be negative" {
            t.Errorf("Expected Message 'Index cannot be negative', got '%s'", caughtEx.Message)
        }
    })
    
    t.Run("ThrowInvalidOperation creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx InvalidOperationException
        
        Try(func() {
            ThrowInvalidOperation("Operation not supported")
        }).Catch(func(ex InvalidOperationException) {
            caught = true
            caughtEx = ex
        })
        
        if !caught {
            t.Error("ThrowInvalidOperation should throw InvalidOperationException")
        }
        if caughtEx.Message != "Operation not supported" {
            t.Errorf("Expected Message 'Operation not supported', got '%s'", caughtEx.Message)
        }
    })
    
    t.Run("ThrowFileError creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx FileException
        
        Try(func() {
            ThrowFileError("data.txt", "Permission denied")
        }).Catch(func(ex FileException) {
            caught = true
            caughtEx = ex
        })
        
        if !caught {
            t.Error("ThrowFileError should throw FileException")
        }
        if caughtEx.FileName != "data.txt" {
            t.Errorf("Expected FileName 'data.txt', got '%s'", caughtEx.FileName)
        }
        if caughtEx.Message != "Permission denied" {
            t.Errorf("Expected Message 'Permission denied', got '%s'", caughtEx.Message)
        }
    })
    
    t.Run("ThrowNetworkError creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx NetworkException
        
        Try(func() {
            ThrowNetworkError("Connection refused", 503)
        }).Catch(func(ex NetworkException) {
            caught = true
            caughtEx = ex
        })
        
        if !caught {
            t.Error("ThrowNetworkError should throw NetworkException")
        }
        if caughtEx.Message != "Connection refused" {
            t.Errorf("Expected Message 'Connection refused', got '%s'", caughtEx.Message)
        }
        if caughtEx.StatusCode != 503 {
            t.Errorf("Expected StatusCode 503, got %d", caughtEx.StatusCode)
        }
    })
}

// ============================================================================
// EXCEPTION INTERFACE COMPLIANCE TESTS
// ============================================================================

func TestExceptionInterfaceCompliance(t *testing.T) {
    exceptions := []ExceptionType{
        ArgumentNullException{ParamName: "test", Message: "test"},
        ArgumentOutOfRangeException{ParamName: "test", Value: 0, Message: "test"},
        InvalidOperationException{Message: "test"},
        FileException{FileName: "test.txt", Message: "test"},
        NetworkException{Message: "test", StatusCode: 200},
    }
    
    for _, ex := range exceptions {
        t.Run("Exception "+ex.TypeName(), func(t *testing.T) {
            // Test Error() method
            errorMsg := ex.Error()
            if errorMsg == "" {
                t.Errorf("Error() should return non-empty string for %s", ex.TypeName())
            }
            
            // Test TypeName() method
            typeName := ex.TypeName()
            if typeName == "" {
                t.Errorf("TypeName() should return non-empty string for %s", ex.TypeName())
            }
            
            // Test that it can be thrown and caught
            var caught bool
            Try(func() {
                Throw(ex)
            }).Any(func(caughtEx Exception) {
                caught = true
                if caughtEx.TypeName() != typeName {
                    t.Errorf("Expected type %s, got %s", typeName, caughtEx.TypeName())
                }
            })
            
            if !caught {
                t.Errorf("Exception %s was not caught", typeName)
            }
        })
    }
}
