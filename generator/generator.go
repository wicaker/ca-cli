/*
Package generator used to generate codes and create a file based on the layer
*/
package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/wicaker/cacli/domain"
)

type caGen struct {
	gen domain.Generator
}

// NewGeneratorService will create new a caGen object representation of domain.Generator interface
func NewGeneratorService() domain.GeneratorService {
	return &caGen{}
}

func genParamList(i domain.Method) []jen.Code {
	var param []jen.Code
	for _, j := range i.ParameterList {
		param = append(param, jen.Id(j.Name).Op(j.Type))
	}
	return param
}

func genReturnList(i domain.Method) (returnType []jen.Code, returnValue []jen.Code) {
	for _, k := range i.ResultList {
		returnType = append(returnType, jen.Id(k.Name).Op(k.Type))
		switch k.Type {
		case "string":
			returnValue = append(returnValue, jen.Op(`""`))
		case "bool":
			returnValue = append(returnValue, jen.Op("false"))
		case "float32", "float64", "complex64", "complex128":
			returnValue = append(returnValue, jen.Op("0.0"))
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintpr", "byte", "rune":
			returnValue = append(returnValue, jen.Op("0"))
		default:
			if len(k.Type) > 7 {
				if k.Type[:7] == "domain." {
					returnValue = append(returnValue, jen.Op(k.Type+"{}"))
				} else {
					returnValue = append(returnValue, jen.Nil())
				}
			} else {
				returnValue = append(returnValue, jen.Nil())
			}
		}
	}
	return
}
