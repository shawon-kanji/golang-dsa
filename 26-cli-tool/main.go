package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"
)

// =====================================================
// CLI Tool with Flags - Resource Manager
// =====================================================
// A command-line tool for managing resources with subcommands,
// flags, help text, and proper exit codes.

const (
	Version = "1.0.0"

	// Exit codes
	ExitSuccess       = 0
	ExitError         = 1
	ExitUsageError    = 2
	ExitResourceError = 3
)

// Resource represents a managed resource
type Resource struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created"`
	UpdatedAt time.Time `json:"updated"`
}

// ResourceStore manages resources in memory with optional file persistence
type ResourceStore struct {
	resources  map[string]*Resource
	configPath string
	verbose    bool
}

// NewResourceStore creates a new resource store
func NewResourceStore(configPath string, verbose bool) *ResourceStore {
	store := &ResourceStore{
		resources:  make(map[string]*Resource),
		configPath: configPath,
		verbose:    verbose,
	}

	// Try to load existing resources from config file
	if configPath != "" {
		store.load()
	}

	return store
}

// load reads resources from the config file
func (s *ResourceStore) load() error {
	if s.configPath == "" {
		return nil
	}

	data, err := os.ReadFile(s.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			if s.verbose {
				fmt.Printf("[DEBUG] Config file not found, starting fresh\n")
			}
			return nil
		}
		return fmt.Errorf("failed to read config: %w", err)
	}

	var resources []*Resource
	if err := json.Unmarshal(data, &resources); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	for _, r := range resources {
		s.resources[r.Name] = r
	}

	if s.verbose {
		fmt.Printf("[DEBUG] Loaded %d resources from %s\n", len(resources), s.configPath)
	}

	return nil
}

// save writes resources to the config file
func (s *ResourceStore) save() error {
	if s.configPath == "" {
		return nil
	}

	resources := make([]*Resource, 0, len(s.resources))
	for _, r := range s.resources {
		resources = append(resources, r)
	}

	data, err := json.MarshalIndent(resources, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize resources: %w", err)
	}

	if err := os.WriteFile(s.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	if s.verbose {
		fmt.Printf("[DEBUG] Saved %d resources to %s\n", len(resources), s.configPath)
	}

	return nil
}

// Create adds a new resource
func (s *ResourceStore) Create(name, resourceType string, size int) (*Resource, error) {
	if _, exists := s.resources[name]; exists {
		return nil, fmt.Errorf("resource %q already exists", name)
	}

	now := time.Now()
	resource := &Resource{
		Name:      name,
		Type:      resourceType,
		Size:      size,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.resources[name] = resource

	if err := s.save(); err != nil {
		return nil, err
	}

	return resource, nil
}

// List returns all resources, optionally filtered
func (s *ResourceStore) List(filter string) []*Resource {
	resources := make([]*Resource, 0, len(s.resources))

	for _, r := range s.resources {
		if filter == "" || matchesFilter(r, filter) {
			resources = append(resources, r)
		}
	}

	// Sort by name for consistent output
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Name < resources[j].Name
	})

	return resources
}

// matchesFilter checks if a resource matches the filter string
// Filter format: "field=value" (e.g., "type=postgres")
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

// Get retrieves a resource by name
func (s *ResourceStore) Get(name string) (*Resource, bool) {
	r, exists := s.resources[name]
	return r, exists
}

// Delete removes a resource
func (s *ResourceStore) Delete(name string) error {
	if _, exists := s.resources[name]; !exists {
		return fmt.Errorf("resource %q not found", name)
	}

	delete(s.resources, name)

	return s.save()
}

// Update modifies an existing resource
func (s *ResourceStore) Update(name string, newSize *int, newType *string) (*Resource, error) {
	r, exists := s.resources[name]
	if !exists {
		return nil, fmt.Errorf("resource %q not found", name)
	}

	if newSize != nil {
		r.Size = *newSize
	}
	if newType != nil && *newType != "" {
		r.Type = *newType
	}
	r.UpdatedAt = time.Now()

	if err := s.save(); err != nil {
		return nil, err
	}

	return r, nil
}

