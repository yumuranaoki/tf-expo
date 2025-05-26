# tf-expo (TFX)

A CLI tool for visualizing Terraform plan output differences.

## Overview

tf-expo (command: `tfx`) is a tool that takes the JSON output from Terraform's `terraform plan` command and allows you to interactively select resource changes to view detailed differences.

## Installation

```bash
go install github.com/yumuranaoki/tf-expo@latest
```

## Usage

1. First, output your Terraform plan in JSON format:

```bash
terraform plan -out=tfplan
terraform show -json tfplan > plan.json
```

2. Pass the JSON output to the `tfx` command:

```bash
cat plan.json | tfx
```

3. An interactive selection screen will appear where you can choose resources that will be changed
4. Detailed differences for the selected resource will be displayed

## Filtering Options

The `tfx` command provides the following filtering options:

* `--action`: Filter by action (create, update, delete, replace)
* `--target`: Filter by module/target prefix

Examples:

```bash
# Show only resources being created
cat plan.json | tfx --action create

# Show only resources in a specific module
cat plan.json | tfx --target module.network
```

## Features

* Parse Terraform plan output JSON
* Interactive resource change selection (fuzzy finder)
* Detailed diff display of selected resources
* Filtering by actions or module names

## License

MIT
