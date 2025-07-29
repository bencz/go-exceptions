package tests

import (
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
        
        expectedError := "ArgumentNullException: Parameter 'testParam' cannot be null. Test message"
        if ex.Error() != expectedError {
            t.Errorf("Expected Error '%s', got '%s'", expectedError, ex.Error())
        }
    })
}

func TestArgumentOutOfRangeException(t *testing.T) {
    t.Run("ArgumentOutOfRangeException properties", func(t *testing.T) {
        ex := ArgumentOutOfRangeException{
            ParamName: "age",
            Value:     -5,
            Message:   "Age cannot be negative",
        }
        
        if ex.ParamName != "age" {
            t.Errorf("Expected ParamName 'age', got '%s'", ex.ParamName)
        }
        
        if ex.Value != -5 {
            t.Errorf("Expected Value -5, got %v", ex.Value)
        }
        
        if ex.Message != "Age cannot be negative" {
            t.Errorf("Expected Message 'Age cannot be negative', got '%s'", ex.Message)
        }
        
        if ex.TypeName() != "ArgumentOutOfRangeException" {
            t.Errorf("Expected TypeName 'ArgumentOutOfRangeException', got '%s'", ex.TypeName())
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
        
        expectedError := "InvalidOperationException: Operation not allowed"
        if ex.Error() != expectedError {
            t.Errorf("Expected Error '%s', got '%s'", expectedError, ex.Error())
        }
    })
}

func TestFileException(t *testing.T) {
    t.Run("FileException properties", func(t *testing.T) {
        ex := FileException{
            Filename: "test.txt",
            Message:  "File not found",
        }
        
        if ex.Filename != "test.txt" {
            t.Errorf("Expected Filename 'test.txt', got '%s'", ex.Filename)
        }
        
        if ex.Message != "File not found" {
            t.Errorf("Expected Message 'File not found', got '%s'", ex.Message)
        }
        
        if ex.TypeName() != "FileException" {
            t.Errorf("Expected TypeName 'FileException', got '%s'", ex.TypeName())
        }
    })
}

func TestNetworkException(t *testing.T) {
    t.Run("NetworkException properties", func(t *testing.T) {
        ex := NetworkException{
            URL:     "https://api.example.com",
            Message: "Connection timeout",
        }
        
        if ex.URL != "https://api.example.com" {
            t.Errorf("Expected URL 'https://api.example.com', got '%s'", ex.URL)
        }
        
        if ex.Message != "Connection timeout" {
            t.Errorf("Expected Message 'Connection timeout', got '%s'", ex.Message)
        }
        
        if ex.TypeName() != "NetworkException" {
            t.Errorf("Expected TypeName 'NetworkException', got '%s'", ex.TypeName())
        }
    })
}

// ============================================================================
// HELPER FUNCTION TESTS
// ============================================================================

func TestThrowHelperFunctions(t *testing.T) {
    t.Run("ThrowArgumentNull creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx ArgumentNullException
        
        Try(func() {
            ThrowArgumentNull("param", "Parameter is null")
        }).Handle(
            Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
                caught = true
                caughtEx = ex
            }),
        )
        
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
            ThrowArgumentOutOfRange("index", 10, "Index out of bounds")
        }).Handle(
            Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
                caught = true
                caughtEx = ex
            }),
        )
        
        if !caught {
            t.Error("ThrowArgumentOutOfRange should throw ArgumentOutOfRangeException")
        }
        if caughtEx.ParamName != "index" {
            t.Errorf("Expected ParamName 'index', got '%s'", caughtEx.ParamName)
        }
        if caughtEx.Value != 10 {
            t.Errorf("Expected Value 10, got %v", caughtEx.Value)
        }
    })
    
    t.Run("ThrowInvalidOperation creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx InvalidOperationException
        
        Try(func() {
            ThrowInvalidOperation("Invalid state")
        }).Handle(
            Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
                caught = true
                caughtEx = ex
            }),
        )
        
        if !caught {
            t.Error("ThrowInvalidOperation should throw InvalidOperationException")
        }
        if caughtEx.Message != "Invalid state" {
            t.Errorf("Expected Message 'Invalid state', got '%s'", caughtEx.Message)
        }
    })
    
    t.Run("ThrowFileError creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx FileException
        
        Try(func() {
            ThrowFileError("config.txt", "File not accessible", nil)
        }).Handle(
            Handler[FileException](func(ex FileException, full Exception) {
                caught = true
                caughtEx = ex
            }),
        )
        
        if !caught {
            t.Error("ThrowFileError should throw FileException")
        }
        if caughtEx.Filename != "config.txt" {
            t.Errorf("Expected Filename 'config.txt', got '%s'", caughtEx.Filename)
        }
        if caughtEx.Message != "File not accessible" {
            t.Errorf("Expected Message 'File not accessible', got '%s'", caughtEx.Message)
        }
    })
    
    t.Run("ThrowNetworkError creates correct exception", func(t *testing.T) {
        var caught bool
        var caughtEx NetworkException
        
        Try(func() {
            ThrowNetworkError("https://api.test.com", "Connection failed", nil)
        }).Handle(
            Handler[NetworkException](func(ex NetworkException, full Exception) {
                caught = true
                caughtEx = ex
            }),
        )
        
        if !caught {
            t.Error("ThrowNetworkError should throw NetworkException")
        }
        if caughtEx.URL != "https://api.test.com" {
            t.Errorf("Expected URL 'https://api.test.com', got '%s'", caughtEx.URL)
        }
        if caughtEx.Message != "Connection failed" {
            t.Errorf("Expected Message 'Connection failed', got '%s'", caughtEx.Message)
        }
    })
}
