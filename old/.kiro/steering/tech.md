# Technology Stack

## Core Technologies
- **Language**: Go 1.24.2
- **Runtime**: TinyGo (for microcontroller compilation)
- **Target Platform**: nrf52xxx microcontrollers
- **Dependency Injection**: Google Wire

## Dependencies
- `github.com/google/wire` - Compile-time dependency injection

## Build System

### TinyGo Commands
The project uses TinyGo for microcontroller compilation:

```bash
# Build firmware
tinygo build -target=feather-nrf52840 --size short -o build/firmware.uf2 .

# Flash to device
tinygo flash -target=feather-nrf52840 --size short .

# Flash with monitor
tinygo flash -target=feather-nrf52840 --size short -monitor
```

### Make Commands
Use the Makefile in the `example/` directory:

```bash
# Build firmware (creates build/firmware.uf2)
make build

# Flash firmware to device
make flash

# Flash with serial monitor
make flash-monitor

# Clean build artifacts
make clean
```

### Testing
Standard Go testing:

```bash
# Run all tests
go test ./tests/...

# Run tests with verbose output
go test -v ./tests/...

# Run specific test
go test ./tests/ -run TestTGKManager
```

### Code Generation
Wire dependency injection requires code generation:

```bash
# Generate wire dependencies (run from di/ directory)
wire
```

## Configuration
- **Keyboard Config**: JSON files (e.g., `keyboard.json`)
- **Build Target**: Specified via TinyGo `-target` flag
- **Log Levels**: Configurable via JSON config (`debug`, `info`, `warn`, `error`)