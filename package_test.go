package goexceptions

import (
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
        typeCache = make(map[string]bool)
        
        // Test caching behavior
        for i := 0; i < 10; i++ {
            Try(func() {
                ThrowArgumentNull("param", "test")
            }).Catch(func(ex ArgumentNullException) {
                // Handler
            })
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
            ExceptionType: ex,
            StackTrace_:   "test stack trace",
            InnerException: nil,
        }
        
        if wrapper.Error() != ex.Error() {
            t.Error("Wrapper should delegate Error() to underlying exception")
        }
        
        if wrapper.TypeName() != ex.TypeName() {
            t.Error("Wrapper should delegate TypeName() to underlying exception")
        }
        
        if wrapper.StackTrace() != "test stack trace" {
            t.Error("Wrapper should return correct stack trace")
        }
    })
}

// ============================================================================
// BENCHMARK TESTS FOR PACKAGE PERFORMANCE
// ============================================================================

func BenchmarkTypeCache(b *testing.B) {
    // Clear cache
    typeCache = make(map[string]bool)
    
    b.ResetTimer()
    
    for i := 0; i < b.N; i++ {
        Try(func() {
            ThrowArgumentNull("param", "test")
        }).Catch(func(ex ArgumentNullException) {
            // Handle
        })
    }
}

func BenchmarkWithoutCache(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Clear cache each time to simulate no caching
        typeCache = make(map[string]bool)
        
        Try(func() {
            ThrowArgumentNull("param", "test")
        }).Catch(func(ex ArgumentNullException) {
            // Handle
        })
    }
}
