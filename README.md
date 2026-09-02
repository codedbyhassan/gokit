# GoKit

A collection of robust, reusable Go utilities and components — built from simple ideas into production-ready building blocks for larger Go applications.

## Packages

- Calculator
- Number Guessing
- Temperature Conversion
- Unit Conversion
- Interest Calculations
- BMI
- Age Calculations
- Quiz Engine

## Intelligent Interpretation

GoKit is designed to accept useful human input without forcing every caller into one exact format.

The `interpret` package provides a shared interpretation layer that can:

1. Normalize input.
2. Detect supported patterns.
3. Generate possible interpretations.
4. Score confidence.
5. Report ambiguity instead of silently guessing.
6. Return a strongly typed Go value that domain packages can consume.

For example, the date interpreter can understand:

```text
11-11-2011
2011-11-11
11-2011-11
11 November 2011
```

The `interpret/intent` package adds a higher-level command layer that composes those interpreters into requests such as:

```text
what's 20% of 500?
convert 10 miles to km
how old is someone born 11-11-2011
calculate 5kg plus 200g
```

The goal is not to pretend GoKit is an LLM. It is to provide deterministic, explainable, testable intelligence that can be extended by domain-specific interpreters.

## Architecture

```text
Human Input
    ↓
Normalizer
    ↓
Intent / Domain Interpreter
    ↓
Candidates + Confidence
    ↓
Typed Value
    ↓
GoKit Package
```

This keeps input interpretation separate from business logic. For example:

```text
"how old is someone born 11-11-2011"
                 ↓
        interpret/intent
                 ↓
          interpret/date
                 ↓
             time.Time
                 ↓
            age.Calculate
                 ↓
            age.Result
```

Another example composes measurements before arithmetic:

```text
"calculate 5kg plus 200g"
          ↓
    parse quantities
          ↓
     normalize units
          ↓
       calculate
          ↓
      5.2 kilograms
```

## Philosophy

GoKit turns small programming exercises into reusable, testable, exported Go packages that can be consumed by CLIs, APIs, web applications, and larger systems.

The intelligence layer follows a **bounded intelligence** principle: prefer deterministic interpretation, explicit confidence, reusable domain logic, and safe failure over silent guessing.

## Status

🚧 Early development
