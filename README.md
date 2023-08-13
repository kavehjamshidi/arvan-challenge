# arvan-challenge

ArvanCloud's challenge for backend developer position.

## Prerequisites

- Go version 1.20 was used as the programming language.
- Redis is required as temporary data store. Redis version 7.0.8 was used.
- PostgreSQL is required as persistent data store. PostgreSQL version 15.4 was used.

## Clone this project

```shell
$ git clone https://github.com/kavehjamshidi/arvan-challenge.git

$ cd arvan-challenge
```

## Installing Go Dependencies

All dependencies can be downloaded using the command below:

```shell
$ go mod download
```

## Setting up Environment Variables

For the sake of simplicity, Environment Variables are used for configurations. If no environment variable is provided,
the fallback values which are defined as constants in the code are used.
The table below includes all required Environment Variables and their respective fallback values:

|                                        | Environment Variable Name | Default Value                                                          |
|----------------------------------------|---------------------------|------------------------------------------------------------------------|
| Environment                            | `ENV`                     | `dev`                                                                  |
| Server Address                         | `SERVER_ADDRESS`          | `:4000`                                                                |
| Redis Address                          | `REDIS_ADDRESS`           | `localhost:6379`                                                       |
| Redis Password                         | `REDIS_PASSWORD`          |                                                                        |
| Test Redis Address (Integration Test)  | `TEST_REDIS_ADDRESS`      | `localhost:6379`                                                       |
| Test Redis Password (Integration Test) | `TEST_REDIS_PASSWORD`     |                                                                        |
| Database URI                           | `DB_URI`                  | `postgres://postgres:very-secret@127.0.0.1:5432/arvan?sslmode=disable` |
| Test Database URI (Integration Test)   | `TEST_DB_URI`             | `postgres://postgres:very-secret@127.0.0.1:5432/arvan?sslmode=disable` |
| Database Migrations Table              | `DB_MIGRATION_TABLE`      | `migrations`                                                           |

## Build and Test

To run all tests, run the command below:

```shell
$ go test ./...
```

To run the project using Docker Compose, use this command:

```shell
$ docker compose up
```

## Endpoints

A Postman collection is included in the /docs directory as the documentation for this project.

