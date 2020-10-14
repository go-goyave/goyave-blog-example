<p align="center">
    <img src="https://raw.githubusercontent.com/System-Glitch/goyave/master/resources/img/logo/goyave_text.png" alt="Goyave Logo" width="550"/>
</p>

## Goyave Blog Example
### 🚧 WORK IN PROGRESS

![https://github.com/System-Glitch/goyave-blog-example/actions](https://github.com/System-Glitch/goyave-blog-example/workflows/Test/badge.svg)

This codebase was created to demonstrate a fully fledged fullstack application built with **[Goyave](https://github.com/System-Glitch/goyave)** including CRUD operations, authentication, routing, pagination, and more.

## Getting Started

### Requirements

- Go 1.13+
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
│   │   ├── placeholder.go   // Placeholders for validation messages of custom rules 
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
├── config.json
├── config.test.json         // Config file used for tests
├── go.mod
└── kernel.go                // Application entrypoint
```

### Running the project

First, make your own configuration for your local environment. You can copy `config.example.json` to `config.json`.

Run `go run kernel.go` in your project's directory to start the server, then try to request the `hello` route.
```
$ curl http://localhost:8080/hello
Hi!
```

There is also an `echo` route, with a basic body validation.
```
$ curl -H "Content-Type: application/json" -X POST -d '{"text":"abc 123"}' http://localhost:8080/echo
abc 123
```

**Using docker:**

```
docker-compose up
```

**Run tests with docker:**

```
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

## Learning Goyave

The Goyave framework has an extensive documentation covering in-depth subjects and teaching you how to run a project using Goyave from setup to deployment.

<a href="https://system-glitch.github.io/goyave/guide/installation"><h3 align="center">Read the documentation</h3></a>

<a href="https://pkg.go.dev/github.com/System-Glitch/goyave/v3"><h3 align="center">pkg.go.dev</h3></a>

## License

This example project is MIT Licensed. Copyright © 2020 Jérémy LAMBERT (SystemGlitch) 

The Goyave framework is MIT Licensed. Copyright © 2019 Jérémy LAMBERT (SystemGlitch)
