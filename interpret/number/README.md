# Number Interpreter

The number interpreter converts human-friendly numeric input into typed `float64` values with confidence and format metadata.

## Supported forms

### Numeric

- `1000`
- `1,000`
- `1_000`
- `1 000`
- `1.5`
- `-42`
- `.75`

### Scaled

- `1.5k`
- `2.5m`
- `3b`
- `4 trillion`
- `2.5 million`
- `two million five hundred thousand`

### English number words

- `five`
- `twenty five`
- `one hundred and twenty five`
- `three hundred forty two`
- `one thousand two hundred and thirty four`
- `two point five`
- `one hundred twenty five point five`
- `negative one hundred`
- `minus one hundred`

Hyphenated forms such as `twenty-five` are normalized as well.

## Design

The parser is deterministic and dependency-free. It does not silently interpret arbitrary text as a number. Unsupported or malformed input returns an error.

The returned `interpret.Result[float64]` contains the parsed value, detected format, confidence, and original input.
