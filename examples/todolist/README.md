GO Todolist

## Database
- `soda g config` to create `database.yml` configuration
- `soda create -e development` (to create database development) [more](https://gobuffalo.io/en/docs/db/toolbox/)
- `soda generate -p ./database/migrations fizz name_of_migration` to create new migrations.
- `soda drop -e development` (to drop or delete database) [more](https://gobuffalo.io/en/docs/db/toolbox/)
- migration up : `soda migrate -p database up` [more](https://gobuffalo.io/en/docs/db/migrations/)
- migration down : `soda migrate -p database down -s {number of database want to down}`. For example: `soda migrate -p database down -s 9`

## protobuf
- protoc --go_out=plugins=grpc:. proto/*.proto