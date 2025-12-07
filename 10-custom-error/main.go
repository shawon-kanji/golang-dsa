package main

import (
	"errors"
	"fmt"
	"log"
)

// ============ Custom Error Types ============

type User struct {
	name  string
	email string
	age   int
}

type ValidationError struct {
	Field string
	Value interface{}
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: field '%s' with value '%v': %v",
		e.Field, e.Value, e.Err)
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during '%s': %v", e.Operation, e.Err)
}

func (e *DatabaseError) Unwrap() error {
	return e.Err
}

type ServiceError struct {
	Service string
	Err     error
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("service '%s' error: %v", e.Service, e.Err)
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// ============ Sentinel Errors ============

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrInvalidInput = errors.New("invalid input")
)

// ============ Error Handler Chain ============

type ErrorHandler func(error) error

// Chain of error handlers
type ErrorHandlerChain struct {
	handlers []ErrorHandler
}

func NewErrorHandlerChain() *ErrorHandlerChain {
	return &ErrorHandlerChain{
		handlers: make([]ErrorHandler, 0),
	}
}

func (c *ErrorHandlerChain) Add(handler ErrorHandler) *ErrorHandlerChain {
	c.handlers = append(c.handlers, handler)
	return c // Return self for chaining
}

func (c *ErrorHandlerChain) Handle(err error) error {
	if err == nil {
		return nil
	}

	for _, handler := range c.handlers {
		err = handler(err)
		if err == nil {
			return nil // Error was fully handled
		}
	}
	return err // Return remaining error
}

// ============ Example Handlers ============

// LoggingHandler - logs the error
func LoggingHandler(err error) error {
	log.Printf("[ERROR] %v", err)
	return err // Pass through
}

// ValidationErrorHandler - handles validation errors specifically
func ValidationErrorHandler(err error) error {
	var valErr *ValidationError
	if errors.As(err, &valErr) {
		log.Printf("[VALIDATION] Field: %s, Value: %v", valErr.Field, valErr.Value)
		// Could return nil to "consume" the error, or return modified error
		return err
	}
	return err
}

// DatabaseErrorHandler - handles database errors
func DatabaseErrorHandler(err error) error {
	var dbErr *DatabaseError
	if errors.As(err, &dbErr) {
		log.Printf("[DATABASE] Operation: %s failed", dbErr.Operation)
		// Retry logic could go here
		return err
	}
	return err
}

// NotFoundHandler - handles not found errors
func NotFoundHandler(err error) error {
	if errors.Is(err, ErrNotFound) {
		log.Printf("[NOT FOUND] Resource not found, returning default")
		return nil // Consume the error - handled!
	}
	return err
}

// ============ Business Logic ============

func ValidateUser(u User) error {
	if u.name == "" {
		return &ValidationError{
			Field: "name",
			Value: u.name,
			Err:   ErrInvalidInput,
		}
	}
	if u.age < 18 {
		return &ValidationError{
			Field: "age",
			Value: u.age,
			Err:   fmt.Errorf("must be 18 or older"),
		}
	}
	if u.email == "" {
		return &ValidationError{
			Field: "email",
			Value: u.email,
			Err:   ErrInvalidInput,
		}
	}
	return nil
}

func SaveUser(u User) error {
	// Simulate database error
	return &DatabaseError{
		Operation: "INSERT",
		Err:       fmt.Errorf("connection timeout"),
	}
}

func GetUser(id int) (User, error) {
	// Simulate not found
	if id == 0 {
		return User{}, &ServiceError{
			Service: "UserService",
			Err:     ErrNotFound,
		}
	}
	return User{name: "John", email: "john@example.com", age: 25}, nil
}

// ============ Main ============

func main() {
	fmt.Println("=== Error Handler Chain Demo ===\n")

	// Create error handler chain
	chain := NewErrorHandlerChain().
		Add(LoggingHandler).
		Add(ValidationErrorHandler).
		Add(DatabaseErrorHandler).
		Add(NotFoundHandler)

	// Test 1: Validation Error
	fmt.Println("--- Test 1: Validation Error ---")
	user1 := User{name: "kanji", age: 10, email: "test@test.com"}
	if err := ValidateUser(user1); err != nil {
		remainingErr := chain.Handle(err)
		if remainingErr != nil {
			fmt.Printf("Unhandled error: %v\n", remainingErr)
		}
	}

	// Test 2: Database Error
	fmt.Println("\n--- Test 2: Database Error ---")
	user2 := User{name: "John", age: 25, email: "john@example.com"}
	if err := SaveUser(user2); err != nil {
		remainingErr := chain.Handle(err)
		if remainingErr != nil {
			fmt.Printf("Unhandled error: %v\n", remainingErr)
		}
	}

	// Test 3: Not Found Error (will be consumed by NotFoundHandler)
	fmt.Println("\n--- Test 3: Not Found Error ---")
	_, err := GetUser(0)
	if err != nil {
		remainingErr := chain.Handle(err)
		if remainingErr == nil {
			fmt.Println("Error was handled successfully!")
		}
	}

	// Test 4: Error wrapping and unwrapping
	fmt.Println("\n--- Test 4: Error Chain Unwrapping ---")
	wrappedErr := &ServiceError{
		Service: "AuthService",
		Err: &DatabaseError{
			Operation: "SELECT",
			Err:       ErrNotFound,
		},
	}

	fmt.Printf("Full error: %v\n", wrappedErr)
	fmt.Printf("Is ErrNotFound? %v\n", errors.Is(wrappedErr, ErrNotFound))

	var dbErr *DatabaseError
	if errors.As(wrappedErr, &dbErr) {
		fmt.Printf("Found DatabaseError: %v\n", dbErr)
	}
}
