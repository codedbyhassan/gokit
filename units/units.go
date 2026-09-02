// Package units provides reusable unit definitions and conversions.
package units

import (
	"fmt"
	"sort"
	"strings"
)

type Dimension string
const (
	Length Dimension = "length"
	Mass Dimension = "mass"
	Temperature Dimension = "temperature"
	Speed Dimension = "speed"
	Volume Dimension = "volume"
)

type Unit struct { Name string; Symbol string; Dimension Dimension }

var definitions = map[string]Unit{
	"meter":{"meter","m",Length}, "m":{"meter","m",Length}, "meters":{"meter","m",Length},
	"centimeter":{"centimeter","cm",Length}, "cm":{"centimeter","cm",Length}, "centimeters":{"centimeter","cm",Length},
	"millimeter":{"millimeter","mm",Length}, "mm":{"millimeter","mm",Length}, "millimeters":{"millimeter","mm",Length},
	"kilometer":{"kilometer","km",Length}, "km":{"kilometer","km",Length}, "kilometers":{"kilometer","km",Length},
	"mile":{"mile","mi",Length}, "mi":{"mile","mi",Length}, "miles":{"mile","mi",Length},
	"yard":{"yard","yd",Length}, "yd":{"yard","yd",Length}, "yards":{"yard","yd",Length},
	"foot":{"foot","ft",Length}, "ft":{"foot","ft",Length}, "feet":{"foot","ft",Length},
	"inch":{"inch","in",Length}, "in":{"inch","in",Length}, "inches":{"inch","in",Length},
	"kilogram":{"kilogram","kg",Mass}, "kg":{"kilogram","kg",Mass}, "kilograms":{"kilogram","kg",Mass},
	"gram":{"gram","g",Mass}, "g":{"gram","g",Mass}, "grams":{"gram","g",Mass},
	"milligram":{"milligram","mg",Mass}, "mg":{"milligram","mg",Mass}, "milligrams":{"milligram","mg",Mass},
	"tonne":{"tonne","t",Mass}, "tonnes":{"tonne","t",Mass}, "metric ton":{"tonne","t",Mass},
	"pound":{"pound","lb",Mass}, "lb":{"pound","lb",Mass}, "lbs":{"pound","lb",Mass}, "pounds":{"pound","lb",Mass},
	"ounce":{"ounce","oz",Mass}, "oz":{"ounce","oz",Mass}, "ounces":{"ounce","oz",Mass},
	"celsius":{"celsius","°C",Temperature}, "centigrade":{"celsius","°C",Temperature}, "c":{"celsius","°C",Temperature}, "°c":{"celsius","°C",Temperature},
	"fahrenheit":{"fahrenheit","°F",Temperature}, "f":{"fahrenheit","°F",Temperature}, "°f":{"fahrenheit","°F",Temperature},
	"kelvin":{"kelvin","K",Temperature}, "k":{"kelvin","K",Temperature},
	"kilometers per hour":{"kilometers per hour","km/h",Speed}, "km/h":{"kilometers per hour","km/h",Speed},
	"miles per hour":{"miles per hour","mph",Speed}, "mph":{"miles per hour","mph",Speed},
	"meter per second":{"meter per second","m/s",Speed}, "meters per second":{"meter per second","m/s",Speed}, "m/s":{"meter per second","m/s",Speed},
	"liter":{"liter","L",Volume}, "l":{"liter","L",Volume}, "liters":{"liter","L",Volume},
	"milliliter":{"milliliter","mL",Volume}, "ml":{"milliliter","mL",Volume}, "milliliters":{"milliliter","mL",Volume},
	"gallon":{"gallon","gal",Volume}, "gal":{"gallon","gal",Volume}, "gallons":{"gallon","gal",Volume},
}

func Lookup(input string) (Unit,error) { key:=strings.ToLower(strings.TrimSpace(input));u,ok:=definitions[key];if !ok{return Unit{},fmt.Errorf("unsupported unit %q",input)};return u,nil }

// Aliases returns known aliases longest-first, useful to deterministic parsers.
func Aliases() []string { out:=make([]string,0,len(definitions));for k:=range definitions{out=append(out,k)};sort.Slice(out,func(i,j int)bool{if len(out[i])==len(out[j]){return out[i]<out[j]};return len(out[i])>len(out[j])});return out }

func Convert(value float64,from,to Unit)(float64,error){if from.Dimension!=to.Dimension{return 0,fmt.Errorf("cannot convert %s to %s",from.Name,to.Name)};if from.Name==to.Name{return value,nil};if from.Dimension==Temperature{return convertTemperature(value,from,to),nil};return fromBase(toBase(value,from),to),nil}
func toBase(v float64,u Unit)float64{switch u.Name{case "kilometer":return v*1000;case "centimeter":return v*.01;case "millimeter":return v*.001;case "mile":return v*1609.344;case "yard":return v*.9144;case "foot":return v*.3048;case "inch":return v*.0254;case "kilogram":return v*1000;case "milligram":return v*.001;case "tonne":return v*1000000;case "pound":return v*453.59237;case "ounce":return v*28.349523125;case "milliliter":return v*.001;case "gallon":return v*3.785411784;case "miles per hour":return v*.44704;case "kilometers per hour":return v/3.6;case "meter per second":return v;default:return v}}
func fromBase(v float64,u Unit)float64{switch u.Name{case "kilometer":return v/1000;case "centimeter":return v/.01;case "millimeter":return v/.001;case "mile":return v/1609.344;case "yard":return v/.9144;case "foot":return v/.3048;case "inch":return v/.0254;case "kilogram":return v/1000;case "milligram":return v/.001;case "tonne":return v/1000000;case "pound":return v/453.59237;case "ounce":return v/28.349523125;case "milliliter":return v/.001;case "gallon":return v/3.785411784;case "miles per hour":return v/.44704;case "kilometers per hour":return v*3.6;case "meter per second":return v;default:return v}}
func convertTemperature(v float64,from,to Unit)float64{var c float64;switch from.Name{case "celsius":c=v;case "fahrenheit":c=(v-32)*5/9;case "kelvin":c=v-273.15};switch to.Name{case "celsius":return c;case "fahrenheit":return c*9/5+32;case "kelvin":return c+273.15};return c}
