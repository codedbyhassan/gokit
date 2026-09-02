// Package date interprets common human date formats into time.Time values.
package date

import (
    "fmt"
    "regexp"
    "sort"
    "strconv"
    "strings"
    "time"
    "github.com/codedbyhassan/gokit/interpret"
)
var numericDate=regexp.MustCompile(`^(\d{1,4})[-/.](\d{1,4})[-/.](\d{1,4})$`)
func Parse(input string)(interpret.Result[time.Time],error){original:=input;input=interpret.Normalize(input);if input==""{return interpret.Result[time.Time]{OriginalInput:original},interpret.ErrEmptyInput};if value,layout,ok:=parseNamed(input);ok{return interpret.Result[time.Time]{Value:value,Format:layout,Confidence:interpret.HighConfidence,OriginalInput:original},nil};parts:=numericDate.FindStringSubmatch(input);if parts==nil{return interpret.Result[time.Time]{OriginalInput:original},fmt.Errorf("%w: %q",interpret.ErrUnrecognizedInput,input)};n:=make([]int,3);for i:=0;i<3;i++{v,e:=strconv.Atoi(parts[i+1]);if e!=nil{return interpret.Result[time.Time]{OriginalInput:original},fmt.Errorf("%w: %q",interpret.ErrUnrecognizedInput,input)};n[i]=v};c:=buildCandidates(n);if len(c)==0{return interpret.Result[time.Time]{OriginalInput:original},fmt.Errorf("%w: %q",interpret.ErrUnrecognizedInput,input)};if len(c)>1&&c[0].score==c[1].score{r:=interpret.Result[time.Time]{Format:"ambiguous",Confidence:c[0].score,Ambiguous:true,OriginalInput:original};for _,x:=range c{r.Candidates=append(r.Candidates,interpret.Candidate{Format:x.format,Confidence:x.score,Reason:x.reason})};return r,interpret.ErrAmbiguousInput};best:=c[0];r:=interpret.Result[time.Time]{Value:best.value,Format:best.format,Confidence:best.score,OriginalInput:original};if best.assumption!=""{r.Assumptions=[]string{best.assumption}};for _,x:=range c{r.Candidates=append(r.Candidates,interpret.Candidate{Format:x.format,Confidence:x.score,Reason:x.reason})};return r,nil}
type candidate struct{value time.Time;format string;score interpret.Confidence;reason,assumption string}
func buildCandidates(n []int)[]candidate{perms:=[]struct{o [3]int;f,r string}{{[3]int{0,1,2},"DD-MM-YYYY","day-month-year"},{[3]int{1,0,2},"MM-DD-YYYY","month-day-year"},{[3]int{2,1,0},"YYYY-MM-DD","year-month-day"},{[3]int{0,2,1},"DD-YYYY-MM","day-year-month"},{[3]int{1,2,0},"MM-YYYY-DD","month-year-day"},{[3]int{2,0,1},"YYYY-DD-MM","year-day-month"}};var out []candidate;for _,p:=range perms{year,month,day:=n[p.o[2]],n[p.o[1]],n[p.o[0]];if year<1000||year>9999||month<1||month>12||day<1||day>31{continue};v:=time.Date(year,time.Month(month),day,0,0,0,0,time.UTC);if v.Year()!=year||int(v.Month())!=month||v.Day()!=day{continue};score:=interpret.Confidence(.70);assumption:="";if year>=1000{score=.95;assumption=fmt.Sprintf("%d identified as the year",year)};out=append(out,candidate{v,p.f,score,p.r,assumption})};sort.SliceStable(out,func(i,j int)bool{return out[i].score>out[j].score});return out}
func parseNamed(input string)(time.Time,string,bool){layouts:=[]struct{layout,name string}{{"2 January 2006","D Month YYYY"},{"2 Jan 2006","D Mon YYYY"},{"January 2 2006","Month D YYYY"},{"Jan 2 2006","Mon D YYYY"},{"2006 January 2","YYYY Month D"}};for _,x:=range layouts{if v,e:=time.Parse(x.layout,input);e==nil{return v,x.name,true}};return time.Time{},"",false}
