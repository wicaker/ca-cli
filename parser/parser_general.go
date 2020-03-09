package parser

import (
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"

	"github.com/wicaker/cacli/domain"
)

type generalParser struct {
	par domain.Parser
}

// NewParserGeneral will create new a generalParser object representation of domain.ParserGeneral interface
// Unlike parser_domain, parser_general can parsing usecase layer, repository layer, transport layer
func NewParserGeneral() domain.ParserGeneral {
	return &generalParser{}
}

// GeneralParser function
// Call this method will starting parsing process
func (v *generalParser) GeneralParser(filePath string) (*domain.Parser, error) {
	var (
		s  scanner.Scanner
		fs = token.NewFileSet()
	)

	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	f := fs.AddFile(filePath, fs.Base(), len(b))
	s.Init(f, b, nil, scanner.ScanComments)

	p, err := parser.ParseFile(fs, filePath, nil, parser.AllErrors)
	if err != nil {
		log.Printf("could not parse %s: %v", filePath, err)
		return nil, err
	}

	ast.Walk(v, p)
	return &v.par, nil
}

// Visit method
// implement the Visitor interface from go/ast package to get syntax trees
func (v *generalParser) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch d := n.(type) {
	case *ast.File:
		for _, j := range d.Decls {
			v.getFuncDecl(j)
		}
	}
	return v
}

func (v *generalParser) getFuncDecl(n ast.Node) {
	switch d := n.(type) {
	case *ast.FuncDecl:
		// spew.Dump(d)
		if len(d.Name.String()) > 7 {
			if d.Name.Name[:3] == "New" && d.Name.Name[len(d.Name.Name)-7:len(d.Name.Name)] == "Usecase" {
				v.initiateUsecase(d)
			}
		}

		if len(d.Name.String()) > 10 {
			if d.Name.Name[:3] == "New" && d.Name.Name[len(d.Name.Name)-10:len(d.Name.Name)] == "Repository" {
				v.initiateRepository(d)
			}
		}

		if len(d.Name.String()) > 7 {
			if d.Name.Name[:3] == "New" && d.Name.Name[len(d.Name.Name)-7:len(d.Name.Name)] == "Handler" {
				v.initiateHandler(d)
			}
		}

	}
}

// initiateRepository will initiate detail information
func (v *generalParser) initiateRepository(d *ast.FuncDecl) {
	method := domain.Method{
		Name:          d.Name.String(),
		ParameterList: []domain.MethodValue{},
		ResultList:    []domain.MethodValue{},
	}
	fType := getFuncType(d.Type)
	for _, k := range fType.Params.List {
		r := getTypeExpr(k.Type, "")

		methodValue := domain.MethodValue{Type: r}
		if len(k.Names) > 0 {
			methodValue.Name = k.Names[0].String()
		}
		method.ParameterList = append(method.ParameterList, methodValue)
	}
	for _, k := range fType.Results.List {
		r := getTypeExpr(k.Type, "")

		methodValue := domain.MethodValue{Type: r}
		if len(k.Names) > 0 {
			methodValue.Name = k.Names[0].String()
		}
		method.ResultList = append(method.ResultList, methodValue)
	}
	v.par.Repository.Method = append(v.par.Repository.Method, method)
}

// initiateUsecase will initiate detail information
func (v *generalParser) initiateUsecase(d *ast.FuncDecl) {
	method := domain.Method{
		Name:          d.Name.String(),
		ParameterList: []domain.MethodValue{},
		ResultList:    []domain.MethodValue{},
	}
	fType := getFuncType(d.Type)
	for _, k := range fType.Params.List {
		r := getTypeExpr(k.Type, "")

		methodValue := domain.MethodValue{Type: r}
		if len(k.Names) > 0 {
			methodValue.Name = k.Names[0].String()
		}
		method.ParameterList = append(method.ParameterList, methodValue)
	}
	for _, k := range fType.Results.List {
		r := getTypeExpr(k.Type, "")

		methodValue := domain.MethodValue{Type: r}
		if len(k.Names) > 0 {
			methodValue.Name = k.Names[0].String()
		}
		method.ResultList = append(method.ResultList, methodValue)
	}
	v.par.Usecase.Method = append(v.par.Usecase.Method, method)
}

// initiateHandler will initiate detail information
func (v *generalParser) initiateHandler(d *ast.FuncDecl) {
	method := domain.Method{
		Name:          d.Name.String(),
		ParameterList: []domain.MethodValue{},
		ResultList:    []domain.MethodValue{},
	}
	fType := getFuncType(d.Type)
	for _, k := range fType.Params.List {
		r := getTypeExpr(k.Type, "")

		methodValue := domain.MethodValue{Type: r}
		if len(k.Names) > 0 {
			methodValue.Name = k.Names[0].String()
		}
		method.ParameterList = append(method.ParameterList, methodValue)
	}
	if fType.Results != nil {
		for _, k := range fType.Results.List {
			r := getTypeExpr(k.Type, "")

			methodValue := domain.MethodValue{Type: r}
			if len(k.Names) > 0 {
				methodValue.Name = k.Names[0].String()
			}
			method.ResultList = append(method.ResultList, methodValue)
		}
	}
	v.par.Handler.Method = append(v.par.Handler.Method, method)
}
