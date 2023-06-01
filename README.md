<p align="center">
  <a href="https://go.dev/" target="blank"><img src="https://static.imasters.com.br/wp-content/uploads/2018/10/24174307/0_OWUKWmE-4jdrLpx7.png" width="200" alt="GO Lang Logo" /></a>
</p>

[circleci-image]: https://img.shields.io/circleci/build/github/nestjs/nest/master?token=abc123def456
[circleci-url]: https://circleci.com/gh/nestjs/nest

<p align="center">A simple, secure, scalable api with <a href="https://go.dev/" target="_blank">Golang</a></p>

<p align="center">
    <a href="https://pkg.go.dev/golang.org/x/example" target="_blank"><img src="https://pkg.go.dev/badge/golang.org/x/example.svg" alt="GoLang" /></a>
</p>

## Description

Devbook API with <a href="https://go.dev/" target="_blank">Golang</a>

## Add missing and remove unused modules

```bash
$ go mod tidy
```

## Running the application

```bash
# development
$ go run main.go

# production mode
$ go build
```

## TODO

- Add migrations in the project
- Add dockerfile for development and production envoirments
- Add hot reload in the project
- Improve Postman collection with envs and scripts to provide variables
- Add unit tests on the project
