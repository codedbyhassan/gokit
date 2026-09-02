// Package number interprets human-friendly numeric input into float64 values.
package number

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/codedbyhassan/gokit/interpret"
)

var numeric = regexp.MustCompile(`^[+-]?(?:\d+(?:\.\d+)?|\.\d+)$`)
var scaleWords = map[string]float64{"thousand": 1e3, "million": 1e6, "billion": 1e9, "trillion": 1e12, "quadrillion": 1e15}
var wordNumbers = map[string]float64{"zero": 0, "one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9, "ten": 10, "eleven": 11, "twelve": 12, "thirteen": 13, "fourteen": 14, "fifteen": 15, "sixteen": 16, "seventeen": 17, "eighteen": 18, "nineteen": 19, "twenty": 20, "thirty": 30, "forty": 40, "fifty": 50, "sixty": 60, "seventy": 70, "eighty": 80, "ninety": 90}

func Parse(input string) (interpret.Result[float64], error) {
	original := input
	clean := normalize(input)
	if clean == "" { return interpret.Result[float64]{OriginalInput: original}, interpret.ErrEmptyInput }
	if numeric.MatchString(clean) {
		v, err := strconv.ParseFloat(clean, 64)
		if err == nil && !math.IsNaN(v) && !math.IsInf(v, 0) { return result(v, "numeric", interpret.HighConfidence, original), nil }
	}
	if v, ok := parseSuffixed(clean); ok { return result(v, "scaled numeric", .98, original), nil }
	if v, format, ok := parseWords(clean); ok { return result(v, format, .96, original), nil }
	return interpret.Result[float64]{OriginalInput: original}, fmt.Errorf("%w: %q", interpret.ErrUnrecognizedInput, input)
}

func normalize(input string) string {
	s := strings.ToLower(strings.TrimSpace(input))
	s = strings.ReplaceAll(s, "_", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "–", "-")
	s = strings.ReplaceAll(s, "—", "-")
	if strings.ContainsAny(s, "0123456789") { compact := strings.ReplaceAll(s, " ", ""); if numeric.MatchString(compact) { s = compact } }
	return strings.Join(strings.Fields(s), " ")
}

func result(v float64, f string, c interpret.Confidence, o string) interpret.Result[float64] { return interpret.Result[float64]{Value:v,Format:f,Confidence:c,OriginalInput:o} }

func parseSuffixed(input string) (float64, bool) {
	negative:=false
	if strings.HasPrefix(input,"-"){negative=true;input=strings.TrimSpace(strings.TrimPrefix(input,"-"))}else if strings.HasPrefix(input,"+"){input=strings.TrimSpace(strings.TrimPrefix(input,"+"))}
	for _,x:=range []struct{s string;m float64}{{"k",1e3},{"m",1e6},{"b",1e9},{"t",1e12}} { if strings.HasSuffix(input,x.s){n:=strings.TrimSpace(strings.TrimSuffix(input,x.s));if numeric.MatchString(n){v,err:=strconv.ParseFloat(n,64);if err==nil{if negative{v=-v};return v*x.m,true}}} }
	parts:=strings.Fields(input)
	if len(parts)>=2 { if scale,ok:=scaleWords[parts[len(parts)-1]];ok { base:=strings.Join(parts[:len(parts)-1]," ");if v,_,ok:=parseWords(base);ok{if negative{v=-v};return v*scale,true} } }
	return 0,false
}

func parseWords(input string)(float64,string,bool){
	parts:=strings.Fields(strings.ReplaceAll(input,"-"," "));if len(parts)==0{return 0,"",false}
	negative:=false;if parts[0]=="minus"||parts[0]=="negative"{negative=true;parts=parts[1:]};if len(parts)==0{return 0,"",false}
	if point:=indexOfPoint(parts);point>=0 { if point==0||point==len(parts)-1{return 0,"",false};integer,ok:=parseIntegerWords(parts[:point]);if !ok{return 0,"",false};fraction:="";for _,token:=range parts[point+1:]{if value,ok:=wordNumbers[token];ok&&value>=0&&value<=9{fraction+=strconv.Itoa(int(value));continue};if numeric.MatchString(token)&&!strings.Contains(token,"."){fraction+=token;continue};return 0,"",false};f,err:=strconv.ParseFloat("0."+fraction,64);if err!=nil{return 0,"",false};v:=integer+f;if negative{v=-v};return v,"number words",true }
	value,ok:=parseIntegerWords(parts);if !ok{return 0,"",false};if negative{value=-value};return value,"number words",true
}

func indexOfPoint(parts []string)int{for i,token:=range parts{if token=="point"||token=="dot"{return i}};return -1}

func parseIntegerWords(parts []string)(float64,bool){
	if len(parts)==0{return 0,false}
	var total,group float64
	used:=false
	lastScale:=math.Inf(1)
	lastValue:=-1.0
	lastWasTens:=false
	hundredUsed:=false
	for _,token:=range parts{
		if token=="and"{continue}
		if value,ok:=wordNumbers[token];ok{if lastValue>=0&&!lastWasTens{return 0,false};if lastValue>=0&&lastWasTens&&value>9{return 0,false};group+=value;lastValue=value;lastWasTens=value>=20&&value<=90&&math.Mod(value,10)==0;used=true;continue}
		if token=="hundred"{if group<=0||hundredUsed{return 0,false};group*=100;hundredUsed=true;lastValue=-1;lastWasTens=false;used=true;continue}
		if scale,ok:=scaleWords[token];ok{if group<=0||scale>=lastScale{return 0,false};total+=group*scale;group=0;hundredUsed=false;lastScale=scale;lastValue=-1;lastWasTens=false;used=true;continue}
		return 0,false
	}
	if !used{return 0,false};return total+group,true
}

func parseAtomicNumber(s string)(float64,error){if numeric.MatchString(s){return strconv.ParseFloat(s,64)};if v,_,ok:=parseWords(s);ok{return v,nil};return 0,fmt.Errorf("not a number")}
