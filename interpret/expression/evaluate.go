package expression

import (
	"fmt"
	"math"

	"github.com/codedbyhassan/gokit/calculator"
)

// Evaluate executes an expression AST recursively.
func Evaluate(node Node) (float64, error) {
	switch n := node.(type) {
	case Number:
		return n.Value, nil
	case Unary:
		v, err := Evaluate(n.Value); if err != nil { return 0, err }
		if n.Op == Subtract { return -v, nil }
		return v, nil
	case Binary:
		left, err := Evaluate(n.Left); if err != nil { return 0, err }
		right, err := Evaluate(n.Right); if err != nil { return 0, err }
		switch n.Op {
		case Add: return calculator.AddNumbers(left,right),nil
		case Subtract: return calculator.SubtractNumbers(left,right),nil
		case Multiply: return calculator.MultiplyNumbers(left,right),nil
		case Divide:
			return calculator.DivideNumbers(left,right)
		case Operator("modulo"):
			if right==0{return 0,calculator.ErrDivisionByZero}; return math.Mod(left,right),nil
		default:return 0,fmt.Errorf("unsupported operator %q",n.Op)
		}
	default:return 0,fmt.Errorf("unsupported AST node %T",node)
	}
}
