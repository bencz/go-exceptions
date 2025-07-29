package goexceptions

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// ExceptionType represents an exception type
type ExceptionType interface {
	TypeName() string
	error
}

// Specific exception types with uniform interface
type ArgumentNullException struct {
	ParamName string
	Message   string
}

func (e ArgumentNullException) Error() string {
	return fmt.Sprintf("ArgumentNullException: Parameter '%s' cannot be null. %s", e.ParamName, e.Message)
}

func (e ArgumentNullException) TypeName() string {
	return "ArgumentNullException"
}

// ArgumentOutOfRangeException ( comment to force new release... )
type ArgumentOutOfRangeException struct {
	ParamName string
	Value     interface{}
	Message   string
}

func (e ArgumentOutOfRangeException) Error() string {
	return fmt.Sprintf("ArgumentOutOfRangeException: Parameter '%s' with value '%v' is out of range. %s", e.ParamName, e.Value, e.Message)
}

func (e ArgumentOutOfRangeException) TypeName() string {
	return "ArgumentOutOfRangeException"
}

type InvalidOperationException struct {
	Message string
}

func (e InvalidOperationException) Error() string {
	return fmt.Sprintf("InvalidOperationException: %s", e.Message)
}

func (e InvalidOperationException) TypeName() string {
	return "InvalidOperationException"
}

type FileException struct {
	Filename string
	Message  string
	Cause    error
}

func (e FileException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("FileException: %s (File: %s, Cause: %v)", e.Message, e.Filename, e.Cause)
	}
	return fmt.Sprintf("FileException: %s (File: %s)", e.Message, e.Filename)
}

func (e FileException) TypeName() string {
	return "FileException"
}

type NetworkException struct {
	URL     string
	Message string
	Cause   error
}

func (e NetworkException) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("NetworkException: %s (URL: %s, Cause: %v)", e.Message, e.URL, e.Cause)
	}
	return fmt.Sprintf("NetworkException: %s (URL: %s)", e.Message, e.URL)
}

func (e NetworkException) TypeName() string {
	return "NetworkException"
}

// Exception is the main wrapper
type Exception struct {
	Type       ExceptionType
	StackTrace []string
	Data       map[string]interface{}
	Inner      *Exception // support for nested exceptions
}

func (e Exception) Error() string {
	return e.Type.Error()
}

func (e Exception) TypeName() string {
	return e.Type.TypeName()
}

// Generic throw
func Throw[T ExceptionType](exception T) {
	panic(Exception{
		Type:       exception,
		StackTrace: getStackTrace(),
		Data:       make(map[string]interface{}),
	})
}

// Helper throw functions
func ThrowArgumentNull(paramName, message string) {
	Throw(ArgumentNullException{ParamName: paramName, Message: message})
}

func ThrowArgumentOutOfRange[T any](paramName string, value T, message string) {
	Throw(ArgumentOutOfRangeException{ParamName: paramName, Value: value, Message: message})
}

func ThrowInvalidOperation(message string) {
	Throw(InvalidOperationException{Message: message})
}

func ThrowFileError(filename, message string, cause error) {
	Throw(FileException{Filename: filename, Message: message, Cause: cause})
}

func ThrowNetworkError(url, message string, cause error) {
	Throw(NetworkException{URL: url, Message: message, Cause: cause})
}

func ThrowIf[T ExceptionType](condition bool, exception T) {
	if condition {
		Throw(exception)
	}
}

// ThrowIfNil throws ArgumentNullException if value is nil
func ThrowIfNil(paramName string, value any) {
	if value == nil {
		ThrowArgumentNull(paramName, "")
		return
	}

	// Check if it's a nilable type and if it's nil
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func:
		if v.IsNil() {
			ThrowArgumentNull(paramName, "")
		}
	}
}

// ThrowWithInner throws an exception with an inner exception
func ThrowWithInner[T ExceptionType](exception T, inner *Exception) {
	panic(Exception{
		Type:       exception,
		StackTrace: getStackTrace(),
		Data:       make(map[string]interface{}),
		Inner:      inner,
	})
}

