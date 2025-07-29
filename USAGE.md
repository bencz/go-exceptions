# Using go-exceptions as a Package

## Installation

```bash
go get github.com/bencz/go-exceptions
```

## Basic Usage

### 1. Simple Import

```go
package main

import (
    "fmt"
    goex "github.com/bencz/go-exceptions"
)

func main() {
    goex.Try(func() {
        goex.ThrowArgumentNull("param", "Parameter cannot be null")
    }).Catch(func(ex goex.ArgumentNullException) {
        fmt.Printf("Caught: %s\n", ex.ParamName)
    })
}
```

### 2. Dot Import (Recommended for convenience)

```go
package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"
)

func main() {
    Try(func() {
        ThrowArgumentNull("param", "Parameter cannot be null")
    }).Catch(func(ex ArgumentNullException) {
        fmt.Printf("Caught: %s\n", ex.ParamName)
    })
}
```

## Real-World Examples

### Web API Error Handling

```go
package main

import (
    "fmt"
    "net/http"
    . "github.com/bencz/go-exceptions"
)

func handleUser(w http.ResponseWriter, r *http.Request) {
    Try(func() {
        userID := r.URL.Query().Get("id")
        ThrowIfNil("userID", userID)
        
        if userID == "" {
            ThrowArgumentNull("userID", "User ID is required")
        }
        
        // Simulate user lookup
        if userID == "invalid" {
            ThrowInvalidOperation("User not found")
        }
        
        fmt.Fprintf(w, "User: %s", userID)
        
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            http.Error(w, fmt.Sprintf("Missing parameter: %s", ex.ParamName), http.StatusBadRequest)
        }),
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            http.Error(w, ex.Message, http.StatusNotFound)
        }),
        HandlerAny(func(ex Exception) {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }),
    )
}
```

### Database Operations

```go
package main

import (
    "database/sql"
    "fmt"
    . "github.com/bencz/go-exceptions"
)

// Custom database exception
type DatabaseException struct {
    Query     string
    Message   string
    ErrorCode int
}

func (e DatabaseException) Error() string {
    return fmt.Sprintf("DatabaseException[%d]: %s", e.ErrorCode, e.Message)
}

func (e DatabaseException) TypeName() string {
    return "DatabaseException"
}

func getUserByID(db *sql.DB, userID int) (User, error) {
    var user User
    var err error
    
    Try(func() {
        ThrowIf(userID <= 0, ArgumentOutOfRangeException{
            ParamName: "userID",
            Value:     userID,
            Message:   "User ID must be positive",
        })
        
        query := "SELECT name, email FROM users WHERE id = ?"
        row := db.QueryRow(query, userID)
        
        if scanErr := row.Scan(&user.Name, &user.Email); scanErr != nil {
            if scanErr == sql.ErrNoRows {
                ThrowInvalidOperation("User not found")
            }
            Throw(DatabaseException{
                Query:     query,
                Message:   scanErr.Error(),
                ErrorCode: 1001,
            })
        }
        
    }).Handle(
        Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
            err = fmt.Errorf("invalid user ID: %v", ex.Value)
        }),
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            err = fmt.Errorf("user not found: %s", ex.Message)
        }),
        Handler[DatabaseException](func(ex DatabaseException, full Exception) {
            err = fmt.Errorf("database error [%d]: %s", ex.ErrorCode, ex.Message)
        }),
    )
    
    return user, err
}
```

### Business Logic Validation

```go
package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"
)

type User struct {
    Name  string
    Email string
    Age   int
}

func validateUser(user *User) error {
    var validationErr error
    
    Try(func() {
        ThrowIfNil("user", user)
        
        ThrowIf(user.Name == "", ArgumentNullException{
            ParamName: "Name",
            Message:   "User name is required",
        })
        
        ThrowIf(user.Age < 18, ArgumentOutOfRangeException{
            ParamName: "Age",
            Value:     user.Age,
            Message:   "User must be at least 18 years old",
        })
        
        ThrowIf(user.Age > 120, ArgumentOutOfRangeException{
            ParamName: "Age",
            Value:     user.Age,
            Message:   "Invalid age",
        })
        
        if !isValidEmail(user.Email) {
            ThrowInvalidOperation("Invalid email format")
        }
        
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            validationErr = fmt.Errorf("validation error: %s", ex.Message)
        }),
        Handler[ArgumentOutOfRangeException](func(ex ArgumentOutOfRangeException, full Exception) {
            validationErr = fmt.Errorf("validation error: %s (value: %v)", ex.Message, ex.Value)
        }),
        Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
            validationErr = fmt.Errorf("validation error: %s", ex.Message)
        }),
    )
    
    return validationErr
}

func isValidEmail(email string) bool {
    // Simple email validation
    return len(email) > 0 && contains(email, "@")
}
```

## Best Practices

### 1. Use Dot Import for Cleaner Code

```go
import . "github.com/bencz/go-exceptions"
```

### 2. Create Custom Exception Types for Domain-Specific Errors

```go
type ValidationException struct {
    Field   string
    Value   interface{}
    Message string
}

func (e ValidationException) Error() string {
    return fmt.Sprintf("ValidationException: %s (Field: %s, Value: %v)", e.Message, e.Field, e.Value)
}

func (e ValidationException) TypeName() string {
    return "ValidationException"
}
```

### 3. Use Specific Handlers Before Generic Ones

```go
Try(func() {
    // Some operation
}).Handle(
    Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
        // Handle specific exception
    }),
    Handler[InvalidOperationException](func(ex InvalidOperationException, full Exception) {
        // Handle another specific exception
    }),
    HandlerAny(func(ex Exception) {
        // Generic fallback handler
    }),
)
```

### 4. Convert Exceptions to Standard Go Errors When Needed

```go
func someFunction() error {
    var err error
    
    Try(func() {
        // Operations that might throw
    }).Any(func(ex Exception) {
        err = fmt.Errorf("operation failed: %s", ex.Error())
    })
    
    return err
}
```

## Integration with Existing Go Code

The exception system works seamlessly with existing Go error handling:

```go
func processFile(filename string) error {
    var result error
    
    Try(func() {
        ThrowIfNil("filename", filename)
        
        file, err := os.Open(filename)
        if err != nil {
            Throw(FileException{
                FileName: filename,
                Message:  err.Error(),
            })
        }
        defer file.Close()
        
        // Process file...
        
    }).Handle(
        Handler[ArgumentNullException](func(ex ArgumentNullException, full Exception) {
            result = fmt.Errorf("invalid parameter: %s", ex.ParamName)
        }),
        Handler[FileException](func(ex FileException, full Exception) {
            result = fmt.Errorf("file error: %s", ex.Message)
        }),
    )
    
    return result
}
```
