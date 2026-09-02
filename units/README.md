# Units

`units` provides deterministic unit definitions, dimensions, aliases, and conversions. The interpreter layer in `interpret/unit` adds natural-language parsing on top.

## Dimensions

- Length: m, cm, mm, km, mi, yd, ft, in
- Mass: mg, g, kg, tonne, oz, lb
- Temperature: °C, °F, K
- Speed: m/s, km/h, mph
- Volume: mL, L, gal

## Examples

```go
km, _ := units.Lookup("km")
miles, _ := units.Lookup("miles")
value, _ := units.Convert(10, km, miles)
```

For human input:

```go
quantity, _ := unit.ParseQuantity("five kilograms")
conversion, _ := unit.ParseConversion("convert five kilometers to miles")
```

Conversions reject incompatible dimensions and preserve the requested target unit.
