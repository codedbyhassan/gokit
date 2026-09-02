package intent

import (
	"math"
	"testing"

	"github.com/codedbyhassan/gokit/interpret/expression"
)

func TestIntentNaturalQuantityExpressions(t *testing.T){
	tests:=[]struct{input string;want float64;symbol string}{
		{"calculate five kilograms plus two hundred grams",5.2,"kg"},
		{"what is ten miles plus five kilometers",13.106855,"mi"},
		{"calculate two liters times three",6,"L"},
	}
	for _,tt:=range tests{r,err:=ParseAt(tt.input,nowForTest());if err!=nil{t.Errorf("%q: %v",tt.input,err);continue};q,ok:=r.Value.(expression.QuantityResult);if !ok{t.Errorf("%q: unexpected type %T",tt.input,r.Value);continue};if math.Abs(q.Value-tt.want)>1e-5||q.Unit.Symbol!=tt.symbol{t.Errorf("%q: got %v %s",tt.input,q.Value,q.Unit.Symbol)}}
}

func nowForTest() (t time.Time) { return time.Date(2026,9,2,0,0,0,0,time.UTC) }
