package expression

import (
	"fmt"
	"strings"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/number"
)

// ParsePercentAST parses percent-of and percent-increase expressions into AST nodes.
func ParsePercentAST(input string) (Node, error) {
	clean := normalize(input)
	clean = strings.TrimSpace(strings.TrimPrefix(clean, "what is "))
	clean = strings.TrimSuffix(clean, "?")
	if clean == "" { return nil, interpret.ErrEmptyInput }
	if strings.Contains(clean, " of ") {
		parts := strings.SplitN(clean, " of ", 2)
		p := strings.TrimSuffix(strings.TrimSpace(parts[0]), "%")
		if r, e1 := number.Parse(p); e1 == nil {
			if b, e2 := number.Parse(parts[1]); e2 == nil && strings.HasSuffix(strings.TrimSpace(parts[0]), "%") {
				return Binary{Left:Number{Value:b.Value}, Op:Operator("percent_of"), Right:Number{Value:r.Value}}, nil
			}
		}
	}
	for _, marker := range []string{"% increase on ", " percent increase on "} {
		if parts := strings.SplitN(clean, marker, 2); len(parts)==2 {
			p,e1:=number.Parse(strings.TrimSpace(parts[0])); b,e2:=number.Parse(strings.TrimSpace(parts[1])); if e1==nil&&e2==nil{return Binary{Left:Number{Value:b.Value},Op:Operator("percent_increase"),Right:Number{Value:p.Value}},nil}
		}
	}
	return nil,fmt.Errorf("not a percent expression")
}