// =====================================================
// CLI Application
// =====================================================

type CLI struct {
	store   *ResourceStore
	verbose bool
}

// printHelp displays the main help message
func printHelp() {
	help := `Usage: myapp [command] [flags]

A tool for managing resources

Commands:
  create    Create a new resource
  list      List all resources
  delete    Delete a resource
  update    Update a resource

Global Flags:
  --config string    Config file path for persistence (optional)
  --verbose          Enable verbose output
  --version          Show version
  --help             Show this help message

Run 'myapp [command] --help' for more information on a command.

Examples:
  myapp create --name "database" --type "postgres" --size 100
  myapp list --format json
  myapp list --filter "type=postgres"
  myapp update --name "database" --size 200
  myapp delete --name "database" --force
`
	fmt.Print(help)
}

// printCreateHelp displays help for the create command
func printCreateHelp() {
	help := `Usage: myapp create [flags]

Create a new resource

Flags:
  --name string     Name of the resource (required)
  --type string     Type of the resource (required)
                    Supported types: postgres, mysql, redis, mongodb, elasticsearch
  --size int        Size of the resource in GB (required)
  --help            Show this help message

Examples:
  myapp create --name "database" --type "postgres" --size 100
  myapp create --name "cache" --type "redis" --size 10
`
	fmt.Print(help)
}

// printListHelp displays help for the list command
func printListHelp() {
	help := `Usage: myapp list [flags]

List all resources

Flags:
  --format string   Output format: table, json (default: table)
  --filter string   Filter resources (format: field=value)
                    Supported fields: name, type
  --help            Show this help message

Examples:
  myapp list
  myapp list --format json
  myapp list --filter "type=postgres"
  myapp list --format json --filter "name=db"
`
	fmt.Print(help)
}

// printDeleteHelp displays help for the delete command
func printDeleteHelp() {
	help := `Usage: myapp delete [flags]

Delete a resource

Flags:
  --name string     Name of the resource to delete (required)
  --force           Skip confirmation prompt
  --help            Show this help message

Examples:
  myapp delete --name "database"
  myapp delete --name "database" --force
`
	fmt.Print(help)
}

// printUpdateHelp displays help for the update command
func printUpdateHelp() {
	help := `Usage: myapp update [flags]

Update an existing resource

Flags:
  --name string     Name of the resource to update (required)
  --size int        New size in GB (optional)
  --type string     New type (optional)
  --help            Show this help message

Examples:
  myapp update --name "database" --size 200
  myapp update --name "database" --type "mysql"
  myapp update --name "database" --size 500 --type "mysql"
`
	fmt.Print(help)
}

// runCreate handles the create subcommand
func (cli *CLI) runCreate(args []string) int {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	name := fs.String("name", "", "Name of the resource")
	resourceType := fs.String("type", "", "Type of the resource")
	size := fs.Int("size", 0, "Size of the resource in GB")
	showHelp := fs.Bool("help", false, "Show help")

	if err := fs.Parse(args); err != nil {
		return ExitUsageError
	}

	if *showHelp {
		printCreateHelp()
		return ExitSuccess
	}

	// Validate required flags
	var errors []string
	if *name == "" {
		errors = append(errors, "  --name is required")
	}
	if *resourceType == "" {
		errors = append(errors, "  --type is required")
	}
	if *size <= 0 {
		errors = append(errors, "  --size must be a positive number")
	}

	// Validate resource type
	validTypes := map[string]bool{
		"postgres":      true,
		"mysql":         true,
		"redis":         true,
		"mongodb":       true,
		"elasticsearch": true,
	}
	if *resourceType != "" && !validTypes[strings.ToLower(*resourceType)] {
		errors = append(errors, "  --type must be one of: postgres, mysql, redis, mongodb, elasticsearch")
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Error: invalid arguments\n%s\n\nRun 'myapp create --help' for usage.\n",
			strings.Join(errors, "\n"))
		return ExitUsageError
	}

	// Create the resource
	resource, err := cli.store.Create(*name, strings.ToLower(*resourceType), *size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitResourceError
	}

	fmt.Printf("✓ Created %s resource: %s (%dGB)\n", resource.Type, resource.Name, resource.Size)
	return ExitSuccess
}

