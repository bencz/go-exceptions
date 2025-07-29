package goexceptions

import (
    "reflect"
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
