package parser

import (
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

type repoParser struct {
	par domain.Parser
}

// NewParserRepoService /
func NewParserRepoService() domain.RepoParserService {
	return &repoParser{}
}

func (v *repoParser) RepoParser(filePath string) (*domain.Parser, error) {
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
	return &v.par, nil
}

func (v *repoParser) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.File:
		for _, j := range d.Decls {
			v.getFunc(j)
		}
	}
	return v
}

func (v *repoParser) getFunc(n ast.Node) {
	switch d := n.(type) {
	case *ast.FuncDecl:
		// spew.Dump(d)
		if len(d.Name.String()) > 7 {
			if d.Name.Name[:3] == "New" && d.Name.Name[len(d.Name.Name)-7:len(d.Name.Name)] == "Usecase" {
				// spew.Dump(d.Type)
				method := domain.Method{
					Name:          d.Name.String(),
					ParameterList: []domain.MethodValue{},
					ResultList:    []domain.MethodValue{},
				}
				fType := v.getFuncType(d.Type)
				for _, k := range fType.Params.List {
					// fmt.Println(k)
					r := v.getTypeExpr(k.Type, "")
					// fmt.Println(k.Names[0].String(), r)
					// r := v.getTypeExpr(k.Type, d.Name.Name[:len(d.Name.Name)-7])
					methodValue := domain.MethodValue{Type: r}
					if len(k.Names) > 0 {
						methodValue.Name = k.Names[0].String()
					}
					method.ParameterList = append(method.ParameterList, methodValue)
				}
				for _, k := range fType.Results.List {
					// fmt.Println(k)
					r := v.getTypeExpr(k.Type, "")

					methodValue := domain.MethodValue{Type: r}
					if len(k.Names) > 0 {
						methodValue.Name = k.Names[0].String()
					}
					method.ResultList = append(method.ResultList, methodValue)
				}
				v.par.Usecase.Method = append(v.par.Usecase.Method, method)
			}
			// fmt.Println(d.Name.Name[:3], d.Name.Name[len(d.Name.Name)-7:len(d.Name.Name)])
		}

		if len(d.Name.String()) > 10 {
			if d.Name.Name[:3] == "New" && d.Name.Name[len(d.Name.Name)-10:len(d.Name.Name)] == "Repository" {
				// fmt.Println(d.Type)
				method := domain.Method{
					Name:          d.Name.String(),
					ParameterList: []domain.MethodValue{},
					ResultList:    []domain.MethodValue{},
				}
				fType := v.getFuncType(d.Type)
				for _, k := range fType.Params.List {
					// fmt.Println(k)
					r := v.getTypeExpr(k.Type, "")
					// fmt.Println(k.Names[0].String(), r)
					// r := v.getTypeExpr(k.Type, d.Name.Name[:len(d.Name.Name)-7])
					methodValue := domain.MethodValue{Type: r}
					if len(k.Names) > 0 {
						methodValue.Name = k.Names[0].String()
					}
					method.ParameterList = append(method.ParameterList, methodValue)
				}
				for _, k := range fType.Results.List {
					// fmt.Println(k)
					r := v.getTypeExpr(k.Type, "")

					methodValue := domain.MethodValue{Type: r}
					if len(k.Names) > 0 {
						methodValue.Name = k.Names[0].String()
					}
					method.ResultList = append(method.ResultList, methodValue)
				}
				v.par.Repository.Method = append(v.par.Repository.Method, method)
			}
		}

		if len(d.Name.String()) > 7 {
			if d.Name.Name[:3] == "New" && d.Name.Name[len(d.Name.Name)-7:len(d.Name.Name)] == "Handler" {
				// spew.Dump(d.Type)
				method := domain.Method{
					Name:          d.Name.String(),
					ParameterList: []domain.MethodValue{},
					ResultList:    []domain.MethodValue{},
				}
				fType := v.getFuncType(d.Type)
				for _, k := range fType.Params.List {
					// fmt.Println(k)
					r := v.getTypeExpr(k.Type, "")
					// fmt.Println(k.Names[0].String(), r)
					// r := v.getTypeExpr(k.Type, d.Name.Name[:len(d.Name.Name)-7])
					methodValue := domain.MethodValue{Type: r}
					if len(k.Names) > 0 {
						methodValue.Name = k.Names[0].String()
					}
					method.ParameterList = append(method.ParameterList, methodValue)
				}
				if fType.Results != nil {
					for _, k := range fType.Results.List {
						r := v.getTypeExpr(k.Type, "")

						methodValue := domain.MethodValue{Type: r}
						if len(k.Names) > 0 {
							methodValue.Name = k.Names[0].String()
						}
						method.ResultList = append(method.ResultList, methodValue)
					}
				}
				v.par.Handler.Method = append(v.par.Handler.Method, method)
			}
			// fmt.Println(d.Name.Name[:3], d.Name.Name[len(d.Name.Name)-7:len(d.Name.Name)])
		}

	}
}

func (v *repoParser) getFuncType(n ast.Node) *ast.FuncType {
	switch d := n.(type) {
	case *ast.FuncType:
		return d
	default:
		return nil
	}
}

func (v *repoParser) getTypeExpr(n ast.Node, structName string) string {
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
