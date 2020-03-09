package parser

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"log"
	"strings"

	"github.com/wicaker/cacli/domain"
)

type doParser struct {
	domainFileName string
	par            domain.Parser
}

// NewParserDomain will create new a doParser object representation of domain.ParserDomain interface
// used particularly to parsing domain layer
func NewParserDomain(dfn string) domain.ParserDomain {
	return &doParser{
		domainFileName: dfn,
	}
}

// DomainParser function
// Call this method will starting parsing process
func (v *doParser) DomainParser(filePath string) (*domain.Parser, error) {
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
	if v.par.Repository.Name == "" || v.par.Usecase.Name == "" {
		return nil, errors.New("interface must exported type and interface name must same with name of file in domain dir")
	}

	return &v.par, nil
}

// Visit method
// implement the Visitor interface from go/ast package to get syntax trees
func (v *doParser) Visit(n ast.Node) ast.Visitor {
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

// to get declration with token is type
func (v *doParser) getType(n ast.Node) {
	switch d := n.(type) {
	case *ast.GenDecl:
		if d.Tok.String() == "type" {
			for _, j := range d.Specs {
				v.getTypeSpec(j)
			}
		}
	}
}

func (v *doParser) getTypeSpec(n ast.Node) {
	switch d := n.(type) {
	case *ast.TypeSpec:
		if _, ok := d.Type.(*ast.InterfaceType); ok {
			// get Usecase layer of interface usecase
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
					field := getMethodsInInterface(d.Type)
					v.initiateUsecase(field, d.Name.Name)
					return
				}
			}

			// get Repository layer of interface repository
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
					field := getMethodsInInterface(d.Type)
					v.initiateRepository(field, d.Name.Name)
					return
				}
			}
		}
	}
}

// initiateRepository will initiate detail information which needed in repository layer
func (v *doParser) initiateRepository(field []*ast.Field, nameType string) {
	for _, j := range field {
		method := domain.Method{
			Name:          j.Names[0].String(),
			ParameterList: []domain.MethodValue{},
			ResultList:    []domain.MethodValue{},
		}
		fType := getFuncType(j.Type)
		for _, k := range fType.Params.List {
			r := getTypeExpr(k.Type, nameType[:len(nameType)-10])
			methodValue := domain.MethodValue{Type: r}
			if len(k.Names) > 0 {
				methodValue.Name = k.Names[0].String()
			}
			method.ParameterList = append(method.ParameterList, methodValue)
		}
		for _, k := range fType.Results.List {
			r := getTypeExpr(k.Type, nameType[:len(nameType)-10])
			methodValue := domain.MethodValue{Type: r}
			if len(k.Names) > 0 {
				methodValue.Name = k.Names[0].String()
			}
			method.ResultList = append(method.ResultList, methodValue)
		}
		v.par.Repository.Method = append(v.par.Repository.Method, method)
	}
}

// initiateUsecase will initiate detail information which needed in usecase layer
func (v *doParser) initiateUsecase(field []*ast.Field, nameType string) {
	for _, j := range field {
		method := domain.Method{
			Name:          j.Names[0].String(),
			ParameterList: []domain.MethodValue{},
			ResultList:    []domain.MethodValue{},
		}
		fType := getFuncType(j.Type)
		for _, k := range fType.Params.List {
			r := getTypeExpr(k.Type, nameType[:len(nameType)-7])
			methodValue := domain.MethodValue{Type: r}
			if len(k.Names) > 0 {
				methodValue.Name = k.Names[0].String()
			}
			method.ParameterList = append(method.ParameterList, methodValue)
		}
		for _, k := range fType.Results.List {
			r := getTypeExpr(k.Type, nameType[:len(nameType)-7])
			methodValue := domain.MethodValue{Type: r}
			if len(k.Names) > 0 {
				methodValue.Name = k.Names[0].String()
			}
			method.ResultList = append(method.ResultList, methodValue)
		}
		v.par.Usecase.Method = append(v.par.Usecase.Method, method)
	}
}
