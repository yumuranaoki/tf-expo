# tfx

A CLI tool for visualizing Terraform plan output differences with an interactive interface.

## Overview

tfx is a tool that takes the JSON output from Terraform's `terraform plan` command and allows you to interactively select resource changes to view detailed differences. It provides a user-friendly way to explore large Terraform plans with filtering and fuzzy search capabilities.

## Installation

### Using go install (Recommended)

```bash
go install github.com/yumuranaoki/tfx@latest
```

This will install the `tfx` binary to your `$GOPATH/bin` directory.

### From Source

```bash
git clone https://github.com/yumuranaoki/tfx.git
cd tfx
go build -o tfx .
```

## Usage

### Basic Usage

1. Run tfx directly with Terraform plan output:

```bash
# Direct piping (recommended)
terraform plan -json | tfx
```

2. Alternative: Generate plan file first, then view:

```bash
terraform plan -out=tfplan
terraform show -json tfplan | tfx
```

3. Use the interactive interface:
   - Browse through resource changes with fuzzy search
   - Select a resource to view detailed differences
   - Press `Enter` to return to the list, `q` to quit
   - Press `Ctrl+C` to exit gracefully

### Filtering Options

tfx provides filtering options to narrow down the resources you want to review:

* `--action`: Filter by action type (create, update, delete, replace)
* `--target`: Filter by module/resource name prefix

Examples:

```bash
# Show only resources being created
terraform plan -json | tfx --action create

# Show only resources in a specific module
terraform plan -json | tfx --target module.network

# Combine filters
terraform plan -json | tfx --action update --target aws_instance
```

## Features

* **Interactive Interface**: Fuzzy search through resource changes
* **Detailed Diffs**: Color-coded differences showing before/after states  
* **Filtering**: Filter by action type or resource/module names
* **Graceful Exit**: Press `Ctrl+C` to exit cleanly
* **No-op Filtering**: Automatically excludes unchanged resources
* **Comprehensive Testing**: Full test coverage following Uber Go guidelines

## Requirements

* Go 1.21 or later
* Terminal with color support

## Contributing

Contributions are welcome! Please ensure all tests pass:

```bash
go test ./...
```

## License

MIT