// runList handles the list subcommand
func (cli *CLI) runList(args []string) int {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	format := fs.String("format", "table", "Output format: table, json")
	filter := fs.String("filter", "", "Filter resources (format: field=value)")
	showHelp := fs.Bool("help", false, "Show help")

	if err := fs.Parse(args); err != nil {
		return ExitUsageError
	}

	if *showHelp {
		printListHelp()
		return ExitSuccess
	}

	// Validate format
	if *format != "table" && *format != "json" {
		fmt.Fprintf(os.Stderr, "Error: --format must be 'table' or 'json'\n")
		return ExitUsageError
	}

	resources := cli.store.List(*filter)

	if len(resources) == 0 {
		if *filter != "" {
			fmt.Println("No resources found matching filter")
		} else {
			fmt.Println("No resources found")
		}
		return ExitSuccess
	}

	switch *format {
	case "json":
		data, err := json.MarshalIndent(resources, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to format output: %v\n", err)
			return ExitError
		}
		fmt.Println(string(data))

	case "table":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tSIZE\tCREATED\tUPDATED")
		fmt.Fprintln(w, "----\t----\t----\t-------\t-------")
		for _, r := range resources {
			fmt.Fprintf(w, "%s\t%s\t%dGB\t%s\t%s\n",
				r.Name,
				r.Type,
				r.Size,
				r.CreatedAt.Format("2006-01-02 15:04"),
				r.UpdatedAt.Format("2006-01-02 15:04"),
			)
		}
		w.Flush()
	}

	return ExitSuccess
}

// runDelete handles the delete subcommand
func (cli *CLI) runDelete(args []string) int {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	name := fs.String("name", "", "Name of the resource to delete")
	force := fs.Bool("force", false, "Skip confirmation prompt")
	showHelp := fs.Bool("help", false, "Show help")

	if err := fs.Parse(args); err != nil {
		return ExitUsageError
	}

	if *showHelp {
		printDeleteHelp()
		return ExitSuccess
	}

	if *name == "" {
		fmt.Fprintf(os.Stderr, "Error: --name is required\n\nRun 'myapp delete --help' for usage.\n")
		return ExitUsageError
	}

	// Check if resource exists
	resource, exists := cli.store.Get(*name)
	if !exists {
		fmt.Fprintf(os.Stderr, "Error: resource %q not found\n", *name)
		return ExitResourceError
	}

	// Confirmation prompt unless --force is used
	if !*force {
		fmt.Printf("Warning: This will permanently delete the resource '%s' (%s, %dGB)\n",
			resource.Name, resource.Type, resource.Size)
		fmt.Print("Are you sure? [y/N]: ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Cancelled")
			return ExitSuccess
		}
	} else if cli.verbose {
		fmt.Println("Warning: This will permanently delete the resource")
	}

	if err := cli.store.Delete(*name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitResourceError
	}

	fmt.Printf("✓ Deleted resource: %s\n", *name)
	return ExitSuccess
}

