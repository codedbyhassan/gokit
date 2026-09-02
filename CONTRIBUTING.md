# Contributing to GoKit

Thank you for contributing to GoKit. The project values small, composable APIs, deterministic behavior, clear errors, strong tests, and documentation that makes packages easy to reuse.

## Before you start

1. Read the README and the package documentation relevant to your change.
2. Search existing issues and pull requests before opening a new one.
3. For substantial changes, open an issue first so the design can be discussed before implementation.

## Development requirements

- Go 1.24 or newer
- Git

GoKit intentionally avoids unnecessary runtime and frontend dependencies. The web interface is built with Go's standard library and browser APIs.

## Getting started

```bash
git clone https://github.com/codedbyhassan/gokit.git
cd gokit
go test ./...
```

Run the complete local checks:

```bash
make check
```

Or individually:

```bash
go test ./...
go test -race ./...
go vet ./...
gofmt -w .
```

Run the web application:

```bash
go run ./cmd/gokit --web --addr :8080
```

## Project structure

```text
calculator/          reusable arithmetic primitives
units/               unit definitions and conversions
interpret/            interpretation and execution engine
api/http/             HTTP/JSON adapter
cmd/gokit/             CLI and application entrypoint
web/frontend/          Go-native browser interface
tests/integration/     cross-package integration coverage
```

Keep reusable domain logic in the appropriate root package. Do not move public packages into `pkg/` simply for convention; GoKit's root-level package layout is intentional.

## Design principles

### Prefer deterministic behavior

Interpretation should be predictable, explainable, and testable. Do not introduce probabilistic or network-dependent behavior into the core engine without a strong architectural reason.

### Never silently guess

If an input has multiple materially different interpretations, represent that ambiguity rather than choosing an arbitrary answer.

### Keep APIs small

Prefer focused exported types and functions. Avoid adding abstractions that are not justified by multiple real use cases.

### Validate at boundaries

Reject malformed input early and return useful, stable errors. Do not let invalid domain state leak into execution.

### Preserve explainability

When interpretation involves assumptions, confidence, or multiple execution steps, preserve that information for callers.

### Standard library first

Use the Go standard library when it provides a clean solution. New dependencies should have a clear maintenance and capability benefit.

## Adding a new interpreter

A new domain interpreter should generally:

1. Normalize input without destroying meaningful information.
2. Parse into a typed domain value.
3. Return confidence and assumptions where appropriate.
4. Reject malformed input rather than accepting arbitrary text.
5. Provide unit tests for valid, invalid, boundary, and ambiguous cases.
6. Integrate with the router or intent layer only when appropriate.
7. Add documentation and examples for the public API.

## Testing expectations

Every behavioral change should include tests. Prefer table-driven tests for parser and conversion behavior.

At minimum, cover:

- normal inputs
- malformed inputs
- boundary values
- ambiguous inputs where applicable
- error behavior
- integration behavior when multiple packages interact

For parser-heavy changes, add fuzz coverage when it provides meaningful protection.

## Commit messages

Use concise, imperative commit messages. Conventional-style prefixes are encouraged:

```text
feat: add duration interpreter
fix: reject incomplete decimal input
test: cover ambiguous date formats
docs: clarify pipeline API
refactor: simplify router registration
```

## Pull requests

A good pull request should:

- explain what changed and why
- keep unrelated refactors out of scope
- include tests for new behavior
- update documentation when public behavior changes
- pass formatting, vet, test, and race checks
- call out compatibility or API changes clearly

Keep pull requests focused. Small, reviewable changes are preferred over large mixed commits.

## Public API compatibility

GoKit is still evolving. Before the project reaches a stable major release, public APIs may change. Even so, avoid unnecessary breaking changes and document intentional ones.

## Reporting security issues

Do not disclose suspected vulnerabilities in a public issue. See [SECURITY.md](SECURITY.md) for the reporting process.

## Code of conduct

Participation in the project is governed by [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md).

## License

By contributing to GoKit, you agree that your contributions are provided under the project's [MIT License](LICENSE).
