package generator

import (
	"io/ioutil"
)

func (gen *caGen) GenEnv(dirName string) error {
	configEnv := []byte(`DB_USER_TEST=
DB_PASSWORD_TEST=
DB_HOST_TEST=
DB_PORT_TEST=
DB_NAME_TEST=test

GIN_MODE=release

DATABASE_HOST=
DATABASE_PORT=
DATABASE_USER=
DATABASE_PASSWORD=
DATABASE_NAME=

DATABASE_MONGO_URL=
DATABASE_MONGO_NAME=

SERVER_ECHO_PORT=9090
SERVER_GIN_PORT=8090
SERVER_GORILLA_MUX_PORT=7090
SERVER_NET_HTTP_SERVER_MUX_PORT=6090
SERVER_GRAPHQL_SERVER_MUX_PORT=5090
SERVER_GRPC_PORT=50051`)

	err := ioutil.WriteFile("./"+dirName+"/.env", configEnv, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("./"+dirName+"/.env.example", configEnv, 0644)
	if err != nil {
		return err
	}

	return nil
}
