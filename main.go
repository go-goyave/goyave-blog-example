package main

import (
	"os"

	"github.com/go-goyave/goyave-blog-example/database/seeder"
	"github.com/go-goyave/goyave-blog-example/http/route"

	_ "github.com/go-goyave/goyave-blog-example/http/validation"

	"github.com/System-Glitch/goyave/v3"
	_ "github.com/System-Glitch/goyave/v3/database/dialect/mysql"
)

func main() {
	goyave.RegisterStartupHook(seeder.Run)

	if err := goyave.Start(route.Register); err != nil {
		os.Exit(err.(*goyave.Error).ExitCode)
	}
}
