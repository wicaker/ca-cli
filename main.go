package main

import (
	"bytes"
	"fmt"

	. "github.com/dave/jennifer/jen"
)

// ca-cli or cc-cli
func main() {
	f := NewFile("a")
	// f.Group.inter
	f.Comment("Foo returns the string \"foo\"")
	f.Group.Func().Id("main").Params(Id("a").Id("string"), Id("b").Id("int64")).Error().Block(
		Return(Id("nil")),
	)
	f.Comment("Test returns the string \"foo\"")
	f.Group.Func().Id("Test").Params(Id("a").Id("string"), Id("b").Id("int64")).Block()
	f.Comment("Ohe return")
	f.Type().Id("Ohe").Struct()
	buf := &bytes.Buffer{}
	err := f.Save("a.go")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(buf.String())
	}
}
