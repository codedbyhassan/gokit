# Intent Interpreter

`interpret/intent` is GoKit's deterministic natural-language command layer. It sits above the domain interpreters and composes them into useful requests without introducing an LLM or external service.

## Examples

```go
result, err := intent.Parse("what's 20% of 500?")

result, err := intent.Parse("convert 10 miles to km")

result, err := intent.Parse("how old is someone born 11-11-2011")

result, err := intent.Parse("calculate 5kg plus 200g")
```

For reproducible age calculations, use `ParseAt` with an explicit reference date.

## Pipeline

```text
Natural Language
      ↓
Normalization
      ↓
Intent Detection
      ↓
Domain Interpreter
      ↓
Typed Domain Result
      ↓
GoKit Calculation / Conversion
```

The layer deliberately reuses existing packages rather than duplicating their rules. Quantity arithmetic converts compatible units before applying the arithmetic operation.

## Current intents

| Intent | Example | Result |
|---|---|---|
| `calculate` | `what's 20% of 500?` | `calculator.Result` |
| `calculate` | `calculate 5kg plus 200g` | `calculator.Result` |
| `convert` | `convert 10 miles to km` | `unit.Conversion` |
| `age` | `how old is someone born 11-11-2011` | `age.Result` |

The interpreter is intentionally bounded: unsupported or genuinely unclear requests return an error instead of inventing an interpretation.
