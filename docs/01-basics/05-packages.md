# Packages

## Package Basics

### Package Declaration
```go
// Every Go file starts with a package declaration
package main  // Executable programs use package main

package mypackage  // Library packages use descriptive names
```

### Package main
```go
// main package and main function = entry point
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

## Importing Packages

### Basic Import
```go
import "fmt"
import "math"
import "strings"

// Grouped imports (preferred)
import (
    "fmt"
    "math"
    "strings"
)
```

### Import Aliases
```go
import (
    "fmt"
    m "math"  // Alias
    "crypto/rand"
    mrand "math/rand"  // Disambiguate
)

func main() {
    fmt.Println(m.Pi)  // Use alias
}
```

### Blank Import
```go
// Import for side effects only (init functions)
import _ "github.com/lib/pq"

// Common for database drivers
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)
```

### Dot Import (Discouraged)
```go
import . "fmt"

func main() {
    Println("No need for fmt prefix")  // Not recommended
}
```

## Package Structure

### Simple Package
```
myproject/
├── go.mod
├── main.go
└── utils/
    └── helper.go
```

```go
// utils/helper.go
package utils

func Add(a, b int) int {
    return a + b
}

// main.go
package main

import (
    "fmt"
    "myproject/utils"
)

func main() {
    result := utils.Add(3, 4)
    fmt.Println(result)
}
```

### Package with Subdirectories
```
myproject/
├── go.mod
├── main.go
└── pkg/
    ├── math/
    │   └── operations.go
    └── string/
        └── helpers.go
```

```go
// pkg/math/operations.go
package math

func Multiply(a, b int) int {
    return a * b
}

// main.go
package main

import "myproject/pkg/math"

func main() {
    result := math.Multiply(3, 4)
}
```

## Exported vs Unexported

### Naming Convention
```go
package mypackage

// Exported (public) - starts with uppercase
func PublicFunction() {}
type PublicType struct {}
const PublicConst = 10
var PublicVar = "hello"

// Unexported (private) - starts with lowercase
func privateFunction() {}
type privateType struct {}
const privateConst = 10
var privateVar = "hidden"
```

### Usage
```go
// In another package
import "mypackage"

func main() {
    mypackage.PublicFunction()  // OK
    mypackage.privateFunction() // Error: undefined
}
```

### Exported Struct Fields
```go
package models

type User struct {
    Name  string  // Exported
    Email string  // Exported
    password string  // Unexported (private)
}

func (u *User) SetPassword(pwd string) {
    u.password = pwd  // Package can access private field
}
```

## init Function

```go
package main

import "fmt"

// init runs automatically before main
func init() {
    fmt.Println("Initialization 1")
}

// Multiple init functions execute in order
func init() {
    fmt.Println("Initialization 2")
}

func main() {
    fmt.Println("Main function")
}
// Output:
// Initialization 1
// Initialization 2
// Main function
```

### Init Order
```go
// Order of initialization:
// 1. Import packages (recursively)
// 2. Initialize package-level variables
// 3. Run init functions
// 4. Run main (if main package)

package config

var DatabaseURL string

func init() {
    DatabaseURL = loadFromEnv()
}

func loadFromEnv() string {
    // Load configuration
    return "localhost:5432"
}
```

## Internal Packages

```go
// Special "internal" directory - only accessible by parent
myproject/
├── go.mod
├── main.go
└── internal/
    └── auth/
        └── token.go

// Can import from same or parent
import "myproject/internal/auth"  // OK in myproject

// Cannot import from outside
import "github.com/other/myproject/internal/auth"  // Error!
```

## Package Documentation

```go
// Package math provides basic mathematical operations.
//
// This package includes functions for arithmetic operations
// and mathematical calculations.
package math

// Add returns the sum of a and b.
//
// Example:
//   result := Add(2, 3)
//   fmt.Println(result)  // Output: 5
func Add(a, b int) int {
    return a + b
}

// Multiply returns the product of a and b.
func Multiply(a, b int) int {
    return a * b
}
```

## Go Modules

### Creating a Module
```bash
go mod init github.com/username/myproject
```

### go.mod File
```go
module github.com/username/myproject

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/lib/pq v1.10.9
)

require (
    // Indirect dependencies
    github.com/bytedance/sonic v1.9.1 // indirect
)
```

### Common Commands
```bash
# Add dependencies
go get github.com/gin-gonic/gin

# Update dependencies
go get -u github.com/gin-gonic/gin

# Remove unused dependencies
go mod tidy

