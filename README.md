# concat & opt

![Version](https://img.shields.io/github/v/release/nessaee/concat?style=flat-square)
![Go Version](https://img.shields.io/github/go-mod/go-version/nessaee/concat?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)

**concat** is a high-performance CLI suite designed to streamline the gathering and preparation of codebase context for Large Language Models (LLMs). It adheres to the Unix philosophy of doing one thing well and piping output to the next tool.

The suite consists of:
1.  **`concat` (The Gatherer):** A blazing fast tool for traversing directories, respecting `.gitignore`, and gathering code files.
2.  **`opt` (The Refiner):** A stream optimizer that reduces token count by stripping whitespace and license headers.

---

## ⚡ The Power Pipe (Quick Start)

The most effective way to use these tools is by chaining them. This workflow gathers your code, optimizes it for token efficiency, and prepares it for your prompt.

```bash
# Gather all Go files, optimize content, and copy to clipboard
concat -p go | opt --compact --strip-headers
```

**What just happened?**
1.  `concat` found all `.go` files (ignoring `.gitignore` rules).
2.  It streamed them to stdout.
3.  `opt` received the stream, removed excessive newlines and license headers.
4.  `opt` detected it wasn't being piped elsewhere, so it copied the result to your **clipboard**.

---

## Installation

### Go Install (Recommended)
Requires Go 1.21+.

```bash
go install github.com/nessaee/concat/cmd/concat@latest
go install github.com/nessaee/concat/cmd/opt@latest
```

### Binary Download
Pre-compiled binaries are available on the [Releases Page](https://github.com/nessaee/concat/releases/latest).

---

## `concat`: The Gatherer

`concat` is designed to be smart about what it includes. It automatically ignores `.git`, `node_modules`, and binary files unless explicitly requested.

### Basic Usage

```bash
# Syntax
concat -p <extension> [flags]
```

### Common Scenarios

**1. Gather specific languages:**
```bash
concat -p js -p ts -p css
```

**2. Exclude test files:**
Great for reducing noise when asking for feature implementation.
```bash
concat -p go --no-tests
```

**3. Generate an XML Prompt:**
Output in XML format (`<file path="...">`) for structured prompting.
```bash
concat -p py --xml
```

**4. Include a Directory Tree:**
Help the LLM understand the project structure.
```bash
concat -p rs --tree
```

### Flag Reference

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--pattern` | `-p` | **Required.** File extension to include (e.g., `go`). Can be repeated. | `[]` |
| `--ignore` | `-i` | Glob pattern to ignore (e.g., `tests/*`). Can be repeated. | `[]` |
| `--no-tests`| `-n` | Exclude common test files (`_test.go`, `.spec.ts`, etc). | `false` |
| `--tree` | `-t` | Prepend a directory tree structure to the output. | `false` |
| `--xml` | `-x` | Output in XML format (`<file path="...">`) instead of Markdown. | `false` |
| `--output` | `-o` | Write output to a specific file. | `(Clipboard)` |
| `--stdout` | `-s` | Force printing to stdout (auto-enabled when piping). | `false` |

---

## `opt`: The Refiner

`opt` is a stream processor. It doesn't read files directly; it reads text from stdin. This makes it compatible with `concat`, `cat`, `grep`, or any other stdout tool.

### Basic Usage

```bash
# Syntax
<command> | opt [flags]
```

### Common Scenarios

**1. Check Token Cost:**
Estimate the token count of your codebase without copying it.
```bash
concat -p go | opt --cost
```
*Output (stderr): `Estimated Tokens: ~1420`*

**2. Maximize Context Window:**
Aggressively strip noise to fit more code into the prompt.
```bash
concat -p ts | opt --compact --strip-headers
```

### Flag Reference

| Flag | Short | Description |
|------|-------|-------------|
| `--compact` | `-c` | Reduce multiple newlines to a single newline. |
| `--strip-headers` | | Remove C-style, Go-style, and Hash-style license headers. |
| `--cost` | | Print estimated token count and cost to stderr. |
| `--stdout` | `-s` | Force print to stdout instead of clipboard. |

---

## Configuration

You can persist your preferences in a `.concat.yaml` or `.forge.yaml` file in your project root.

```yaml
# .concat.yaml
extensions:
  - go
  - mod
ignore:
  - "vendor/*"
  - "migration/*"
tree: true
no_tests: true
```

## Contributing

We welcome contributions! Please see the [RUNBOOK](RUNBOOK.md) for details on building and testing the project locally.

## License

MIT
