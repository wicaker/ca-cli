/*
Package parser used to parsing golang file and get specific information
*/
package parser

import (
	"fmt"
	"go/ast"
	"strings"
)

func getMethodsInInterface(n ast.Node) []*ast.Field {
	switch d := n.(type) {
	case *ast.InterfaceType:
		return d.Methods.List
	default:
		return nil
	}
}

func getFuncType(n ast.Node) *ast.FuncType {
	switch d := n.(type) {
	case *ast.FuncType:
		return d
	default:
		return nil
	}
}

func getTypeExpr(n ast.Node, structName string) string {
	switch d := n.(type) {
	case *ast.StarExpr:
		return "*" + getTypeExpr(d.X, structName)
	case *ast.SelectorExpr:
		result := fmt.Sprintf("%s.%s", d.X, d.Sel)
		return result
	case *ast.Ident:
		name := d.Name
		if string(d.Name[0]) == strings.ToUpper(string(name[0])) {
			return "domain." + d.Name
		}
		return d.Name
	case *ast.ArrayType:
		return "[]" + getTypeExpr(d.Elt, structName)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return ""
	}
}
