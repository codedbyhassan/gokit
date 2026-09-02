package expression

import (
	"fmt"

	"github.com/codedbyhassan/gokit/interpret/unit"
	"github.com/codedbyhassan/gokit/units"
)

type QuantityNode struct { Quantity unit.Quantity }
func (QuantityNode) node() {}

type QuantityResult struct { Value float64; Unit units.Unit }

func EvaluateQuantity(node Node) (QuantityResult, error) {
	switch n := node.(type) {
	case QuantityNode:
		return QuantityResult{Value:n.Quantity.Value, Unit:n.Quantity.Unit}, nil
	case Number:
		return QuantityResult{Value:n.Value}, nil
	case Unary:
		v,err:=EvaluateQuantity(n.Value);if err!=nil{return QuantityResult{},err};if n.Op==Subtract{v.Value=-v.Value};return v,nil
	case Binary:
		left,err:=EvaluateQuantity(n.Left);if err!=nil{return QuantityResult{},err}
		right,err:=EvaluateQuantity(n.Right);if err!=nil{return QuantityResult{},err}
		lh:=left.Unit.Name!="";rh:=right.Unit.Name!=""
		switch n.Op {
		case Add,Subtract:
			if !lh&&!rh { if n.Op==Add{return QuantityResult{Value:left.Value+right.Value},nil};return QuantityResult{Value:left.Value-right.Value},nil }
			if !lh||!rh{return QuantityResult{},fmt.Errorf("%s requires two quantities",n.Op)}
			if left.Unit.Dimension!=right.Unit.Dimension{return QuantityResult{},fmt.Errorf("cannot %s %s and %s",n.Op,left.Unit.Name,right.Unit.Name)}
			converted,err:=units.Convert(right.Value,right.Unit,left.Unit);if err!=nil{return QuantityResult{},err};if n.Op==Add{left.Value+=converted}else{left.Value-=converted};return left,nil
		case Multiply:
			if lh&&rh{return QuantityResult{},fmt.Errorf("multiplication of two quantities is not supported")}
			if lh{left.Value*=right.Value;return left,nil};if rh{right.Value*=left.Value;return right,nil};return QuantityResult{Value:left.Value*right.Value},nil
		case Divide:
			if right.Value==0{return QuantityResult{},fmt.Errorf("division by zero")};if lh&&rh{return QuantityResult{},fmt.Errorf("division of two quantities is not supported")};if lh{left.Value/=right.Value;return left,nil};return QuantityResult{Value:left.Value/right.Value},nil
		default:return QuantityResult{},fmt.Errorf("operator %q is not supported for quantities",n.Op)
		}
	default:return QuantityResult{},fmt.Errorf("unsupported quantity AST node %T",node)
	}
}