func getStackTrace() []string {
	var traces []string
	for i := 3; i < 15; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fn := runtime.FuncForPC(pc)
		if fn == nil {
			continue
		}

		funcName := fn.Name()
		if strings.Contains(funcName, "runtime.") || strings.Contains(funcName, "panic") {
			continue
		}

		traces = append(traces, fmt.Sprintf("%s:%d %s", file, line, funcName))
	}
	return traces
}

// ============================================================================
// EXPANDABLE SOLUTION: Using Type Constraints and Reflection
// ============================================================================

// TryResult with expandable system
type TryResult struct {
	exception *Exception
	handled   bool
}

// Try executes a block that can throw exceptions
func Try(tryBlock func()) *TryResult {
	var exception *Exception

	// Internal function to ensure defer is executed correctly
	func() {
		defer func() {
			if r := recover(); r != nil {
				switch e := r.(type) {
				case Exception:
					exception = &e
				case ExceptionType:
					exception = &Exception{
						Type:       e,
						StackTrace: getStackTrace(),
						Data:       make(map[string]interface{}),
					}
				case error:
					exception = &Exception{
						Type:       InvalidOperationException{Message: e.Error()},
						StackTrace: getStackTrace(),
						Data:       make(map[string]interface{}),
					}
				default:
					exception = &Exception{
						Type:       InvalidOperationException{Message: fmt.Sprintf("%v", r)},
						StackTrace: getStackTrace(),
						Data:       make(map[string]interface{}),
					}
				}
			}
		}()

		tryBlock()
	}()

	return &TryResult{exception: exception}
}

// ============================================================================
// PERFORMANCE: Type cache to avoid repeated reflection
// ============================================================================

var typeCache = make(map[reflect.Type]bool)
var typeCacheMutex sync.RWMutex

func getTypeOf[T any]() reflect.Type {
	// Use reflect.TypeOf((*T)(nil)).Elem() to capture the correct type
	// even when T is an interface
	return reflect.TypeOf((*T)(nil)).Elem()
}

func isTypeMatch[T any](actualType reflect.Type) bool {
	expectedType := getTypeOf[T]()

	// Cache lookup for performance
	typeCacheMutex.RLock()
	cacheKey := expectedType
	if cached, exists := typeCache[cacheKey]; exists {
		typeCacheMutex.RUnlock()
		return cached && actualType == expectedType
	}
	typeCacheMutex.RUnlock()

	// Calculate and store in cache
	match := actualType == expectedType
	typeCacheMutex.Lock()
	typeCache[cacheKey] = match
	typeCacheMutex.Unlock()

	return match
}

// ============================================================================
// APPROACH 1: Catch with Type Parameter using external function
// ============================================================================

func Catch[T ExceptionType](tr *TryResult, handler func(T, Exception)) *TryResult {
	if tr == nil || tr.exception == nil || tr.handled {
		return tr
	}

	// Check if exception type is compatible using cache
	actualType := reflect.TypeOf(tr.exception.Type)

	if isTypeMatch[T](actualType) {
		exceptionValue := tr.exception.Type.(T)
		handler(exceptionValue, *tr.exception)
		tr.handled = true
	}

	return tr
}

// ============================================================================
// APPROACH 2: Cleaner Builder Pattern
// ============================================================================

type CatchBuilder struct {
	result *TryResult
}

func (tr *TryResult) When() *CatchBuilder {
	return &CatchBuilder{result: tr}
}

func On[T ExceptionType](cb *CatchBuilder, handler func(T, Exception)) *CatchBuilder {
	if cb.result == nil || cb.result.exception == nil || cb.result.handled {
		return cb
	}

	actualType := reflect.TypeOf(cb.result.exception.Type)

	if isTypeMatch[T](actualType) {
		exceptionValue := cb.result.exception.Type.(T)
		handler(exceptionValue, *cb.result.exception)
		cb.result.handled = true
	}

	return cb
}