# Vendor dependencies
go mod vendor

# List modules
go list -m all

# Download dependencies
go mod download
```

### Importing Local Modules
```go
// go.mod
module myapp

require mylib v0.0.0

replace mylib => ../mylib  // Local replacement
```

### Version Selection
```bash
# Specific version
go get github.com/gin-gonic/gin@v1.9.1

# Latest
go get github.com/gin-gonic/gin@latest

# Specific commit
go get github.com/gin-gonic/gin@abc123

# Specific branch
go get github.com/gin-gonic/gin@master
```

## Standard Library Packages

### Common Packages
```go
import (
    "fmt"           // Formatted I/O
    "strings"       // String manipulation
    "strconv"       // String conversions
    "time"          // Time operations
    "math"          // Mathematical functions
    "os"            // OS functionality
    "io"            // I/O primitives
    "bufio"         // Buffered I/O
    "encoding/json" // JSON encoding
    "net/http"      // HTTP client/server
    "database/sql"  // SQL database interface
    "context"       // Context for cancellation
    "sync"          // Synchronization primitives
    "errors"        // Error handling
    "log"           // Logging
)
```

### Package Organization
```go
// Use subdomain imports for clarity
import (
    // Standard library
    "fmt"
    "strings"

    // External packages
    "github.com/gin-gonic/gin"
    "github.com/lib/pq"

    // Internal packages
    "myproject/internal/auth"
    "myproject/pkg/utils"
)
```

## Package Patterns

### Single Package
```go
// Simple projects can use single package
myproject/
├── go.mod
├── main.go
├── handler.go
├── model.go
└── database.go

package main  // All files use package main
```

### Package per Feature
```go
myproject/
├── go.mod
├── main.go
├── user/
│   ├── handler.go
│   ├── model.go
│   └── repository.go
└── product/
    ├── handler.go
    ├── model.go
    └── repository.go
```

### Layered Architecture
```go
myproject/
├── go.mod
├── main.go
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── model/
└── pkg/
    └── utils/
```

## Common Patterns

### Package-level Variables
```go
package config

var (
    AppName    = "MyApp"
    Version    = "1.0.0"
    Debug      = false
)

// Usage in other files
import "myapp/config"

func main() {
    fmt.Println(config.AppName)
}
```

### Package Constructor Pattern
```go
package database

type DB struct {
    conn *sql.DB
}

// New creates a new database connection
func New(connStr string) (*DB, error) {
    conn, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }
    return &DB{conn: conn}, nil
}
```

### Package with Options
```go
package server

type Server struct {
    host string
    port int
}

type Option func(*Server)

func WithHost(host string) Option {
    return func(s *Server) {
        s.host = host
    }
}

func New(opts ...Option) *Server {
    s := &Server{host: "localhost", port: 8080}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

## Gotchas

❌ **Circular imports not allowed**
```go
// package a imports package b
// package b imports package a
// Error: import cycle not allowed
```

❌ **Cannot have multiple packages in same directory**
```go
// mydir/file1.go
package pkg1

// mydir/file2.go
package pkg2  // Error!
```

❌ **Package name vs directory name**
```go
// Directory: myutils
// File: myutils/helper.go
package utils  // Package name can differ from directory

// Import uses directory path
import "myproject/myutils"

// Use with package name
utils.Helper()
```

❌ **Exported type with unexported fields**
```go
package models

type User struct {
    name string  // Unexported - won't marshal to JSON!
}

// Fix: Export the field
type User struct {
    Name string
}
```

## Best Practices

1. **Package names should be short and clear** - lowercase, no underscores
2. **One package per directory** - keeps code organized
3. **Avoid package name stuttering** - `http.Server` not `http.HTTPServer`
4. **Use internal for private packages** - prevents external use
5. **Group imports logically** - stdlib, external, internal
6. **Document exported identifiers** - help users understand API
7. **Avoid circular dependencies** - restructure if needed
8. **Keep main package small** - delegate to other packages
9. **Use go modules** - modern dependency management
10. **Run `go mod tidy` regularly** - keep dependencies clean

## Package Documentation Examples

```go
/*
Package calculator provides basic arithmetic operations.

This package supports addition, subtraction, multiplication,
and division operations on integers and floating-point numbers.

Example usage:

	import "myproject/calculator"

	result := calculator.Add(2, 3)
	fmt.Println(result)  // Output: 5

For more complex operations, see the Advanced function.
*/
package calculator

// Add returns the sum of two integers.
func Add(a, b int) int {
    return a + b
}
```
