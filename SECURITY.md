# Security Policy

## Supported versions

GoKit is currently in active development. Security fixes are applied to the default development branch. Released versions may not receive fixes indefinitely until a formal support policy is established.

## Reporting a vulnerability

Please do not report security vulnerabilities through public GitHub issues.

Use GitHub's private security reporting feature for this repository when available. If private reporting is unavailable, contact the maintainer privately before disclosing the issue publicly.

Include:

- a clear description of the vulnerability
- affected package, endpoint, or feature
- reproduction steps or a minimal proof of concept
- potential impact
- any suggested mitigation

Please avoid including secrets, personal information, or live credentials in reports.

## Response process

The maintainer will assess the report, attempt to reproduce the issue, determine impact, and coordinate a fix or mitigation. Once a fix is available, the project may publish an advisory or release note when appropriate.

## Scope

Security reports are especially relevant to the interpreter, HTTP API, CLI, and web frontend. GoKit is designed to run locally by default; deployments that expose the HTTP server to untrusted networks are responsible for adding appropriate authentication, authorization, TLS, rate limiting, and network controls.
