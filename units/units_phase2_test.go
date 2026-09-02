package units

import (
	"math"
	"testing"
)

func TestPhase2Conversions(t *testing.T) {
	tests:=[]struct{value float64;from,to string;want float64}{
		{150,"cm","m",1.5},{2,"km","m",2000},{2,"kg","g",2000},{500,"mg","g",.5},{1,"tonne","kg",1000},{1,"gallon","liter",3.785411784},{1,"liter","ml",1000},{60,"mph","km/h",96.56064},{100,"celsius","fahrenheit",212},{273.15,"kelvin","celsius",0},
	}
	for _,tt:=range tests{from,e:=Lookup(tt.from);if e!=nil{t.Fatal(e)};to,e:=Lookup(tt.to);if e!=nil{t.Fatal(e)};got,e:=Convert(tt.value,from,to);if e!=nil{t.Fatal(e)};if math.Abs(got-tt.want)>1e-8{t.Errorf("%g %s -> %s: got %g want %g",tt.value,tt.from,tt.to,got,tt.want)}}
}

func TestIncompatibleDimensions(t *testing.T){a,_:=Lookup("kg");b,_:=Lookup("m");if _,err:=Convert(1,a,b);err==nil{t.Fatal("expected dimension error")}}

func TestUnknownUnit(t *testing.T){if _,err:=Lookup("banana");err==nil{t.Fatal("expected unknown unit error")}}
