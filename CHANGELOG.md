# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.5] - 2025-11-23

### Added
- **XML Output Support (`--xml`):** `concat` can now output files wrapped in XML tags (`<file path="...">`) for better LLM context handling.
- **Directory Tree (`--tree`):** Added flag to prepend a directory tree structure to the output.
- **Test Exclusion (`--no-tests`):** Added flag to automatically filter out test files (`_test.go`, `.spec.ts`, etc.).
- **Clipboard Integration:** `concat` now copies to clipboard by default if no output is specified.
- **Stream Optimization (`opt`):** Introduced `opt` tool for whitespace compaction and license header stripping.

### Changed
- **Performance:** Improved directory traversal speed.
- **Configuration:** Added support for `.concat.yaml` and `.forge.yaml` configuration files.
