# GoKit

> Deterministic, explainable Go utilities with a human-input interpretation layer.

[![Go Reference](https://pkg.go.dev/badge/github.com/codedbyhassan/gokit.svg)](https://pkg.go.dev/github.com/codedbyhassan/gokit)
[![CI](https://github.com/codedbyhassan/gokit/actions/workflows/ci.yml/badge.svg)](https://github.com/codedbyhassan/gokit/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/codedbyhassan/gokit)](https://goreportcard.com/report/github.com/codedbyhassan/gokit)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

GoKit is a collection of reusable Go packages that started with small utilities and evolved into a deterministic interpretation and execution engine.

It is designed for applications that need to accept useful human input without turning the core logic into an opaque AI system. GoKit normalizes input, parses it into typed values, ranks interpretations with confidence, reports ambiguity, executes supported operations, and preserves an explainable plan.

**No cloud service. No LLM dependency. No database required.**

## Why GoKit?

Most utility libraries assume perfectly formatted input. Real applications rarely receive it.

GoKit lets callers move from rigid calls such as:

```go
units.Convert(5, "kilometers", "miles")
```

toward human-friendly input such as:

```text
convert five kilometers to miles
```

while keeping the underlying behavior deterministic and inspectable.

The project follows a **bounded intelligence** principle:

- deterministic parsing over probabilistic guessing
- typed values over unstructured strings
- explicit confidence over hidden assumptions
- explicit ambiguity over silent guesses
- reusable packages over application-specific logic
- standard-library solutions where practical

## Highlights

- Reusable arithmetic, units, dates, age, BMI, interest, and quiz packages
- Natural-language number interpretation
- Arithmetic and percentage expressions
- Quantity-aware arithmetic and unit conversion
- Intent detection and confidence-ranked routing
- Explainable execution plans
- HTTP/JSON API
- Command-line interface
- Go-native browser frontend
- Browser-local IndexedDB history and saved results
- Cross-package integration tests
- Fuzz tests and benchmarks
- CI for formatting, vetting, tests, race detection, and fuzz smoke coverage

## Local development

GoKit is a Go application and library, not a Node/npm frontend project. The browser interface is embedded into the Go binary, so there is **no frontend package installation, npm install, or separate frontend dev server**.

### Requirements

- Go 1.24+
- Git
- A C compiler is required on Windows if you want to run Go's race detector with `go test -race`.

### Clone the repository

```bash
git clone https://github.com/codedbyhassan/gokit.git
cd gokit
```

### 1. Verify the Go environment

```bash
go version
go env CGO_ENABLED
```

The normal test suite and `go vet` do not require CGO. On Windows, `go test -race` does.

### 2. Run the normal test suite

```bash
go test ./...
```

This is the first check to run after cloning or pulling changes.

### 3. Run static analysis

```bash
go vet ./...
```

### 4. Run race detection

The race detector requires CGO and a working C compiler.

On Windows, one supported setup is [MSYS2](https://www.msys2.org/) with the UCRT64 MinGW GCC toolchain.

Open **MSYS2 UCRT64** (not PowerShell) and install/update the toolchain:

```bash
pacman -Syu
```

If MSYS2 asks you to close the terminal, reopen **MSYS2 UCRT64** and run:

```bash
pacman -Su
pacman -S --needed mingw-w64-ucrt-x86_64-gcc
```

Make sure `gcc.exe` is available to Windows. The default UCRT64 location is:

```text
C:\msys64\ucrt64\bin
```

Add that directory to your Windows PATH, then open a **new PowerShell** and verify:

```powershell
gcc --version
```

Enable CGO for the current PowerShell session and run the race tests:

```powershell
$env:CGO_ENABLED="1"
go test -race ./...
```

If PowerShell reports `gcc is not recognized`, GCC is either not installed or `C:\msys64\ucrt64\bin` is not on PATH.

### 5. Run the complete local quality check

The repository Makefile combines formatting, vetting, and tests:

```bash
make check
```

The individual commands are still useful when diagnosing a failure:

```bash
go test ./...
go vet ./...
go test -race ./...
```

### 6. Start the web application

The recommended development mode is the embedded web application:

```bash
go run ./cmd/gokit --web --addr :8080
```

Then open:

```text
http://localhost:8080
```

The same Go process serves the browser frontend and the API. **Do not start `--serve` separately when using `--web`.**

The browser application provides:

- Natural-language input
- Live GoKit execution
- Confidence and assumptions
- Explainable execution steps
- Copy and save actions
- Local history
- Saved results
- Dark/light theme preference
- Responsive mobile layout

Browser data is stored in **IndexedDB**. GoKit does not require an online database or account for the web experience.

### 7. Test the HTTP API directly

If you specifically want API-only mode instead of the browser application:

```bash
go run ./cmd/gokit --serve --addr :8080
```

Then, in another terminal:

```bash
curl -X POST http://localhost:8080/v1/interpret \
  -H 'Content-Type: application/json' \
  -d '{"input":"what is 20% of 500"}'
```

On PowerShell, the equivalent is:

```powershell
Invoke-RestMethod `
  -Uri "http://localhost:8080/v1/interpret" `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"input":"what is 20% of 500"}'
```

Health check:

```text
http://localhost:8080/health
```

### 8. Run the CLI directly

The CLI can interpret a single request without starting a server:

```powershell
go run ./cmd/gokit "what is 20% of 500"
```

Other examples:

```powershell
go run ./cmd/gokit "convert 10 miles to km"
go run ./cmd/gokit "how old is someone born 11-11-2011"
go run ./cmd/gokit "five kilograms plus two hundred grams"
```

See the available commands and flags with:

```powershell
go run ./cmd/gokit --help
```

Current server flags include:

```text
-addr string
      HTTP listen address (default ":8080")
-serve
      start the GoKit HTTP API
-web
      start the GoKit web frontend
```

### Recommended development workflow

For day-to-day development, use this sequence:

```text
Clone / pull changes
      ↓
go test ./...
      ↓
go vet ./...
      ↓
go test -race ./...
      ↓
go run ./cmd/gokit --web --addr :8080
      ↓
Test the browser UI
      ↓
Commit focused changes
```

When working on a parser or interpreter, also run its package-specific tests and fuzz tests before pushing.

## Quick start

If Go and Git are already installed and you do not need race detection yet:

```bash
git clone https://github.com/codedbyhassan/gokit.git
cd gokit
go test ./...
go run ./cmd/gokit --web --addr :8080
```

Then open `http://localhost:8080`.

## Use GoKit as a library

The unified pipeline is the easiest entry point when an application wants to accept human input:

```go
package main

import (
    "fmt"
    "github.com/codedbyhassan/gokit/interpret/pipeline"
)

func main() {
    result, err := pipeline.Parse("what is 20% of 500")
    if err != nil {
        panic(err)
    }

    fmt.Println(result.Value)
    fmt.Println(result.Confidence)
}
```

For deterministic contextual operations, use `pipeline.ParseAt` with an explicit reference time.

## Human input

GoKit accepts several useful forms without requiring one exact syntax.

### Numbers

```text
1,000                    → 1000
1.5k                     → 1500
2.5 million              → 2500000
twenty five               → 25
one hundred twenty five   → 125
two point five million    → 2500000
negative one hundred      → -100
```

Malformed constructions are rejected rather than silently interpreted.

### Expressions

```text
20% of 500
25 + 17 * 2
what is 1000 / 8
20% increase on 500
```

The expression engine respects normal arithmetic precedence and exposes typed AST nodes for callers that need more control.

### Quantities

```text
5kg + 200g
five kilograms plus two hundred grams
25 km to miles
convert five kilometers to miles
2m * 4
5kg / 2
```

Quantity operations validate dimensions and perform compatible-unit conversion where appropriate.

### Dates and age

GoKit can interpret common numeric and named date formats, including formats such as:

```text
11-11-2011
2011-11-11
11-2011-11
11 November 2011
```

When a date is genuinely ambiguous, the interpretation layer can report that ambiguity instead of silently choosing a meaning.

## Architecture

```text
                           GoKit
                             │
                      Human Input
                             │
                       Normalization
                             │
                 ┌───────────┴───────────┐
                 │                       │
             Intent Layer          Domain Router
                 │                       │
                 └───────────┬───────────┘
                             │
                  Candidates + Confidence
                             │
                       Typed Value
                             │
                         Execution
                             │
                    Explainable Plan
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
         CLI             HTTP/JSON          Web UI
                                              │
                                          IndexedDB
                                        browser-local data
```

Repository layout:

```text
gokit/
├── calculator/             arithmetic primitives
├── units/                  unit definitions and conversions
├── interpret/
│   ├── date/              date interpretation
│   ├── number/            human number interpretation
│   ├── expression/        arithmetic and quantity ASTs
│   ├── unit/              quantity/conversion interpretation
│   ├── intent/            natural-language intent composition
│   ├── router/            confidence-ranked routing
│   ├── plan/              explainable execution plans
│   └── pipeline/          unified interpretation boundary
├── api/http/               HTTP/JSON adapter
├── cmd/gokit/              application entrypoint
├── web/frontend/           Go-native web interface
├── examples/               usage examples
├── docs/                   project documentation
├── internal/               private implementation helpers
├── tests/integration/      cross-package integration tests
├── .github/                CI and contribution automation
├── LICENSE
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
└── SECURITY.md
```

## HTTP API

Run the API:

```bash
go run ./cmd/gokit --serve --addr :8080
```

Interpret input:

```bash
curl -X POST http://localhost:8080/v1/interpret \
  -H 'Content-Type: application/json' \
  -d '{"input":"what is 20% of 500"}'
```

The API exposes the interpreted result together with confidence, assumptions, and an execution plan.

See [`api/http/README.md`](api/http/README.md) for the transport contract.

## CLI

The command can also interpret input directly:

```bash
go run ./cmd/gokit "what is 20% of 500"
```

For a long-running service:

```bash
go run ./cmd/gokit --serve --addr :8080
```

For the browser application:

```bash
go run ./cmd/gokit --web --addr :8080
```

## Testing and quality

GoKit treats tests as part of the package API contract.

```bash
go test ./...
go test -race ./...
go vet ./...
```

Formatting:

```bash
gofmt -w .
```

Benchmarks:

```bash
go test -bench=. ./...
```

Fuzzing examples:

```bash
go test ./interpret/number -run '^$' -fuzz Fuzz -fuzztime=10s
```

CI runs formatting checks, `go vet`, the complete test suite, race detection, and a short fuzz smoke test on pushes to `main` and pull requests.

## Extending GoKit

New domain interpreters should normally:

1. normalize input
2. parse into a typed value
3. expose useful confidence/assumption metadata
4. reject malformed input
5. cover normal, invalid, boundary, and ambiguous cases
6. integrate with the router or intent layer when appropriate
7. document public APIs and examples

See [`CONTRIBUTING.md`](CONTRIBUTING.md) before opening a pull request.

## Project status

GoKit is actively evolving and is **not yet a stable v1 API**. Public APIs may change before the first stable major release.

Current development areas:

- core reusable utilities — established
- interpretation engine — established
- unified pipeline — established
- HTTP/CLI surface — established
- local-first web interface — established
- broader domain interpreters — ongoing
- API stabilization — ongoing

## License

GoKit is released under the [MIT License](LICENSE).

## Security

Please do not disclose vulnerabilities in public issues. See [`SECURITY.md`](SECURITY.md) for the reporting policy.

## Contributing

Contributions, bug reports, tests, documentation improvements, and focused feature proposals are welcome. See [`CONTRIBUTING.md`](CONTRIBUTING.md) and [`CODE_OF_CONDUCT.md`](CODE_OF_CONDUCT.md).
