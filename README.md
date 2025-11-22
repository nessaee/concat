# concat & opt

![Version](https://img.shields.io/github/v/release/nessaee/concat?style=flat-square)
![Go Version](https://img.shields.io/github/go-mod/go-version/nessaee/concat?style=flat-square)
![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)

**concat** (v0.1.1) is a high-performance CLI suite designed to prepare codebases for Large Language Models (LLMs). It consists of two powerful tools that follow the Unix Philosophy:

1.  **`concat` (The Gatherer):** Blazing fast directory traversal, file gathering, and `.gitignore` filtering.
2.  **`opt` (The Refiner):** Semantic stream optimizer for token reduction, cost estimation, and privacy/license stripping.

## Features

- 🚀 **High Performance:** Built with Go for instant result on massive repos.
- 🛡️ **Smart Filtering:** Automatically respects `.gitignore`, excludes binary files, and offers strict inclusion lists.
- 📉 **Cost Estimation:** `opt` calculates estimated token count and API cost (Gemini Flash pricing).
- 🧹 **Context Optimization:** `opt` strips excess whitespace and corporate license headers to save context window.
- 📋 **Clipboard Integration:** `concat` copies to clipboard by default on all platforms.

## Installation

### Go Install (Recommended)

You can install both tools with a single command:

```bash
go install github.com/nessaee/concat/cmd/concat@latest
go install github.com/nessaee/concat/cmd/opt@latest
```

### Binary Download
Download pre-compiled binaries from the [Releases Page](https://github.com/nessaee/concat/releases/latest).

## Usage

### 1. The "Power Pipe" (Recommended)
Combine both tools for the ultimate LLM context preparation workflow. This pipeline finds files, minimizes them, strips noise, and estimates cost.

```bash
# Gather Go files -> Optimize content -> Output to stdout
concat -p go | opt --compact --strip-headers
```

### 2. `concat` (Gatherer)
Standalone usage for simple file aggregation.

```bash
# Syntax
concat -p <extension> [flags]

# Example: Copy all JS/TS files to clipboard (ignoring tests)
concat -p js -p ts --no-tests
```

**Common Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--pattern` | `-p` | **Required.** Extension to include (e.g., `go`, `py`). |
| `--ignore` | `-i` | Glob pattern to ignore (e.g., `tests/*`). |
| `--no-tests`| `-n` | Exclude test files (`_test.go`, `.spec.ts`, etc). |
| `--tree` | `-t` | Include directory tree at the top. |
| `--output` | `-o` | Write to file. |
| `--stdout` | `-s` | Force print to stdout (auto-detected in pipes). |

### 3. `opt` (Refiner)
Standalone usage for stream optimization.

```bash
# Syntax
cat file.txt | opt [flags]

# Example: Check cost of a file without outputting it
concat -p py | opt --cost > /dev/null
```

**Common Flags:**
| Flag | Short | Description |
|------|-------|-------------|
| `--compact` | `-c` | Reduce vertical whitespace. |
| `--strip-headers` | | Remove copyright/license headers. |
| `--cost` | | Print estimated token count and cost to stderr. |

## Default Behavior

`concat` is opinionated but flexible:
- **Ignored by default:** `.git`, `node_modules`, `__pycache__`, `vendor`, lockfiles (`go.sum`, `yarn.lock`), and binaries.
- **Override:** If you explicitly request a file type (e.g., `-p lock`), `concat` will fetch it even if it's usually ignored.

## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes.
4. Push to the branch.
5. Open a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.