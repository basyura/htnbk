# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is "htnbk" (hatena blog backup) - a Go tool for extracting articles from Hatena Blog and outputting each entry to separate files for backup purposes.

## Development Commands

### Building and Running
- `go run main.go` - Run the application directly
- `go build` - Build the binary
- `go build -o htnbk` - Build with specific binary name

### Testing and Quality
- `go test ./...` - Run all tests
- `go vet ./...` - Run Go's static analysis tool
- `go fmt ./...` - Format all Go files

### Dependencies
- `go mod tidy` - Clean up module dependencies
- `go mod download` - Download dependencies

## Architecture

This is a simple Go command-line tool with a single main.go file. The current implementation is minimal (prints "hello") and appears to be in early development phase for what will become a Hatena Blog backup utility.

## Project Structure

- `main.go` - Main entry point
- `go.mod` - Go module definition (Go 1.24.5)
- `README.md` - Project documentation in Japanese