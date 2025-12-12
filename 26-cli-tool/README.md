# CLI Tool - Resource Manager

A command-line tool for managing resources with subcommands, flags, help text, and proper exit codes. This documentation explains the code architecture, flag parsing, and resource management in detail.

---

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Data Structures](#data-structures)
4. [Flag Parsing](#flag-parsing)
5. [Commands & Flags Reference](#commands--flags-reference)
6. [Resource Management](#resource-management)
7. [Exit Codes](#exit-codes)
8. [Usage Examples](#usage-examples)

---

## Overview

This CLI tool (`myapp`) is a resource manager that allows you to create, list, update, and delete resources. Resources can be database types like PostgreSQL, MySQL, Redis, MongoDB, or Elasticsearch. The tool supports:

- **Subcommands**: `create`, `list`, `delete`, `update`
- **Global flags**: `--config`, `--verbose`, `--version`, `--help`
- **Persistence**: Optional JSON file storage
- **Filtering**: Query resources by field values
- **Multiple output formats**: Table and JSON

---

## Architecture

The application is structured into several components:

```
┌─────────────────────────────────────────────────────────────────┐
│                          main()                                  │
│  - Parses global flags (--config, --verbose, --version, --help) │
│  - Identifies subcommand and passes remaining args              │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                           CLI                                    │
│  - Routes subcommands to appropriate handlers                    │
│  - Methods: Run(), runCreate(), runList(), runDelete(),         │
│             runUpdate()                                          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      ResourceStore                               │
│  - In-memory resource storage with optional file persistence    │
│  - Methods: Create(), List(), Get(), Update(), Delete()         │
│  - Persistence: load(), save()                                  │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                     resources.json                               │
│  - JSON file for persistent storage (optional)                  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Data Structures

### Resource

The `Resource` struct represents a managed resource:

```go
type Resource struct {
    Name      string    `json:"name"`      // Unique identifier for the resource
    Type      string    `json:"type"`      // Resource type (postgres, mysql, redis, etc.)
    Size      int       `json:"size"`      // Size in GB
    CreatedAt time.Time `json:"created"`   // Timestamp when created
    UpdatedAt time.Time `json:"updated"`   // Timestamp when last updated
}
```

### ResourceStore

The `ResourceStore` manages resources in memory with optional file persistence:

```go
type ResourceStore struct {
    resources  map[string]*Resource  // In-memory storage (key = resource name)
    configPath string                 // Path to JSON config file for persistence
    verbose    bool                   // Enable debug logging
}
```

### CLI

The `CLI` struct wraps the store and handles command execution:

```go
type CLI struct {
    store   *ResourceStore  // Reference to the resource store
    verbose bool            // Verbose output flag
}
```

---

## Flag Parsing

### How Go's `flag` Package Works

The `flag` package provides command-line flag parsing. This application uses `flag.FlagSet` to create separate flag sets for global flags and each subcommand.

### Global Flag Parsing

The `main()` function first separates global flags from subcommand arguments:

```go
// Global flags
var configPath string
var verbose bool
var showVersion bool
var showHelp bool

// Parse global flags using a dedicated FlagSet
globalFS := flag.NewFlagSet("global", flag.ContinueOnError)
globalFS.StringVar(&configPath, "config", "", "Config file path for persistence")
globalFS.BoolVar(&verbose, "verbose", false, "Enable verbose output")
globalFS.BoolVar(&showVersion, "version", false, "Show version")
globalFS.BoolVar(&showHelp, "help", false, "Show help")
```

**How arguments are separated:**

```
myapp --config resources.json --verbose create --name "db" --type "postgres" --size 100
       └──────── Global Flags ────────┘ └────────── Subcommand + Args ───────────────┘
```

The code iterates through `os.Args[1:]` and identifies where the subcommand begins:

```go
subcommands := map[string]bool{
    "create": true,
    "list":   true,
    "delete": true,
    "update": true,
    "help":   true,
}

for i, arg := range os.Args[1:] {
    if !foundSubcommand {
        if subcommands[arg] || arg == "--version" || arg == "-v" || arg == "--help" || arg == "-h" {
            foundSubcommand = true
            subcommandArgs = os.Args[i+1:]
            break
        }
        globalArgs = append(globalArgs, arg)
    }
}
```

### Subcommand Flag Parsing

Each subcommand creates its own `FlagSet`:

```go
func (cli *CLI) runCreate(args []string) int {
    fs := flag.NewFlagSet("create", flag.ContinueOnError)
    fs.SetOutput(os.Stderr)  // Error messages go to stderr

    name := fs.String("name", "", "Name of the resource")
    resourceType := fs.String("type", "", "Type of the resource")
    size := fs.Int("size", 0, "Size of the resource in GB")
    showHelp := fs.Bool("help", false, "Show help")

    if err := fs.Parse(args); err != nil {
        return ExitUsageError
    }
    // ... rest of the logic
}
```

### Flag Types and Methods

| Method | Return Type | Usage |
|--------|-------------|-------|
| `fs.String(name, default, usage)` | `*string` | String flags like `--name "value"` |
| `fs.Int(name, default, usage)` | `*int` | Integer flags like `--size 100` |
| `fs.Bool(name, default, usage)` | `*bool` | Boolean flags like `--force` |

### Parsing Example Walkthrough

Given the command:
```bash
./myapp --config resources.json create --name "database" --type "postgres" --size 100
```

**Step 1: Raw arguments**
```
os.Args = ["./myapp", "--config", "resources.json", "create", "--name", "database", "--type", "postgres", "--size", "100"]
```

**Step 2: Split into global and subcommand args**
```
globalArgs = ["--config", "resources.json"]
subcommandArgs = ["create", "--name", "database", "--type", "postgres", "--size", "100"]
```

**Step 3: Parse global flags**
```
configPath = "resources.json"
verbose = false (default)
```

**Step 4: Route to subcommand handler**
```
cli.runCreate(["--name", "database", "--type", "postgres", "--size", "100"])
```

**Step 5: Parse subcommand flags**
```
name = "database"
resourceType = "postgres"
size = 100
```

---

## Commands & Flags Reference

### Global Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--config` | string | `""` (empty) | Path to JSON config file for persistence |
| `--verbose` | bool | `false` | Enable debug output |
| `--version` | bool | `false` | Show version and exit |
| `--help` | bool | `false` | Show help message and exit |

### `create` Command

Creates a new resource.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | ✅ Yes | Unique name for the resource |
| `--type` | string | ✅ Yes | Type: `postgres`, `mysql`, `redis`, `mongodb`, `elasticsearch` |
| `--size` | int | ✅ Yes | Size in GB (must be positive) |
| `--help` | bool | No | Show create command help |

### `list` Command

Lists all resources with optional filtering.

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | `"table"` | Output format: `table` or `json` |
| `--filter` | string | `""` | Filter by field: `field=value` (e.g., `type=postgres`) |
| `--help` | bool | `false` | Show list command help |

### `delete` Command

Deletes a resource.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | ✅ Yes | Name of the resource to delete |
| `--force` | bool | No | Skip confirmation prompt |
| `--help` | bool | No | Show delete command help |

### `update` Command

Updates an existing resource.

| Flag | Type | Required | Description |
|------|------|----------|-------------|
| `--name` | string | ✅ Yes | Name of the resource to update |
| `--size` | int | One required | New size in GB |
| `--type` | string | One required | New type |
| `--help` | bool | No | Show update command help |

---

## Resource Management

### Creating Resources

The `Create` method adds a new resource to the store:

```go
func (s *ResourceStore) Create(name, resourceType string, size int) (*Resource, error) {
    // 1. Check if resource already exists
    if _, exists := s.resources[name]; exists {
        return nil, fmt.Errorf("resource %q already exists", name)
    }

    // 2. Create resource with current timestamp
    now := time.Now()
    resource := &Resource{
        Name:      name,
        Type:      resourceType,
        Size:      size,
        CreatedAt: now,
        UpdatedAt: now,
    }

    // 3. Store in memory
    s.resources[name] = resource

    // 4. Persist to file (if config path is set)
    if err := s.save(); err != nil {
        return nil, err
    }

    return resource, nil
}
```

### Persistence (Load/Save)

Resources are persisted to a JSON file when `--config` is specified:

**Loading resources:**
```go
func (s *ResourceStore) load() error {
    data, err := os.ReadFile(s.configPath)
    if err != nil {
        if os.IsNotExist(err) {
            return nil  // Start fresh if file doesn't exist
        }
        return err
    }

    var resources []*Resource
    json.Unmarshal(data, &resources)

    for _, r := range resources {
        s.resources[r.Name] = r  // Populate in-memory map
    }
    return nil
}
```

**Saving resources:**
```go
func (s *ResourceStore) save() error {
    resources := make([]*Resource, 0, len(s.resources))
    for _, r := range s.resources {
        resources = append(resources, r)
    }

    data, _ := json.MarshalIndent(resources, "", "  ")
    return os.WriteFile(s.configPath, data, 0644)
}
```

### Filtering Resources

The `List` method supports filtering with the format `field=value`:

```go
func matchesFilter(r *Resource, filter string) bool {
    parts := strings.SplitN(filter, "=", 2)
    if len(parts) != 2 {
        return false
    }

    field := strings.ToLower(strings.TrimSpace(parts[0]))
    value := strings.ToLower(strings.TrimSpace(parts[1]))

    switch field {
    case "name":
        return strings.Contains(strings.ToLower(r.Name), value)
    case "type":
        return strings.ToLower(r.Type) == value
    default:
        return false
    }
}
```

**Filter examples:**
- `--filter "type=postgres"` - Exact match on type
- `--filter "name=db"` - Partial match on name (contains)

---

## Exit Codes

The application uses meaningful exit codes:

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | `ExitSuccess` | Operation completed successfully |
| 1 | `ExitError` | General error (e.g., JSON formatting failed) |
| 2 | `ExitUsageError` | Invalid command-line usage |
| 3 | `ExitResourceError` | Resource operation failed (not found, already exists) |

---

## Usage Examples

### Build the Application

```bash
go build -o myapp main.go
```

### Basic Commands

```bash
# Show help
./myapp --help

# Show version
./myapp --version

# Create a resource (no persistence)
./myapp create --name "database" --type "postgres" --size 100
```

### With Persistence

```bash
# Create resources with file persistence
./myapp --config resources.json create --name "database" --type "postgres" --size 100
./myapp --config resources.json create --name "cache" --type "redis" --size 10
./myapp --config resources.json create --name "search" --type "elasticsearch" --size 50

# List all resources
./myapp --config resources.json list

# Output:
# NAME       TYPE           SIZE   CREATED           UPDATED
# ----       ----           ----   -------           -------
# cache      redis          10GB   2025-12-10 12:58  2025-12-10 12:58
# database   postgres       100GB  2025-12-10 12:58  2025-12-10 12:58
# search     elasticsearch  50GB   2025-12-10 12:58  2025-12-10 12:58
```

### Filtering and Formatting

```bash
# List as JSON
./myapp --config resources.json list --format json

# Filter by type
./myapp --config resources.json list --filter "type=postgres"

# Filter by name (partial match)
./myapp --config resources.json list --filter "name=data"
```

### Updating Resources

```bash
# Update size
./myapp --config resources.json update --name "database" --size 200

# Update type
./myapp --config resources.json update --name "database" --type "mysql"

# Update both
./myapp --config resources.json update --name "database" --size 500 --type "mysql"
```

### Deleting Resources

```bash
# Delete with confirmation prompt
./myapp --config resources.json delete --name "database"
# Output: Warning: This will permanently delete the resource 'database' (postgres, 100GB)
#         Are you sure? [y/N]:

# Delete without confirmation
./myapp --config resources.json delete --name "database" --force
```

### Verbose Mode

```bash
# Enable debug output
./myapp --config resources.json --verbose create --name "db" --type "mysql" --size 50

# Output:
# [DEBUG] Loaded 3 resources from resources.json
# [DEBUG] Saved 4 resources to resources.json
# ✓ Created mysql resource: db (50GB)
```

### Example JSON Output

```json
[
  {
    "name": "database",
    "type": "postgres",
    "size": 200,
    "created": "2025-12-10T12:58:03.230842+08:00",
    "updated": "2025-12-10T12:58:51.052267+08:00"
  },
  {
    "name": "cache",
    "type": "redis",
    "size": 10,
    "created": "2025-12-10T12:58:13.531907+08:00",
    "updated": "2025-12-10T12:58:13.531907+08:00"
  }
]
```

---

## Summary

This CLI tool demonstrates several Go best practices:

1. **Subcommand pattern** - Using `flag.FlagSet` for each subcommand
2. **Separation of concerns** - CLI logic vs. storage logic
3. **Proper error handling** - Meaningful exit codes
4. **Persistence layer** - JSON file storage
5. **User experience** - Help text, confirmation prompts, verbose mode
6. **Validation** - Required flags, type validation, error messages
