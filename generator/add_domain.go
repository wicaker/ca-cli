package generator

import (
	"bytes"
	"fmt"

	. "github.com/dave/jennifer/jen"
)

func NewGenerateDomain() {
	f := NewFile("a")
	// f.Group.inter
	f.Group.Func().Id("main").Params(Id("a").Id("string"), Id("b").Id("int64")).Error().Block(
		Return(Id("nil")),
	)
	f.Group.Func().Id("test").Params(Id("a").Id("string"), Id("b").Id("int64")).Block()
	buf := &bytes.Buffer{}
	err := f.Save("make/a.go")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(buf.String())
	}
}
