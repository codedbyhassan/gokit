package unit

import (
	"math"
	"testing"
)

func TestWordQuantities(t *testing.T){
	tests:=[]struct{input string;want float64;symbol string}{
		{"five kilograms",5,"kg"},{"two hundred grams",200,"g"},{"1.25 kg",1.25,"kg"},{"three miles",3,"mi"},{"twenty five celsius",25,"°C"},
	}
	for _,tt:=range tests{r,err:=Parse(tt.input);if err!=nil{t.Errorf("%q: %v",tt.input,err);continue};q:=r.Value.(Quantity);if math.Abs(q.Value-tt.want)>1e-9||q.Unit.Symbol!=tt.symbol{t.Errorf("%q: got %v %s",tt.input,q.Value,q.Unit.Symbol)}}
}

func TestWordConversion(t *testing.T){r,err:=Parse("convert five kilometers to miles");if err!=nil{t.Fatal(err)};c:=r.Value.(Conversion);if math.Abs(c.Value-3.106855961)>1e-8{t.Fatalf("got %v",c.Value)}}

func TestTemperatureAliases(t *testing.T){for _,input:=range []string{"25°C","25°F","25 celsius","25 fahrenheit","273.15 kelvin"}{if _,err:=Parse(input);err!=nil{t.Errorf("%q: %v",input,err)}}}