func (cb *CatchBuilder) Any(handler func(Exception)) *CatchBuilder {
	if cb.result != nil && cb.result.exception != nil && !cb.result.handled {
		handler(*cb.result.exception)
		cb.result.handled = true
	}
	return cb
}

func (cb *CatchBuilder) Finally(cleanup func()) *TryResult {
	if cb.result != nil {
		cleanup()
	}
	return cb.result
}

func (cb *CatchBuilder) End() *TryResult {
	return cb.result
}

// ============================================================================
// APPROACH 3: Using interfaces and smarter type switching
// ============================================================================

type ExceptionHandler interface {
	Handle(ex Exception) bool
}

// TypedHandler for any type
type TypedHandler[T ExceptionType] struct {
	handler func(T, Exception)
}

func (th *TypedHandler[T]) Handle(ex Exception) bool {
	actualType := reflect.TypeOf(ex.Type)

	if isTypeMatch[T](actualType) {
		typedEx := ex.Type.(T)
		th.handler(typedEx, ex)
		return true
	}
	return false
}

func Handler[T ExceptionType](handler func(T, Exception)) ExceptionHandler {
	return &TypedHandler[T]{handler: handler}
}

// GenericHandler for catching any Exception type
type GenericHandler struct {
	handler func(Exception)
}

func (gh *GenericHandler) Handle(ex Exception) bool {
	gh.handler(ex)
	return true
}

// HandlerAny creates a generic handler that catches any exception
func HandlerAny(handler func(Exception)) ExceptionHandler {
	return &GenericHandler{handler: handler}
}

func (tr *TryResult) Handle(handlers ...ExceptionHandler) *TryResult {
	if tr == nil || tr.exception == nil || tr.handled {
		return tr
	}

	for _, handler := range handlers {
		if handler.Handle(*tr.exception) {
			tr.handled = true
			break
		}
	}

	return tr
}

func (tr *TryResult) Finally(cleanup func()) *TryResult {
	if tr != nil {
		cleanup()
	}
	return tr
}

func (tr *TryResult) Any(handler func(Exception)) *TryResult {
	if tr != nil && tr.exception != nil && !tr.handled {
		handler(*tr.exception)
		tr.handled = true
	}
	return tr
}

// HasException checks if there was an exception
func (tr *TryResult) HasException() bool {
	return tr != nil && tr.exception != nil
}

// GetException returns the exception if any
func (tr *TryResult) GetException() *Exception {
	if tr == nil {
		return nil
	}
	return tr.exception
}

// Rethrow re-throws the exception if it wasn't handled
func (tr *TryResult) Rethrow() {
	if tr != nil && tr.exception != nil && !tr.handled {
		panic(*tr.exception)
	}
}

// ============================================================================
// HELPER METHODS FOR NESTED EXCEPTIONS
// ============================================================================

// HasInnerException checks if the exception has an inner exception
func (e *Exception) HasInnerException() bool {
	return e.Inner != nil
}

// GetInnerException returns the inner exception
func (e *Exception) GetInnerException() *Exception {
	return e.Inner
}

// GetFullMessage returns the full message including inner exceptions
func (e *Exception) GetFullMessage() string {
	message := e.Error()
	if e.Inner != nil {
		message += " --> " + e.Inner.GetFullMessage()
	}
	return message
}

// GetAllExceptions returns all exceptions in the chain
func (e *Exception) GetAllExceptions() []*Exception {
	var exceptions []*Exception
	current := e
	for current != nil {
		exceptions = append(exceptions, current)
		current = current.Inner
	}
	return exceptions
}

// FindInnerException finds the first inner exception of the specified type
func FindInnerException[T ExceptionType](e *Exception) *T {
	current := e
	for current != nil {
		if isTypeMatch[T](reflect.TypeOf(current.Type)) {
			if typed, ok := current.Type.(T); ok {
				return &typed
			}
		}
		current = current.Inner
	}
	return nil
}
