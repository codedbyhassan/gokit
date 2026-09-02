// Package unit interprets human-friendly measurements and conversions.
package unit

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/codedbyhassan/gokit/interpret"
	"github.com/codedbyhassan/gokit/interpret/number"
	"github.com/codedbyhassan/gokit/units"
)

type Quantity struct { Value float64; Unit units.Unit }
type Conversion struct { Value float64; From Quantity; To units.Unit; OriginalInput string }

var conversionPattern=regexp.MustCompile(`^(.+?)\s+(?:to|into|in)\s+(.+)$`)

// Parse interprets numeric or word-based quantities and natural conversion requests.
func Parse(input string)(interpret.Result[any],error){
	original:=input;clean:=normalize(input);if clean==""{return interpret.Result[any]{OriginalInput:original},interpret.ErrEmptyInput}
	for _,prefix:=range []string{"please convert ","convert ","calculate the conversion of ","what is ","calculate "}{if strings.HasPrefix(clean,prefix){clean=strings.TrimSpace(strings.TrimPrefix(clean,prefix));break}}
	clean=strings.TrimSuffix(strings.TrimSpace(clean),"?")
	if m:=conversionPattern.FindStringSubmatch(clean);len(m)==3{from,err:=parseQuantity(m[1]);if err!=nil{return interpret.Result[any]{OriginalInput:original},err};to,err:=units.Lookup(m[2]);if err!=nil{return interpret.Result[any]{OriginalInput:original},err};value,err:=units.Convert(from.Value,from.Unit,to);if err!=nil{return interpret.Result[any]{OriginalInput:original},err};return interpret.Result[any]{Value:Conversion{Value:value,From:from,To:to,OriginalInput:original},Format:"unit conversion",Confidence:.99,OriginalInput:original},nil}
	q,err:=parseQuantity(clean);if err!=nil{return interpret.Result[any]{OriginalInput:original},err};return interpret.Result[any]{Value:q,Format:"quantity",Confidence:.97,OriginalInput:original},nil
}

func normalize(input string)string{clean:=strings.ToLower(strings.TrimSpace(input));clean=strings.NewReplacer("°c"," celsius","°f"," fahrenheit").Replace(clean);return strings.Join(strings.Fields(clean)," ")}

func parseQuantity(input string)(Quantity,error){clean:=strings.TrimSpace(input);for _,alias:=range units.Aliases(){if !strings.HasSuffix(clean,alias){continue};numberText:=strings.TrimSpace(strings.TrimSuffix(clean,alias));if numberText==""{continue};value,err:=number.Parse(numberText);if err!=nil{continue};u,err:=units.Lookup(alias);if err==nil{return Quantity{Value:value.Value,Unit:u},nil}}
	return Quantity{},fmt.Errorf("invalid measurement %q",input)
}
