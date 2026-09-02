package expression

import (
	"fmt"
	"strings"

	"github.com/codedbyhassan/gokit/interpret/unit"
)

// ParseQuantityAST parses expressions containing measurements and scalar values.
// Quantity literals are recognized by trying the longest suffix at each factor.
func ParseQuantityAST(input string) (Node, error) {
	clean := normalize(input)
	clean = strings.TrimSpace(strings.TrimPrefix(clean, "what is "))
	clean = strings.TrimSuffix(clean, "?")
	if clean=="" {return nil,fmt.Errorf("empty expression")}
	p:=&quantityParser{input:clean}
	n,err:=p.parseExpression();if err!=nil{return nil,err}
	p.skipSpace();if p.pos!=len(p.input){return nil,fmt.Errorf("unexpected input near %q",p.input[p.pos:])}
	return n,nil
}

type quantityParser struct{input string;pos int}
func(p *quantityParser)skipSpace(){for p.pos<len(p.input)&&p.input[p.pos]==' '{p.pos++}}
func(p *quantityParser)parseExpression()(Node,error){left,e:=p.parseTerm();if e!=nil{return nil,e};for{p.skipSpace();if p.pos>=len(p.input){break};c:=p.input[p.pos];if c!='+'&&c!='-'{break};p.pos++;right,e:=p.parseTerm();if e!=nil{return nil,e};op:=Add;if c=='-'{op=Subtract};left=Binary{Left:left,Op:op,Right:right}};return left,nil}
func(p *quantityParser)parseTerm()(Node,error){left,e:=p.parseFactor();if e!=nil{return nil,e};for{p.skipSpace();if p.pos>=len(p.input){break};c:=p.input[p.pos];if c!='*'&&c!='/'{break};p.pos++;right,e:=p.parseFactor();if e!=nil{return nil,e};op:=Multiply;if c=='/'{op=Divide};left=Binary{Left:left,Op:op,Right:right}};return left,nil}
func(p *quantityParser)parseFactor()(Node,error){p.skipSpace();if p.pos>=len(p.input){return nil,fmt.Errorf("expected a value")};if p.input[p.pos]=='(' {p.pos++;n,e:=p.parseExpression();if e!=nil{return nil,e};p.skipSpace();if p.pos>=len(p.input)||p.input[p.pos]!=')'{return nil,fmt.Errorf("missing closing parenthesis")};p.pos++;return n,nil};start:=p.pos;if p.input[p.pos]=='+'||p.input[p.pos]=='-' {p.pos++};for p.pos<len(p.input)&&((p.input[p.pos]>='0'&&p.input[p.pos]<='9')||p.input[p.pos]=='.'||p.input[p.pos]==','||p.input[p.pos]=='_'){p.pos++};if start==p.pos{return nil,fmt.Errorf("expected a value near %q",p.input[p.pos:])};numberText:=p.input[start:p.pos];unitStart:=p.pos;p.skipSpace();for p.pos<len(p.input)&&isUnitChar(p.input[p.pos]){p.pos++};if p.pos>unitStart {q,err:=unit.Parse(numberText+" "+strings.TrimSpace(p.input[unitStart:p.pos]));if err==nil{return QuantityNode{Quantity:q.Value.(unit.Quantity)},nil}};p.pos=start;n,e:=ParseAST(p.input[p.pos:]);if e!=nil{return nil,e};return n,nil}
func isUnitChar(c byte)bool{return (c>='a'&&c<='z')||c=='°'||c=='/' }
