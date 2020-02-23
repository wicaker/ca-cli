package generator

import (
	"io/ioutil"
)

func (gen *caGen) GenEnv(dirName string) error {
	configEnv := []byte(`DB_USER_TEST=
DB_PASSWORD_TEST=
DB_NAME_TEST=test

GIN_MODE=release

DATABASE_HOST_PG=127.0.0.1
DATABASE_PORT_PG=5432
DATABASE_USER_PG=
DATABASE_PASSWORD_PG=
DATABASE_NAME_PG=

DATABASE_HOST_MYSQL=127.0.0.1
DATABASE_PORT_MYSQL=3306
DATABASE_USER_MYSQL=
DATABASE_PASSWORD_MYSQL=
DATABASE_NAME_MYSQL=

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
