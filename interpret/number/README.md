# Number Interpreter

The number interpreter converts common human-friendly numeric input into `float64` values.

Supported forms include:

- `1000`
- `1,000`
- `1 000`
- `1.5`
- `1.5k`
- `5k`
- `2 million`
- `2.5 million`
- `five`
- `twenty five`
- `one hundred`

The parser returns a typed value together with GoKit interpretation metadata such as confidence and the detected format.
