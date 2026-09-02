# Unit Interpretation

The `interpret/unit` package turns human-friendly measurements into typed quantities and conversions.

```go
result, err := unit.Parse("20 km to miles")
if err != nil {
    // handle invalid or unsupported input
}

conversion := result.Value.(unit.Conversion)
// conversion.Value ≈ 12.4274
```

Supported input shapes include:

- `5 kg`
- `10 kilometers`
- `25 celsius`
- `100 miles per hour`
- `5kg to pounds`
- `20 km to miles`
- `25 celsius to fahrenheit`

The parser is deliberately deterministic: unsupported units and incompatible dimensions return errors instead of silently guessing.
