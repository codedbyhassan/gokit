package expression

import (
	"fmt"
	"strings"

	"github.com/codedbyhassan/gokit/interpret/number"
	"github.com/codedbyhassan/gokit/interpret/unit"
)

// ParseQuantityAST parses arithmetic containing quantities and scalar values.
func ParseQuantityAST(input string) (Node,error){clean:=normalizeQuantityExpression(input);clean=strings.TrimSpace(strings.TrimPrefix(clean,"what is "));clean=strings.TrimSuffix(clean,"?");if clean==""{return nil,fmt.Errorf("empty expression")};p:=&quantityParser{input:clean};n,err:=p.parseExpression();if err!=nil{return nil,err};p.skipSpace();if p.pos!=len(p.input){return nil,fmt.Errorf("unexpected input near %q",p.input[p.pos:])};return n,nil}

func normalizeQuantityExpression(input string)string{clean:=strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(input)))," ");clean=strings.TrimSuffix(clean,"?");clean=strings.NewReplacer("multiplied by","*","divided by","/","plus","+","added to","+","add to","+","minus","-","subtract from","-","times","*").Replace(clean);return strings.Join(strings.Fields(clean)," ")}

type quantityParser struct{input string;pos int}
func(p *quantityParser)skipSpace(){for p.pos<len(p.input)&&p.input[p.pos]==' '{p.pos++}}
func(p *quantityParser)parseExpression()(Node,error){left,err:=p.parseTerm();if err!=nil{return nil,err};for{p.skipSpace();if p.pos>=len(p.input)||(p.input[p.pos]!='+'&&p.input[p.pos]!='-'){break};op:=Add;if p.input[p.pos]=='-'{op=Subtract};p.pos++;right,err:=p.parseTerm();if err!=nil{return nil,err};left=Binary{Left:left,Op:op,Right:right}};return left,nil}
func(p *quantityParser)parseTerm()(Node,error){left,err:=p.parseFactor();if err!=nil{return nil,err};for{p.skipSpace();if p.pos>=len(p.input)||(p.input[p.pos]!='*'&&p.input[p.pos]!='/'){break};op:=Multiply;if p.input[p.pos]=='/'{op=Divide};p.pos++;right,err:=p.parseFactor();if err!=nil{return nil,err};left=Binary{Left:left,Op:op,Right:right}};return left,nil}
func(p *quantityParser)parseFactor()(Node,error){p.skipSpace();if p.pos>=len(p.input){return nil,fmt.Errorf("expected a value")};if p.input[p.pos]=='(' {p.pos++;n,err:=p.parseExpression();if err!=nil{return nil,err};p.skipSpace();if p.pos>=len(p.input)||p.input[p.pos]!=')'{return nil,fmt.Errorf("missing closing parenthesis")};p.pos++;return n,nil};if p.input[p.pos]=='+'||p.input[p.pos]=='-'{op:=p.input[p.pos];p.pos++;n,err:=p.parseFactor();if err!=nil{return nil,err};if op=='-'{return Unary{Op:Subtract,Value:n},nil};return n,nil}
	start:=p.pos;for p.pos<len(p.input)&&p.input[p.pos]!='+'&&p.input[p.pos]!='-'&&p.input[p.pos]!='*'&&p.input[p.pos]!='/'&&p.input[p.pos]!='('&&p.input[p.pos]!=')'{p.pos++};text:=strings.TrimSpace(p.input[start:p.pos]);if text==""{return nil,fmt.Errorf("expected a value")}
	if q,err:=unit.ParseQuantity(text);err==nil{return QuantityNode{Quantity:q},nil};if parsed,err:=number.Parse(text);err==nil{return Number{Value:parsed.Value},nil};return nil,fmt.Errorf("invalid expression value %q",text)
}
