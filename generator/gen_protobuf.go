package generator

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/wicaker/cacli/domain"
)

func (gen *caGen) GenProtobuf(dirName string, domainFile string, gomodName string, parser *domain.Parser) error {
	var (
		file            = path.Base(domainFile)
		domainName      = strings.TrimSuffix(file, filepath.Ext(file))
		domainNameInCap = strings.ToUpper(string(domainName[0])) + domainName[1:]
		protoUsecase    = ``
		protoService    = ``
	)
	for _, i := range parser.Usecase.Method {
		protoUsecase += `message ` + i.Name + domainNameInCap + `Req{}`
		protoUsecase += `
`
		protoUsecase += `message ` + i.Name + domainNameInCap + `Resp{}`
		protoUsecase += `

`
		protoService += `
`
		protoService += `	rpc ` + i.Name + domainNameInCap + `(` + i.Name + domainNameInCap + `Req) returns (` + i.Name + domainNameInCap + `Resp);`
	}

	genProtobuf := []byte(`syntax = "proto3";

package proto;

import "google/protobuf/timestamp.proto";

message ` + domainNameInCap + ` {
	uint64 id = 1;
	google.protobuf.Timestamp createdAt = 2;
	google.protobuf.Timestamp updatedAt = 3;
}

` + protoUsecase + `service ExampleService {` + protoService + `
}
`)

	err := ioutil.WriteFile("./"+dirName+"/"+domainName+".proto", genProtobuf, 0644)
	if err != nil {
		return err
	}
	return nil
}
