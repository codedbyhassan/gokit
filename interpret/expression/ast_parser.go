package expression

import (
	"fmt"
	"strings"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/number"
)

// ParseAST parses arithmetic with standard precedence and parentheses.
func ParseAST(input string) (Node, error) {
	clean := normalize(input)
	clean = strings.TrimSpace(strings.TrimPrefix(clean, "what is "))
	clean = strings.TrimSuffix(clean, "?")
	if clean == "" { return nil, interpret.ErrEmptyInput }
	p := &astParser{input: clean}
	n, err := p.parseExpression()
	if err != nil { return nil, err }
	p.skipSpace()
	if p.pos != len(p.input) { return nil, fmt.Errorf("unexpected input near %q", p.input[p.pos:]) }
	return n, nil
}

type astParser struct { input string; pos int }
func (p *astParser) skipSpace() { for p.pos < len(p.input) && p.input[p.pos] == ' ' { p.pos++ } }
func (p *astParser) parseExpression() (Node,error) {
	left,err:=p.parseTerm(); if err!=nil{return nil,err}
	for { p.skipSpace(); if p.pos>=len(p.input){break}; c:=p.input[p.pos]; if c!='+'&&c!='-'{break}; p.pos++; right,e:=p.parseTerm(); if e!=nil{return nil,e}; op:=Add;if c=='-'{op=Subtract};left=Binary{Left:left,Op:op,Right:right} }
	return left,nil
}
func (p *astParser) parseTerm() (Node,error) {
	left,err:=p.parseFactor(); if err!=nil{return nil,err}
	for { p.skipSpace(); if p.pos>=len(p.input){break}; c:=p.input[p.pos]; if c!='*'&&c!='/'&&c!='%'{break}; p.pos++; right,e:=p.parseFactor();if e!=nil{return nil,e};op:=Multiply;if c=='/'{op=Divide};if c=='%'{op=Operator("modulo")};left=Binary{Left:left,Op:op,Right:right} }
	return left,nil
}
func (p *astParser) parseFactor() (Node,error) {
	p.skipSpace(); if p.pos>=len(p.input){return nil,fmt.Errorf("expected a number")}
	if p.input[p.pos]=='(' { p.pos++; n,e:=p.parseExpression();if e!=nil{return nil,e};p.skipSpace();if p.pos>=len(p.input)||p.input[p.pos]!=')'{return nil,fmt.Errorf("missing closing parenthesis")};p.pos++;return n,nil }
	if p.input[p.pos]=='+'||p.input[p.pos]=='-' { c:=p.input[p.pos];p.pos++;n,e:=p.parseFactor();if e!=nil{return nil,e};if c=='-' {return Unary{Op:Subtract,Value:n},nil};return n,nil }
	start:=p.pos
	for p.pos<len(p.input) && (p.input[p.pos]>='0'&&p.input[p.pos]<='9'||p.input[p.pos]=='.'||p.input[p.pos]==','||p.input[p.pos]=='_'){p.pos++}
	if start==p.pos{return nil,fmt.Errorf("expected a number near %q",p.input[p.pos:])}
	r,err:=number.Parse(p.input[start:p.pos]);if err!=nil{return nil,err};return Number{Value:r.Value},nil
}
