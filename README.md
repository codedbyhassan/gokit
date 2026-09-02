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

GoKit accepts useful human input without forcing every caller into one exact format.

The `interpret` package provides a shared interpretation layer that can normalize input, detect patterns, generate interpretations, score confidence, report ambiguity, and return typed Go values.

The `interpret/intent` package composes those interpreters into requests such as:

```text
what's 20% of 500?
convert 10 miles to km
how old is someone born 11-11-2011
calculate 5kg plus 200g
five kilograms plus two hundred grams
```

## Quantity Engine

GoKit's quantity engine supports natural measurements, compatible-unit arithmetic, conversion, scalar multiplication/division, and dimensional validation.

Examples:

```text
5kg + 200g                         → 5.2 kg
five kilograms plus two hundred grams → 5.2 kg
10 miles + 5 km                    → 13.106855 mi
2m * 4                             → 8 m
5kg / 2                            → 2.5 kg
convert five kilometers to miles   → 3.106855... mi
```

Supported dimensions include length, mass, temperature, speed, and volume. Unit aliases include symbols, singular/plural names, and common natural-language forms.

The typed APIs are available through `interpret/unit.ParseQuantity`, `interpret/unit.ParseConversion`, and `interpret/expression.EvaluateQuantity`.

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

Quantity expressions use a typed AST:

```text
Input
 ↓
Normalize natural operators
 ↓
Parse quantities + scalar values
 ↓
Build expression AST
 ↓
Validate dimensions / arithmetic rules
 ↓
Convert compatible units
 ↓
Evaluate
 ↓
Typed QuantityResult
```

The goal is not to pretend GoKit is an LLM. It provides deterministic, explainable, testable intelligence that can be extended by domain-specific interpreters.

## Philosophy

GoKit turns small programming exercises into reusable, testable, exported Go packages that can be consumed by CLIs, APIs, web applications, and larger systems.

The intelligence layer follows a **bounded intelligence** principle: prefer deterministic interpretation, explicit confidence, reusable domain logic, and safe failure over silent guessing.

## Status

Phase 1 — Expression Engine: **complete**

Phase 2 — Quantity Engine: **complete**

Overall project: 🚧 Early development
