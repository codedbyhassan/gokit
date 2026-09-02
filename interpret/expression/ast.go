package expression

// Node is an expression tree node.
type Node interface{ node() }

// Number is a numeric literal.
type Number struct { Value float64 }
func (Number) node() {}

// Binary is a binary arithmetic operation.
type Binary struct {
	Left  Node
	Op    Operator
	Right Node
}
func (Binary) node() {}

// Unary is a signed numeric expression.
type Unary struct {
	Op    Operator
	Value Node
}
func (Unary) node() {}
