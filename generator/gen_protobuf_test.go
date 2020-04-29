package generator_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/wicaker/cacli/domain"
	"github.com/wicaker/cacli/fs"
	"github.com/wicaker/cacli/generator"
)

const (
	expected_example_protobuf = `syntax = "proto3";

package proto;

import "google/protobuf/timestamp.proto";

message Example {
	uint64 id = 1;
	google.protobuf.Timestamp createdAt = 2;
	google.protobuf.Timestamp updatedAt = 3;
}

message FetchExampleReq{}
message FetchExampleResp{}

message GetByIDExampleReq{}
message GetByIDExampleResp{}

message StoreExampleReq{}
message StoreExampleResp{}

message UpdateExampleReq{}
message UpdateExampleResp{}

message DeleteExampleReq{}
message DeleteExampleResp{}

service ExampleService {
	rpc FetchExample(FetchExampleReq) returns (FetchExampleResp);
	rpc GetByIDExample(GetByIDExampleReq) returns (GetByIDExampleResp);
	rpc StoreExample(StoreExampleReq) returns (StoreExampleResp);
	rpc UpdateExample(UpdateExampleReq) returns (UpdateExampleResp);
	rpc DeleteExample(DeleteExampleReq) returns (DeleteExampleResp);
}
`
)

func TestGenerateProtobuf(t *testing.T) {
	var (
		serviceName = "test_example_protobuf"
		dirLayer1   = "protobuf"
		domainFile  = "example.go"
		gomodName   = "github.com/example/exampleprotobuf"
		dirName     = fmt.Sprintf("%s/%s", serviceName, dirLayer1)
		newFs       = fs.NewFsService()
	)

	t.Run("success, should generate an example.proto file", func(t *testing.T) {
		// create directory of service
		err := newFs.CreateDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		err = newFs.CreateDir(serviceName + "/" + dirLayer1)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		// generate example.proto file
		gen := generator.NewGeneratorService()
		err = gen.GenProtobuf(dirName, domainFile, gomodName, domain.MockParser)
		resGopg, err := newFs.FindFile(dirName + "/example.proto")
		assert.NoError(t, err)
		assert.NotEqual(t, nil, resGopg)

		data, err := ioutil.ReadFile(dirName + "/example.proto")
		if err != nil {
			log.Error("File reading error", err)
			// remove directory of service
			err = newFs.RemoveDir(serviceName)
			if err != nil {
				log.Error(err)
				os.Exit(1)
			}
			os.Exit(1)
		}
		assert.Equal(t, expected_example_protobuf, string(data))

		// remove directory of service
		err = newFs.RemoveDir(serviceName)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	})

	t.Run("failed, because directory not found", func(t *testing.T) {
		// generate example.proto file
		gen := generator.NewGeneratorService()
		err := gen.GenProtobuf(dirName, domainFile, gomodName, domain.MockParser)

		assert.Error(t, err)
	})
}
