package parser

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/wicaker/cacli/domain"
)

type caParser struct {
	domainFileName string
	par            domain.Parser
}

// NewParserService /
func NewParserService(dFName string) domain.ParserService {
	return &caParser{
		domainFileName: dFName,
	}
}

func (v *caParser) Parser(filePath string) (*domain.Parser, error) {
	fs := token.NewFileSet()

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	f := fs.AddFile(filePath, fs.Base(), len(b))

	var s scanner.Scanner
	s.Init(f, b, nil, scanner.ScanComments)

	p, err := parser.ParseFile(fs, filePath, nil, parser.AllErrors)
	if err != nil {
		log.Printf("could not parse %s: %v", filePath, err)
		return nil, err
	}
	ast.Walk(v, p)
	if v.par.Repository.Name == "" || v.par.Usecase.Name == "" {
		return nil, errors.New("interface must exported type and interface name must same with name of file in domain dir")
	}
	return &v.par, nil
}

func (v *caParser) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.File:
		for _, j := range d.Decls {
			v.getType(j)
		}
	}
	return v
}

func (v *caParser) getType(n ast.Node) {
	switch d := n.(type) {
	case *ast.GenDecl:
		if d.Tok.String() == "type" {
			for _, j := range d.Specs {
				v.getTypeSpec(j)
			}
		}
	}
}

func (v *caParser) getTypeSpec(n ast.Node) {
	switch d := n.(type) {
	case *ast.TypeSpec:
		if _, ok := d.Type.(*ast.InterfaceType); ok {
			if len(d.Name.Name) > 7 {
				if string(d.Name.Name[len(d.Name.Name)-7:len(d.Name.Name)]) == "Usecase" {
					name := string(d.Name.Name[0])
					if strings.ToUpper(name) != string(d.Name.Name[0]) {
						return
					}
					if strings.ToUpper(d.Name.Name[:len(d.Name.Name)-7]) != strings.ToUpper(v.domainFileName) {
						return
					}
					v.par.Usecase.Name = d.Name.Name
					field := v.getMethods(d.Type)

					for _, j := range field {
						method := domain.Method{
							Name:          j.Names[0].String(),
							ParameterList: []domain.MethodValue{},
							ResultList:    []domain.MethodValue{},
						}
						fType := v.getFuncType(j.Type)
						for _, k := range fType.Params.List {
							r := v.getTypeExpr(k.Type, d.Name.Name[:len(d.Name.Name)-7])
							methodValue := domain.MethodValue{Type: r}
							if len(k.Names) > 0 {
								methodValue.Name = k.Names[0].String()
							}
							method.ParameterList = append(method.ParameterList, methodValue)
						}
						for _, k := range fType.Results.List {
							r := v.getTypeExpr(k.Type, d.Name.Name[:len(d.Name.Name)-7])
							methodValue := domain.MethodValue{Type: r}
							if len(k.Names) > 0 {
								methodValue.Name = k.Names[0].String()
							}
							method.ResultList = append(method.ResultList, methodValue)
						}
						v.par.Usecase.Method = append(v.par.Usecase.Method, method)
					}

					return
				}
			}
			if len(d.Name.Name) > 10 {
				if string(d.Name.Name[len(d.Name.Name)-10:len(d.Name.Name)]) == "Repository" {
					name := string(d.Name.Name[0])
					if strings.ToUpper(name) != string(d.Name.Name[0]) {
						return
					}
					if strings.ToUpper(d.Name.Name[:len(d.Name.Name)-10]) != strings.ToUpper(v.domainFileName) {
						return
					}
					v.par.Repository.Name = d.Name.Name
					field := v.getMethods(d.Type)

					for _, j := range field {
						method := domain.Method{
							Name:          j.Names[0].String(),
							ParameterList: []domain.MethodValue{},
							ResultList:    []domain.MethodValue{},
						}
						fType := v.getFuncType(j.Type)
						for _, k := range fType.Params.List {
							r := v.getTypeExpr(k.Type, d.Name.Name[:len(d.Name.Name)-10])
							methodValue := domain.MethodValue{Type: r}
							if len(k.Names) > 0 {
								methodValue.Name = k.Names[0].String()
							}
							method.ParameterList = append(method.ParameterList, methodValue)
						}
						for _, k := range fType.Results.List {
							r := v.getTypeExpr(k.Type, d.Name.Name[:len(d.Name.Name)-10])
							methodValue := domain.MethodValue{Type: r}
							if len(k.Names) > 0 {
								methodValue.Name = k.Names[0].String()
							}
							method.ResultList = append(method.ResultList, methodValue)
						}
						v.par.Repository.Method = append(v.par.Repository.Method, method)
					}

					return
				}
			}
		}
	}
}

func (v *caParser) getMethods(n ast.Node) []*ast.Field {
	switch d := n.(type) {
	case *ast.InterfaceType:
		return d.Methods.List
	default:
		return nil
	}
}

func (v *caParser) getFuncType(n ast.Node) *ast.FuncType {
	switch d := n.(type) {
	case *ast.FuncType:
		return d
	default:
		return nil
	}
}

func (v *caParser) getTypeExpr(n ast.Node, structName string) string {
	switch d := n.(type) {
	case *ast.StarExpr:
		return "*" + v.getTypeExpr(d.X, structName)
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
		return "[]" + v.getTypeExpr(d.Elt, structName)
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return ""
	}
}
