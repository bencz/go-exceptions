package main

import (
    "fmt"
    . "github.com/bencz/go-exceptions"
)

func main() {
    fmt.Println("Go Exception System Demo")
    fmt.Println("========================")
    
    // Run simple demo first
    SimpleDemo()
    
    // Run basic examples
    BasicExamples()
    
    // Run advanced examples
    AdvancedExamples()
    
    // Run improvement examples
    ImprovementExamples()
    
    // Run custom exception examples
    CustomExceptionExamples()
    
    fmt.Println("\nTry/Catch/Throw system working perfectly!")
}
