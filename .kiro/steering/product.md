# TGK Product Overview

TGK is a keyboard firmware framework built with TinyGo, specifically designed for nrf52xxx microcontrollers. It provides a modular architecture for creating custom keyboard firmware with support for various matrix types, HID interfaces, and advanced features like split keyboards and key remapping.

## Key Features
- **TinyGo-based**: Optimized for microcontrollers with minimal resource usage
- **Modular Architecture**: Service-oriented design with dependency injection
- **Multiple Matrix Types**: Support for MX switches and EC (electrostatic capacitive) switches
- **Configurable HID**: Pluggable HID interface system (USB, BLE)
- **Split Keyboard Support**: Built-in support for split keyboard configurations
- **Advanced Key Processing**: Layer management, key remapping, and macro support
- **JSON Configuration**: Keyboard layouts and settings defined in JSON files

## Target Hardware
- Primary: nrf52xxx microcontrollers (Nordic Semiconductor)
- Development board: Adafruit Feather nRF52840
- Custom keyboard PCBs with matrix scanning capabilities

## Architecture Pattern
The firmware follows a layered architecture with clear separation of concerns:
- **Manager Layer**: Central orchestration (TGKManager)
- **Service Layer**: Business logic (HID, KeyScan, Layer, Config, etc.)
- **Repository Layer**: Hardware abstraction (GPIO, Config, Keymap)
- **Interface Layer**: Pluggable implementations for different hardware