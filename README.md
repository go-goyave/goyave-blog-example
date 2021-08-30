<p align="center">
    <img src="https://raw.githubusercontent.com/System-Glitch/goyave/master/resources/img/logo/goyave_text.png" alt="Goyave Logo" width="550"/>
</p>

## Goyave Blog Example

![https://github.com/go-goyave/goyave-blog-example/actions](https://github.com/go-goyave/goyave-blog-example/workflows/Test/badge.svg)

This codebase was created to demonstrate a fully fledged fullstack application built with **[Goyave](https://github.com/System-Glitch/goyave)** including CRUD operations, authentication, routing, pagination, and more.

## Getting Started

### Requirements

- Go 1.16+
- Go modules

### Directory structure

```
.
├── database
│   ├── model                // ORM models
│   |   └── ...
│   └── seeder               // Generators for database testing
│       └── ...
├── http
│   ├── controller           // Business logic of the application
│   │   └── ...
│   ├── middleware           // Logic executed before or after controllers
│   │   └── ...
│   ├── validation
│   │   └── validation.go    // Custom validation rules
│   └── route
│       └── route.go         // Routes definition
│
├── resources
│   └── lang
│       └── en-US            // Overrides to the default language lines
│           ├── fields.json
│           ├── locale.json
│           └── rules.json
│
├── test                     // Functional tests
|   └── ...
|
├── .gitignore
├── .golangci.yml            // Settings for the Golangci-lint linter
├── config.example.json      // Example config for local development
├── config.test.json         // Config file used for tests
├── go.mod
└── main.go                  // Application entrypoint
```

### Running the project

First, make your own configuration for your local environment. You can copy `config.example.json` to `config.json`.

Run `go run main.go` in your project's directory to start the server.

**Using docker:**

```
docker-compose up
```

**Run tests with docker:**

```
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

**Database seeding:**

If `app.environment` is set to `localhost` in the config and if the database is empty (no record in the users table), the seeders will be executed and a random dataset will be generated and inserted into the database.

## Learning Goyave

The Goyave framework has an extensive documentation covering in-depth subjects and teaching you how to run a project using Goyave from setup to deployment.

<a href="https://goyave.dev/guide/installation"><h3 align="center">Read the documentation</h3></a>

<a href="https://pkg.go.dev/goyave.dev/goyave/v4"><h3 align="center">pkg.go.dev</h3></a>

## License

This example project is MIT Licensed. Copyright © 2020 Jérémy LAMBERT (SystemGlitch) 

The Goyave framework is MIT Licensed. Copyright © 2019 Jérémy LAMBERT (SystemGlitch)
