<p align="center">
    <img src="./.github/img/goyave_banner.png#gh-light-mode-only" alt="Goyave Logo" width="550"/>
    <img src="./.github/img/goyave_banner_dark.png#gh-dark-mode-only" alt="Goyave Logo" width="550"/>
</p>

## Goyave Blog Example

![https://github.com/go-goyave/goyave-blog-example/actions](https://github.com/go-goyave/goyave-blog-example/workflows/Test/badge.svg)

This example project was created to demonstrate a simple application built with **[Goyave](https://github.com/go-goyave/goyave)** including CRUD operations, authentication, routing, pagination, and more. With this application, users can register, login and write blog posts (articles) or read the other user's ones.

## Running the project

First, make your own configuration for your local environment.

- Copy `config.example.json` to `config.json`.
- Start the database container with `docker-compose up`.
- Run migrations with [dbmate](https://github.com/amacneil/dbmate): `dbmate -u postgres://dbuser:secret@127.0.0.1:5432/blog?sslmode=disable -d ./database/migrations --no-dump-schema migrate`
- Run `go run main.go` in your project's directory to start the server. If you want to seed your database with random records use the `-seed` flag: `go run main.go -seed`. Users will all be created with the following password: `p4ssW0rd_`

## Resources

- [Documentation](https://goyave.dev)
- [go.pkg.dev](https://pkg.go.dev/goyave.dev/goyave/v5)