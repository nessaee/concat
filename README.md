# concat

![Version](https://img.shields.io/github/v/release/nessaee/concat?style=flat-square)
![Go Version](https://img.shields.io/github/go-mod/go-version/nessaee/concat?style=flat-square)
![License](https://img.shields.io/github/license/nessaee/concat?style=flat-square)

**concat** is a high-performance, cross-platform CLI tool engineered to aggregate project source code into a single, formatted text block. It is designed specifically for developers who need to provide code context to Large Language Models (LLMs) or share code snippets efficiently.

## Features

- 🚀 **High Performance:** Built with Go for blazing fast recursive directory traversal.
- 🛡️ **Smart Filtering:** Automatically respects `.gitignore` files (nested and root) and filters out common noise (e.g., `node_modules`, `.git`, build artifacts).
- 📋 **Clipboard Integration:** Automatically copies output to the system clipboard on Linux, macOS, and Windows.
- 🌳 **Tree Context:** Optionally generates a visual directory tree to preserve structural context for the LLM.
- 📦 **Zero Dependency:** Distributed as a single static binary.

## Installation

### Option 1: Go Install (Recommended for Developers)
If you have Go 1.22+ installed, this is the easiest way to get the latest version.

```bash
go install github.com/nessaee/concat/cmd/concat@latest
```
*Ensure that your Go bin directory (`$(go env GOPATH)/bin`) is in your system's `PATH`.*

### Option 2: Binary Download
Download the pre-compiled binary for your operating system from the [Releases Page](https://github.com/nessaee/concat/releases/latest).

**Linux / macOS:**
1. Download the `.tar.gz` archive.
2. Extract the binary: `tar -xzf concat_*.tar.gz`
3. Move it to your path: `sudo mv concat /usr/local/bin/`

**Windows:**
1. Download the `.zip` archive.
2. Extract `concat.exe`.
3. Place it in a folder included in your system `PATH`.

### Option 3: Build from Source
```bash
git clone https://github.com/nessaee/concat.git
cd concat
go build -ldflags="-s -w" -o concat cmd/concat/main.go
```

## Platform Specifics

### Linux
`concat` relies on standard system tools for clipboard access. Please ensure one of the following is installed:
- **Wayland:** `wl-clipboard` (`sudo apt install wl-clipboard`)
- **X11:** `xclip` or `xsel` (`sudo apt install xclip`)

### WSL (Windows Subsystem for Linux)
`concat` supports WSL by piping output to the Windows clipboard via `clip.exe` automatically if native Linux tools are missing.

## Usage

The basic syntax is:
```bash
concat -p <extension> [flags]
```

### Common Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--pattern` | `-p` | **Required.** File extension to include (e.g., `go`, `py`, `js`). Can be repeated. |
| `--ignore` | `-i` | Glob pattern to ignore (e.g., `tests/*`, `*.log`). Can be repeated. |
| `--tree` | `-t` | Prepend a visual directory tree structure to the output. |
| `--output` | `-o` | Write result to a file instead of the clipboard. |
| `--help` | `-h` | Show help message. |

### Examples

**1. The "LLM Context" Run**
Grab all Go source files and the `go.mod` file, include the directory tree for context, and copy to clipboard.
```bash
concat -p go -p mod -t
```

**2. Frontend Project**
Grab TypeScript and CSS files, but ignore the `tests` folder and `stories` files.
```bash
concat -p ts -p tsx -p css -i "tests/*" -i "*.stories.tsx"
```

**3. Export to File**
Useful for archiving or processing with other tools.
```bash
concat -p py -o codebase.txt
```

## Default Behavior

By default, `concat` enforces these ignore patterns to keep your context clean:
- **Directories:** `.git`, `node_modules`, `__pycache__`, `.venv`, `venv`, `target`, `dist`, `build`
- **Files:** `*.log`, `*.lock`, `*.swp`, `.DS_Store`, `*.exe`, `*.dll`, `*.so`

*It also respects any `.gitignore` files found during traversal.*

## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes.
4. Push to the branch.
5. Open a Pull Request.

---
*Released under the MIT License.*
