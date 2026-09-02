package expression

import (
	"math"
	"testing"
)

func TestNaturalWordQuantityExpressions(t *testing.T){
	tests:=[]struct{input string;want float64;symbol string}{
		{"five kilograms plus two hundred grams",5.2,"kg"},
		{"two hundred grams plus five kilograms",5200,"g"},
		{"three miles plus two kilometers",4.242742,"mi"},
		{"two kilograms times three",6,"kg"},
		{"ten kilograms divided by two",5,"kg"},
		{"(two meters plus three meters) times four",20,"m"},
	}
	for _,tt:=range tests{n,err:=ParseQuantityAST(tt.input);if err!=nil{t.Errorf("%q: %v",tt.input,err);continue};got,err:=EvaluateQuantity(n);if err!=nil{t.Errorf("%q: %v",tt.input,err);continue};if math.Abs(got.Value-tt.want)>1e-5||got.Unit.Symbol!=tt.symbol{t.Errorf("%q: got %v %s want %v %s",tt.input,got.Value,got.Unit.Symbol,tt.want,tt.symbol)}}
}

func TestQuantityOperatorRules(t *testing.T){
	for _,input:=range []string{"5kg * 2kg","5kg / 2kg","5kg + 2"}{n,err:=ParseQuantityAST(input);if err!=nil{t.Fatal(err)};if _,err=EvaluateQuantity(n);err==nil{t.Errorf("%q: expected semantic error",input)}}
}

func TestQuantityTemperatureConversion(t *testing.T){
	n,err:=ParseQuantityAST("0 celsius + 32 fahrenheit");if err!=nil{t.Fatal(err)}
	got,err:=EvaluateQuantity(n);if err!=nil{t.Fatal(err)}
	// Absolute temperatures are converted with their offsets before arithmetic:
	// 32°F is 0°C, so 0°C + 32°F = 0°C.
	if math.Abs(got.Value)>1e-8||got.Unit.Symbol!="°C"{t.Fatalf("got %v %s",got.Value,got.Unit.Symbol)}
}
