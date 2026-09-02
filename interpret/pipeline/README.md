# Pipeline

`interpret/pipeline` is GoKit's unified human-input execution boundary.

It combines the intent layer, domain router, typed interpretation, execution, and explainable plans behind one API.

## Flow

```text
Human Input
    ↓
Intent Interpreter
    ↓
Typed Interpretation
    ↓
Execution
    ↓
Explainable Plan
    │
    └── no intent match → Domain Router
```

## Usage

```go
result, err := pipeline.Parse("what is 20% of 500")
if err != nil {
    // handle interpretation or execution failure
}

fmt.Println(result.Kind)
fmt.Println(result.Value)
```

Use `ParseAt` when contextual calculations need deterministic time:

```go
result, err := pipeline.ParseAt(
    "how old is someone born 11-11-2011",
    time.Date(2026, 9, 2, 0, 0, 0, 0, time.UTC),
)
```

The result preserves both the original typed interpretation and the executed value, alongside confidence, assumptions, and a step-by-step `plan.Plan`.

## Design principles

- deterministic rather than probabilistic
- explicit failures rather than silent guessing
- typed values rather than unstructured strings
- explainable execution plans
- contextual operations can be made deterministic with `ParseAt`
- existing domain packages remain independently reusable
