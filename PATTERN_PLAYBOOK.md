# Pattern Playbook: Golang CLI Architecture

_Extracted from golang-pokedex codebase analysis_

## Table of Contents

1. [Architectural Patterns](#architectural-patterns)
2. [Design Patterns](#design-patterns)
3. [Code Organization](#code-organization)
4. [Best Practices](#best-practices)
5. [Reusable Components](#reusable-components)
6. [Implementation Templates](#implementation-templates)

---

## Architectural Patterns

### 1. REPL (Read-Eval-Print Loop) Architecture

**What it is:** Interactive command-line interface that continuously reads user input, evaluates commands, and prints results.

**Why it's effective:**

-   Provides immediate feedback to users
-   Creates an interactive experience
-   Easy to extend with new commands
-   Natural for CLI tools and debugging interfaces

**How to implement:**

```go
func startRepl(cfg *config) {
    reader := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("YourApp > ")
        reader.Scan()
        words := cleanInput(reader.Text())

        if len(words) == 0 {
            continue
        }

        command := words[0]
        cliCommand, exists := getCommands()[command]

        if exists {
            err := cliCommand.callback(cfg)
            if err != nil {
                fmt.Println("Error:", err)
            }
        } else {
            fmt.Println("Invalid command")
        }
    }
}
```

**When to use:**

-   CLI tools requiring user interaction
-   Developer tools and debuggers
-   Administrative interfaces
-   Learning/tutorial applications

### 2. Clean Architecture with Layered Design

**What it is:** Separation of concerns using internal packages and dependency injection.

**Why it's effective:**

-   Clear separation between business logic and external dependencies
-   Easy to test and mock external services
-   Maintainable and extensible codebase
-   Domain logic is protected from external changes

**How to implement:**

```
project/
├── main.go                    # Entry point & dependency wiring
├── repl.go                    # Presentation layer
├── command_*.go               # Application layer
└── internal/
    └── service/               # Domain/Service layer
        ├── client.go         # External interface
        ├── types.go          # Domain models
        └── service_methods.go # Business logic
```

**When to use:**

-   Medium to large applications
-   Applications with external API dependencies
-   When testability is important
-   Long-term maintainable projects

### 3. Configuration-Driven Architecture

**What it is:** Central configuration object passed through the application containing all dependencies and state.

**Why it's effective:**

-   Single source of truth for application state
-   Easy dependency injection
-   Simplifies testing with mock configurations
-   Clear data flow through the application

**How to implement:**

```go
type Config struct {
    APIClient    APIClient
    Database     Database
    Cache        Cache
    Settings     AppSettings
}

func NewConfig() *Config {
    return &Config{
        APIClient: NewAPIClient(timeout),
        Database:  NewDatabase(connectionString),
        Cache:     NewCache(size),
    }
}
```

**When to use:**

-   Applications with multiple external dependencies
-   When you need to share state across commands/handlers
-   Testing scenarios requiring dependency injection
-   Applications requiring runtime configuration

---

## Design Patterns

### 1. Command Pattern

**What it is:** Encapsulates commands as objects with consistent interface, enabling parameterization and queuing.

**Why it's effective:**

-   Decouples command invoker from command execution
-   Easy to add new commands without modifying existing code
-   Supports undo/redo functionality
-   Enables command composition and chaining

**How to implement:**

```go
type Command struct {
    name        string
    description string
    callback    func(*Config) error
}

func getCommands() map[string]Command {
    return map[string]Command{
        "help": {
            name:        "help",
            description: "Show help information",
            callback:    commandHelp,
        },
        "list": {
            name:        "list",
            description: "List items",
            callback:    commandList,
        },
    }
}

// Command implementation
func commandHelp(cfg *Config) error {
    // Implementation here
    return nil
}
```

**When to use:**

-   CLI applications with multiple commands
-   GUI applications with menu actions
-   Macro recording systems
-   Queue-based processing systems

### 2. Factory Pattern

**What it is:** `getCommands()` function serves as a factory for creating command objects.

**Why it's effective:**

-   Centralizes object creation logic
-   Easy to modify without affecting client code
-   Supports polymorphism and interface-based design
-   Reduces coupling between components

**How to implement:**

```go
type CommandFactory struct {
    config *Config
}

func NewCommandFactory(cfg *Config) *CommandFactory {
    return &CommandFactory{config: cfg}
}

func (f *CommandFactory) GetCommands() map[string]Command {
    return map[string]Command{
        "command1": f.createCommand1(),
        "command2": f.createCommand2(),
    }
}

func (f *CommandFactory) createCommand1() Command {
    return Command{
        name:     "command1",
        callback: func(cfg *Config) error { /* implementation */ },
    }
}
```

**When to use:**

-   Creating families of related objects
-   When object creation logic is complex
-   Plugin systems
-   When you need to abstract object creation

### 3. Dependency Injection Pattern

**What it is:** Dependencies are passed into objects rather than created internally.

**Why it's effective:**

-   Improves testability through mock injection
-   Reduces coupling between components
-   Makes dependencies explicit and clear
-   Supports different implementations (dev vs prod)

**How to implement:**

```go
// Define interfaces for dependencies
type APIClient interface {
    GetData(url string) ([]byte, error)
}

// Inject dependencies through constructor
func NewService(client APIClient, config Config) *Service {
    return &Service{
        client: client,
        config: config,
    }
}

// Or through method parameters
func ProcessCommand(cfg *Config) error {
    data, err := cfg.APIClient.GetData(url)
    // Process data...
}
```

**When to use:**

-   Applications requiring testability
-   Systems with multiple environment configurations
-   When using external services or databases
-   Microservices architectures

### 4. Builder Pattern (HTTP Client)

**What it is:** Step-by-step construction of complex objects with optional parameters.

**Why it's effective:**

-   Handles complex object creation elegantly
-   Provides fluent interface for configuration
-   Allows optional parameters without method overloading
-   Immutable object creation

**How to implement:**

```go
type ClientBuilder struct {
    timeout    time.Duration
    retries    int
    baseURL    string
    headers    map[string]string
}

func NewClientBuilder() *ClientBuilder {
    return &ClientBuilder{
        timeout: 30 * time.Second,
        retries: 3,
        headers: make(map[string]string),
    }
}

func (b *ClientBuilder) WithTimeout(timeout time.Duration) *ClientBuilder {
    b.timeout = timeout
    return b
}

func (b *ClientBuilder) WithRetries(retries int) *ClientBuilder {
    b.retries = retries
    return b
}

func (b *ClientBuilder) Build() Client {
    return Client{
        httpClient: http.Client{Timeout: b.timeout},
        retries:    b.retries,
        baseURL:    b.baseURL,
        headers:    b.headers,
    }
}

// Usage:
client := NewClientBuilder().
    WithTimeout(5*time.Second).
    WithRetries(5).
    Build()
```

**When to use:**

-   Objects with many optional parameters
-   Complex configuration scenarios
-   When you want fluent/chainable APIs
-   Immutable object construction

---

## Code Organization

### 1. Package Structure Pattern

**What it is:** Logical organization using internal packages and clear boundaries.

**Structure:**

```
project/
├── main.go              # Entry point
├── app_layer.go         # Application logic
├── internal/           # Private packages
│   └── domain/         # Business logic
│       ├── client.go   # External interfaces
│       ├── types.go    # Data models
│       └── service.go  # Business methods
└── test/               # Test utilities
```

**Why it's effective:**

-   Enforces encapsulation with internal packages
-   Clear separation of concerns
-   Easy to navigate and understand
-   Prevents circular dependencies

**When to use:**

-   Medium to large Go projects
-   Libraries that need to hide implementation details
-   Team projects requiring clear boundaries

### 2. File Naming Conventions

**Patterns identified:**

-   `command_*.go` - Command implementations
-   `*_test.go` - Test files
-   `types_*.go` - Type definitions
-   Descriptive, action-based naming

**Guidelines:**

-   Use snake_case for multi-word files
-   Group related functionality by prefix
-   Keep file names descriptive but concise
-   Separate test files clearly

### 3. Function Organization Pattern

**What it is:** Single responsibility functions with clear input/output contracts.

**Example:**

```go
// Utility functions
func cleanInput(input string) []string {
    return strings.Fields(input)
}

// Command functions with consistent signature
func commandHelp(cfg *Config) error {
    // Implementation
    return nil
}

// Factory functions
func getCommands() map[string]Command {
    // Return command map
}
```

**Best practices:**

-   Consistent function signatures within categories
-   Clear, action-based function names
-   Single responsibility per function
-   Error handling as return value

---

## Best Practices

### 1. Error Handling Pattern

**What it is:** Consistent error propagation and handling throughout the application.

**Implementation:**

```go
// Commands return errors for centralized handling
func commandProcess(cfg *Config) error {
    data, err := cfg.Client.FetchData()
    if err != nil {
        return fmt.Errorf("failed to fetch data: %w", err)
    }

    err = processData(data)
    if err != nil {
        return fmt.Errorf("failed to process data: %w", err)
    }

    return nil
}

// REPL handles errors gracefully
if err := command.callback(cfg); err != nil {
    fmt.Println("Error executing command:", err)
    // Continue execution, don't crash
}
```

**Benefits:**

-   Graceful degradation
-   Clear error context
-   Non-crashing application behavior
-   Wrapped errors for debugging

### 2. Input Validation and Sanitization

**What it is:** Clean and validate user input before processing.

**Implementation:**

```go
func cleanInput(input string) []string {
    words := strings.Fields(input) // Removes extra whitespace
    return words
}

// Validation in REPL
words := cleanInput(reader.Text())
if len(words) == 0 {
    continue // Skip empty input
}
```

**Benefits:**

-   Prevents injection attacks
-   Handles edge cases gracefully
-   Consistent input format
-   Better user experience

### 3. Resource Management

**What it is:** Proper resource cleanup and timeout management.

**Implementation:**

```go
// HTTP client with timeout
func NewClient(timeout time.Duration) Client {
    return Client{
        httpClient: http.Client{
            Timeout: timeout,
        },
    }
}

// Proper resource cleanup
defer resp.Body.Close()
```

**Benefits:**

-   Prevents resource leaks
-   Handles network timeouts
-   Graceful failure handling
-   Better system stability

### 4. Table-Driven Testing

**What it is:** Test cases defined as data structures for comprehensive coverage.

**Implementation:**

```go
func TestCleanInput(t *testing.T) {
    cases := []struct {
        input    string
        expected []string
    }{
        {
            input:    " hello world ",
            expected: []string{"hello", "world"},
        },
        {
            input:    "single",
            expected: []string{"single"},
        },
        {
            input:    "",
            expected: []string{},
        },
    }

    for _, c := range cases {
        actual := cleanInput(c.input)
        // Assertions...
    }
}
```

**Benefits:**

-   Easy to add new test cases
-   Clear test data organization
-   Comprehensive edge case coverage
-   Maintainable test code

---

## Reusable Components

### 1. Generic Command Framework

**Component:** Command registration and execution system

**Template:**

```go
// Generic command interface
type Command interface {
    Name() string
    Description() string
    Execute(cfg *Config) error
}

// Command registry
type CommandRegistry struct {
    commands map[string]Command
}

func (r *CommandRegistry) Register(cmd Command) {
    r.commands[cmd.Name()] = cmd
}

func (r *CommandRegistry) Execute(name string, cfg *Config) error {
    cmd, exists := r.commands[name]
    if !exists {
        return fmt.Errorf("unknown command: %s", name)
    }
    return cmd.Execute(cfg)
}
```

### 2. HTTP Client Builder

**Component:** Configurable HTTP client with sensible defaults

**Template:**

```go
type HTTPClientConfig struct {
    Timeout    time.Duration
    Retries    int
    BaseURL    string
    Headers    map[string]string
}

func NewHTTPClient(config HTTPClientConfig) *http.Client {
    return &http.Client{
        Timeout: config.Timeout,
        Transport: &http.Transport{
            // Add retry logic, circuit breaker, etc.
        },
    }
}
```

### 3. Configuration Manager

**Component:** Centralized configuration with validation

**Template:**

```go
type AppConfig struct {
    APIEndpoint string        `json:"api_endpoint" validate:"required,url"`
    Timeout     time.Duration `json:"timeout" validate:"required"`
    LogLevel    string        `json:"log_level" validate:"oneof=debug info warn error"`
}

func LoadConfig(path string) (*AppConfig, error) {
    // Load from file, environment, flags
    // Validate configuration
    // Return populated config
}
```

### 4. Input Parser Utility

**Component:** Robust input parsing and validation

**Template:**

```go
type InputParser struct {
    validators map[string]func(string) error
}

func (p *InputParser) Parse(input string) ([]string, error) {
    words := strings.Fields(strings.TrimSpace(input))

    for i, word := range words {
        if validator, exists := p.validators[fmt.Sprintf("arg%d", i)]; exists {
            if err := validator(word); err != nil {
                return nil, fmt.Errorf("invalid argument %d: %w", i, err)
            }
        }
    }

    return words, nil
}
```

---

## Implementation Templates

### 1. CLI Application Template

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type Config struct {
    // Add your dependencies here
}

type Command struct {
    Name        string
    Description string
    Callback    func(*Config) error
}

func main() {
    cfg := &Config{
        // Initialize dependencies
    }
    startREPL(cfg)
}

func startREPL(cfg *Config) {
    reader := bufio.NewScanner(os.Stdin)
    commands := getCommands()

    for {
        fmt.Print("YourApp > ")
        reader.Scan()
        words := parseInput(reader.Text())

        if len(words) == 0 {
            continue
        }

        if words[0] == "exit" {
            break
        }

        cmd, exists := commands[words[0]]
        if !exists {
            fmt.Println("Unknown command. Type 'help' for available commands.")
            continue
        }

        if err := cmd.Callback(cfg); err != nil {
            fmt.Printf("Error: %v\n", err)
        }
    }
}

func getCommands() map[string]Command {
    return map[string]Command{
        "help": {
            Name:        "help",
            Description: "Show available commands",
            Callback:    commandHelp,
        },
        // Add more commands here
    }
}

func parseInput(input string) []string {
    return strings.Fields(strings.TrimSpace(input))
}

func commandHelp(cfg *Config) error {
    fmt.Println("Available commands:")
    for name, cmd := range getCommands() {
        fmt.Printf("  %s - %s\n", name, cmd.Description)
    }
    return nil
}
```

### 2. HTTP Client Template

```go
package client

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    httpClient http.Client
    baseURL    string
}

func NewClient(baseURL string, timeout time.Duration) *Client {
    return &Client{
        httpClient: http.Client{Timeout: timeout},
        baseURL:    baseURL,
    }
}

func (c *Client) Get(endpoint string, result interface{}) error {
    url := c.baseURL + endpoint

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("failed to execute request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("API returned status %d", resp.StatusCode)
    }

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response: %w", err)
    }

    if err := json.Unmarshal(data, result); err != nil {
        return fmt.Errorf("failed to parse JSON: %w", err)
    }

    return nil
}
```

### 3. Test Suite Template

```go
package main

import (
    "testing"
)

func TestFeature(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected interface{}
        wantErr  bool
    }{
        {
            name:     "valid input",
            input:    "test input",
            expected: "expected output",
            wantErr:  false,
        },
        {
            name:     "invalid input",
            input:    "",
            expected: nil,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := functionUnderTest(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("Expected error: %v, got error: %v", tt.wantErr, err)
                return
            }

            if !tt.wantErr && result != tt.expected {
                t.Errorf("Expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

---

## Quick Reference Checklist

### Starting a New CLI Project

-   [ ] Set up REPL with command pattern
-   [ ] Create internal package for business logic
-   [ ] Implement configuration-driven architecture
-   [ ] Add proper error handling throughout
-   [ ] Create table-driven tests
-   [ ] Set up HTTP client with timeouts
-   [ ] Implement input validation
-   [ ] Add help command and usage information

### Adding New Features

-   [ ] Create new command with consistent signature
-   [ ] Add to command registry/factory
-   [ ] Implement with dependency injection
-   [ ] Add comprehensive error handling
-   [ ] Write table-driven tests
-   [ ] Update help documentation

### Code Quality Maintenance

-   [ ] Keep functions focused on single responsibility
-   [ ] Use descriptive naming conventions
-   [ ] Implement proper resource cleanup
-   [ ] Add timeouts to network operations
-   [ ] Validate all user inputs
-   [ ] Maintain consistent error handling patterns

This playbook provides a solid foundation for building maintainable, testable, and extensible CLI applications in Go using the patterns demonstrated in your Pokedex project.
