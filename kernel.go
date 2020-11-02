package main

import (
	"os"

	"github.com/System-Glitch/goyave-blog-example/http/route"

	"github.com/System-Glitch/goyave-blog-example/database/dbutil"
	_ "github.com/System-Glitch/goyave-blog-example/http/validation"

	"github.com/System-Glitch/goyave/v3"
	_ "github.com/System-Glitch/goyave/v3/database/dialect/mysql"
)

func main() {
	goyave.RegisterStartupHook(dbutil.RunSeeders)

	if err := goyave.Start(route.Register); err != nil {
		os.Exit(err.(*goyave.Error).ExitCode)
	}
}
