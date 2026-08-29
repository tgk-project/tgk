# Project Structure

## Root Level Files
- `go.mod` / `go.sum` - Go module definition and dependencies
- `types.go` - Core type definitions and interfaces
- `tgk_manager.go` - Main manager orchestrating all services
- `keycode.go` - Keyboard key code definitions
- `service_*.go` - Service layer implementations
- `repository_*.go` - Repository layer implementations

## Directory Structure

### `/di/`
Dependency injection configuration using Google Wire:
- `wire.go` - Wire provider sets and initialization
- `wire_gen.go` - Generated dependency injection code

### `/example/`
Reference implementation and build configuration:
- `main.go` - Arduino-style main loop example
- `Makefile` - Build commands for TinyGo
- `keyboard.json` - Example keyboard configuration
- `keymap.go` - Example keymap implementation
- `/build/` - Compiled firmware output

### `/tests/`
Unit tests for all components:
- `*_test.go` - Test files following Go conventions
- Tests are organized by component (service, repository, manager)

### `/hid/`
HID interface implementations:
- `usb.go` - USB HID implementation
- `ble.go` - Bluetooth LE HID implementation

### `/keyscan/`
Matrix scanning implementations:
- `factory.go` - Factory pattern for matrix types
- `mx.go` - MX switch matrix scanning
- `ec.go` - Electrostatic capacitive matrix scanning

### `/debugger/`
Development and debugging tools:
- `main.go` - Debug utilities
- `keyboard.json` - Debug configuration

## Naming Conventions

### Files
- `service_*.go` - Service layer implementations
- `repository_*.go` - Repository layer implementations
- `*_test.go` - Test files
- `*.json` - Configuration files

### Interfaces
- Service interfaces: `*Service` (e.g., `HIDService`, `ConfigService`)
- Repository interfaces: `*Repository` (e.g., `ConfigRepository`)
- Hardware interfaces: `*Interface` (e.g., `HIDInterface`)

### Implementations
- Service implementations: `*Service` struct with lowercase name
- Use constructor pattern: `New*Service()` functions
- Repository implementations follow same pattern

## Architecture Layers

### Manager Layer
- `TGKManager` - Central orchestrator
- Coordinates all services and handles main execution loop

### Service Layer
- Business logic and coordination
- Each service has a clear responsibility
- Services communicate through well-defined interfaces

### Repository Layer
- Hardware abstraction
- Direct interaction with GPIO, storage, etc.
- Provides data access for services

### Interface Layer
- Pluggable implementations for different hardware
- Allows runtime selection of HID, matrix types, etc.