// runUpdate handles the update subcommand
func (cli *CLI) runUpdate(args []string) int {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	name := fs.String("name", "", "Name of the resource to update")
	size := fs.Int("size", 0, "New size in GB")
	resourceType := fs.String("type", "", "New type")
	showHelp := fs.Bool("help", false, "Show help")

	if err := fs.Parse(args); err != nil {
		return ExitUsageError
	}

	if *showHelp {
		printUpdateHelp()
		return ExitSuccess
	}

	if *name == "" {
		fmt.Fprintf(os.Stderr, "Error: --name is required\n\nRun 'myapp update --help' for usage.\n")
		return ExitUsageError
	}

	// Check that at least one update field is provided
	if *size == 0 && *resourceType == "" {
		fmt.Fprintf(os.Stderr, "Error: at least one of --size or --type must be provided\n")
		return ExitUsageError
	}

	// Validate new type if provided
	if *resourceType != "" {
		validTypes := map[string]bool{
			"postgres":      true,
			"mysql":         true,
			"redis":         true,
			"mongodb":       true,
			"elasticsearch": true,
		}
		if !validTypes[strings.ToLower(*resourceType)] {
			fmt.Fprintf(os.Stderr, "Error: --type must be one of: postgres, mysql, redis, mongodb, elasticsearch\n")
			return ExitUsageError
		}
	}

	// Prepare update values
	var sizePtr *int
	var typePtr *string
	if *size > 0 {
		sizePtr = size
	}
	if *resourceType != "" {
		lowerType := strings.ToLower(*resourceType)
		typePtr = &lowerType
	}

	resource, err := cli.store.Update(*name, sizePtr, typePtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return ExitResourceError
	}

	// Print what was updated
	var updates []string
	if sizePtr != nil {
		updates = append(updates, fmt.Sprintf("size to %dGB", *sizePtr))
	}
	if typePtr != nil {
		updates = append(updates, fmt.Sprintf("type to %s", *typePtr))
	}

	fmt.Printf("✓ Updated %s %s\n", resource.Name, strings.Join(updates, " and "))
	return ExitSuccess
}

// Run executes the CLI application
func (cli *CLI) Run(args []string) int {
	if len(args) < 1 {
		printHelp()
		return ExitUsageError
	}

	switch args[0] {
	case "create":
		return cli.runCreate(args[1:])
	case "list":
		return cli.runList(args[1:])
	case "delete":
		return cli.runDelete(args[1:])
	case "update":
		return cli.runUpdate(args[1:])
	case "--help", "-h", "help":
		printHelp()
		return ExitSuccess
	case "--version", "-v":
		fmt.Printf("myapp version %s\n", Version)
		return ExitSuccess
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command %q\n\nRun 'myapp --help' for usage.\n", args[0])
		return ExitUsageError
	}
}

func main() {
	// Parse global flags first
	var configPath string
	var verbose bool
	var showVersion bool
	var showHelp bool

	// Find where global flags end and subcommand begins
	globalArgs := []string{}
	subcommandArgs := []string{}
	foundSubcommand := false

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

	// Parse global flags
	globalFS := flag.NewFlagSet("global", flag.ContinueOnError)
	globalFS.StringVar(&configPath, "config", "", "Config file path for persistence")
	globalFS.BoolVar(&verbose, "verbose", false, "Enable verbose output")
	globalFS.BoolVar(&showVersion, "version", false, "Show version")
	globalFS.BoolVar(&showHelp, "help", false, "Show help")

	if err := globalFS.Parse(globalArgs); err != nil {
		os.Exit(ExitUsageError)
	}

	if showVersion {
		fmt.Printf("myapp version %s\n", Version)
		os.Exit(ExitSuccess)
	}

	if showHelp || len(subcommandArgs) == 0 {
		printHelp()
		if len(subcommandArgs) == 0 && !showHelp {
			os.Exit(ExitUsageError)
		}
		os.Exit(ExitSuccess)
	}

	// Create the store and CLI
	store := NewResourceStore(configPath, verbose)
	cli := &CLI{
		store:   store,
		verbose: verbose,
	}

	// Run the CLI
	exitCode := cli.Run(subcommandArgs)
	os.Exit(exitCode)
}
