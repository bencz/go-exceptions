package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"
)

// ============================================================================
// CUSTOM EXCEPTION TYPES
// ============================================================================

// Custom exception for database operations
type DatabaseException struct {
    Query     string
    Message   string
    ErrorCode int
}

func (e DatabaseException) Error() string {
    return fmt.Sprintf("DatabaseException[%d]: %s (Query: %s)", e.ErrorCode, e.Message, e.Query)
}

func (e DatabaseException) TypeName() string {
    return "DatabaseException"
}

// Custom exception for business logic
type BusinessRuleException struct {
    Rule    string
    Value   interface{}
    Message string
}

func (e BusinessRuleException) Error() string {
    return fmt.Sprintf("BusinessRuleException: %s (Rule: %s, Value: %v)", e.Message, e.Rule, e.Value)
}

func (e BusinessRuleException) TypeName() string {
    return "BusinessRuleException"
}

// Custom exception for authentication
type AuthenticationException struct {
    Username string
    Reason   string
}

func (e AuthenticationException) Error() string {
    return fmt.Sprintf("AuthenticationException: %s (User: %s)", e.Reason, e.Username)
}

func (e AuthenticationException) TypeName() string {
    return "AuthenticationException"
}

// ============================================================================
// HELPER FUNCTIONS FOR CUSTOM EXCEPTIONS
// ============================================================================

func ThrowDatabaseError(query, message string, errorCode int) {
    Throw(DatabaseException{
        Query:     query,
        Message:   message,
        ErrorCode: errorCode,
    })
}

func ThrowBusinessRule(rule string, value interface{}, message string) {
    Throw(BusinessRuleException{
        Rule:    rule,
        Value:   value,
        Message: message,
    })
}

func ThrowAuthenticationError(username, reason string) {
    Throw(AuthenticationException{
        Username: username,
        Reason:   reason,
    })
}

// ============================================================================
// EXAMPLES USING CUSTOM EXCEPTIONS
// ============================================================================

func main() {
    fmt.Println("\n=== Custom Exception Examples ===")
    
    // Example 1: Database exception
    fmt.Println("\n1. Database Exception:")
    Try(func() {
        ThrowDatabaseError("SELECT * FROM users WHERE id = ?", "Connection timeout", 1001)
    }).Handle(
        Handler[DatabaseException](func(ex DatabaseException, full Exception) {
            fmt.Printf("   Database error [%d]: %s\n", ex.ErrorCode, ex.Message)
            fmt.Printf("   Failed query: %s\n", ex.Query)
        }),
    )
    
    // Example 2: Business rule exception
    fmt.Println("\n2. Business Rule Exception:")
    Try(func() {
        userAge := 15
        if userAge < 18 {
            ThrowBusinessRule("MinimumAge", userAge, "User must be at least 18 years old")
        }
    }).Handle(
        Handler[BusinessRuleException](func(ex BusinessRuleException, full Exception) {
            fmt.Printf("   Business rule violated: %s\n", ex.Rule)
            fmt.Printf("   Invalid value: %v\n", ex.Value)
            fmt.Printf("   Message: %s\n", ex.Message)
        }),
    )
    
    // Example 3: Authentication exception
    fmt.Println("\n3. Authentication Exception:")
    Try(func() {
        ThrowAuthenticationError("john.doe", "Invalid password")
    }).Handle(
        Handler[AuthenticationException](func(ex AuthenticationException, full Exception) {
            fmt.Printf("   Auth failed for user: %s\n", ex.Username)
            fmt.Printf("   Reason: %s\n", ex.Reason)
        }),
    )
    
    // Example 4: Multiple custom exception types
    fmt.Println("\n4. Multiple Custom Exception Types:")
    for i := 0; i < 3; i++ {
        Try(func() {
            switch i {
            case 0:
                ThrowDatabaseError("UPDATE users SET active = 1", "Deadlock detected", 1205)
            case 1:
                ThrowBusinessRule("MaxLoginAttempts", 5, "Too many login attempts")
            case 2:
                ThrowAuthenticationError("admin", "Account locked")
            }
        }).Handle(
            Handler[DatabaseException](func(ex DatabaseException, full Exception) {
                fmt.Printf("   DB Error [%d]: %s\n", ex.ErrorCode, ex.Message)
            }),
            Handler[BusinessRuleException](func(ex BusinessRuleException, full Exception) {
                fmt.Printf("   Rule Error: %s = %v\n", ex.Rule, ex.Value)
            }),
            Handler[AuthenticationException](func(ex AuthenticationException, full Exception) {
                fmt.Printf("   Auth Error: %s (%s)\n", ex.Username, ex.Reason)
            }),
            HandlerAny(func(ex Exception) {
                fmt.Printf("   Unexpected: %s\n", ex.Error())
            }),
        )
    }
    
    // Example 5: Custom exception with nested inner exception
    fmt.Println("\n5. Custom Exception with Inner Exception:")
    Try(func() {
        var innerException *Exception
        
        // Simulate a database connection failure
        Try(func() {
            ThrowDatabaseError("CONNECT", "Network unreachable", 2003)
        }).Handle(
            Handler[DatabaseException](func(ex DatabaseException, full Exception) {
                innerException = &full
            }),
        )
        
        // Throw a business rule exception with the database error as inner
        ThrowWithInner(BusinessRuleException{
            Rule:    "DatabaseAvailability",
            Value:   "required",
            Message: "Cannot process request due to database unavailability",
        }, innerException)
    }).Handle(
        Handler[BusinessRuleException](func(ex BusinessRuleException, full Exception) {
            fmt.Printf("   Business rule: %s\n", ex.Rule)
            
            if full.HasInnerException() {
                fmt.Printf("   Root cause: %s\n", full.GetFullMessage())
                
                if dbEx := FindInnerException[DatabaseException](&full); dbEx != nil {
                    fmt.Printf("   DB Error Code: %d\n", dbEx.ErrorCode)
                }
            }
        }),
    )
    
    fmt.Println("\nCustom exception examples completed!")
}
