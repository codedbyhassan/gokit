package expression

import (
	"math"
	"testing"
)

func TestParseASTPrecedence(t *testing.T) {
	tests:=[]struct{input string;want float64}{
		{"10 + 5 * 4",30},
		{"(10 + 5) * 4",60},
		{"100 / 5 - 7",13},
		{"-10 + 4",-6},
		{"20 % 6",2},
	}
	for _,tt:=range tests{n,err:=ParseAST(tt.input);if err!=nil{t.Fatalf("%q: %v",tt.input,err)};got,err:=Evaluate(n);if err!=nil{t.Fatalf("%q: %v",tt.input,err)};if math.Abs(got-tt.want)>1e-9{t.Fatalf("%q: got %v want %v",tt.input,got,tt.want)}}
}

func TestParseASTInvalid(t *testing.T) {
	for _,input:=range []string{"", "10 +", "(10 + 5", "hello"}{
		if _,err:=ParseAST(input);err==nil{t.Fatalf("expected error for %q",input)}
	}
}
