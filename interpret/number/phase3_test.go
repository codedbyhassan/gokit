package number

import "testing"

func TestPhase3CompoundNumbers(t *testing.T) {
	tests:=[]struct{input string;want float64}{
		{"twenty five",25},{"one hundred and twenty five",125},{"three hundred forty two",342},
		{"one thousand two hundred and thirty four",1234},{"two thousand and nineteen",2019},
		{"one million two hundred thousand",1200000},{"two million five hundred thousand",2500000},
		{"two point five million",2500000},{"one hundred twenty five point five",125.5},
		{"negative one hundred",-100},{"minus one hundred",-100},
	}
	for _,tt:=range tests { r,e:=Parse(tt.input);if e!=nil{t.Errorf("%q: %v",tt.input,e);continue};if r.Value!=tt.want{t.Errorf("%q: got %v want %v",tt.input,r.Value,tt.want)} }
}

func TestPhase3ScaledNumbers(t *testing.T){
	for _,tt:=range []struct{input string;want float64}{{"1.5k",1500},{"2.5m",2500000},{"3b",3e9},{"4 trillion",4e12},{"2.5 million",2500000}}{r,e:=Parse(tt.input);if e!=nil||r.Value!=tt.want{t.Errorf("%q: got %v err=%v",tt.input,r.Value,e)}}
}

func TestPhase3RejectsMalformedWords(t *testing.T){for _,s:=range []string{"one two hundred hundred","million five","one point","not a number"}{if _,e:=Parse(s);e==nil{t.Errorf("expected rejection for %q",s)}}}
