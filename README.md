# vulcan: simple and ORM-less microservice skeleton

[![license](https://img.shields.io/github/license/exfly/vulcan.svg)](LICENSE)

[vulcan](https://github.com/exfly/vulcan), a simple, ORM-less microservice skeleton, has been written in Golang. No ORM. Grpc for internal service and [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) for external request.

## Table of Contents

- [Technical stack](#Technical)
- [Usage](#usage)
- [API](#api)
- [Contributing](#contributing)
- [License](#license)

## Technical

- Backend building blocks
	- [grpc-ecosystem/grpc-gateway/v2](https://github.com/grpc-ecosystem/grpc-gateway)
	- [kyleconroy/sqlc](https://github.com/kyleconroy/sqlc) generate type-save code from sql
	- [golang-migrate/migrate/v4](https://github.com/golang-migrate/migrate) for schema migration
	- [jmoiron/sqlx](https://github.com/jmoiron/sqlx) advance `database/sql`
	- [Masterminds/squirrel](https://github.com/Masterminds/squirrel) for dynamic sql builder
	- utils
		- [sirupsen/logrus](https://github.com/sirupsen/logrus) for log
		- [spf13/viper](https://github.com/spf13/viper) for manage config
		- [stretchr/testify](https://github.com/stretchr/testify) for test

### No ORM

ORM can speed up development efficiency to a certain extent, but to a certain extent, it reduces the control of details and abandons performance. So no user ORM.

## Usage

> Assuming you have installed docker, docker-compose and golang

```
docker-compose up -d
go build -o vulcan cmd/vulcan/main.go cmd/vulcan/wire_gen.go
```

## Contributing

See [the contributing file](CONTRIBUTING.md)!

PRs accepted.

## License

[Apache © exfly.](./LICENSE)
