package expression

import "testing"

func TestParseASTPublicCases(t *testing.T) {
	cases:=[]string{"10 + 5","10 - 3","10 * 5","100 / 4","10 + 5 * 2","(10 + 5) * 2","-10 + 4"}
	for _,input:=range cases{if _,err:=ParseAST(input);err!=nil{t.Errorf("%q: %v",input,err)}}
}
