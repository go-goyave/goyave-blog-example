package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/go-goyave/goyave-blog-example/database/repository"
	seeders "github.com/go-goyave/goyave-blog-example/database/seed"
	"github.com/go-goyave/goyave-blog-example/http/route"
	"github.com/go-goyave/goyave-blog-example/service/article"
	"github.com/go-goyave/goyave-blog-example/service/storage"
	"github.com/go-goyave/goyave-blog-example/service/user"

	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/util/errors"
	"goyave.dev/goyave/v5/util/fsutil"
	"goyave.dev/goyave/v5/util/fsutil/osfs"
	"goyave.dev/goyave/v5/util/session"

	_ "goyave.dev/goyave/v5/database/dialect/postgres"
)

//go:embed resources
var resources embed.FS

func main() {
	var seed bool
	flag.BoolVar(&seed, "seed", false, "If true, the database will be seeded with random data.")
	flag.Parse()

	resources := fsutil.NewEmbed(resources)
	langFS, err := resources.Sub("resources/lang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.(*errors.Error).String())
		os.Exit(1)
	}

	opts := goyave.Options{
		LangFS: langFS,
	}

	server, err := goyave.New(opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.(*errors.Error).String())
		os.Exit(1)
	}

	server.Logger.Info("Registering hooks")
	server.RegisterSignalHook()

	server.RegisterStartupHook(func(s *goyave.Server) {
		server.Logger.Info("Server is listening", "host", s.Host())
	})

	server.RegisterShutdownHook(func(s *goyave.Server) {
		s.Logger.Info("Server is shutting down")
	})

	registerServices(server)

	server.Logger.Info("Registering routes")
	server.RegisterRoutes(route.Register)

	if seed {
		server.Logger.Info("Seeding database...")
		seeders.Seed(server.DB())
	}

	if err := server.Start(); err != nil {
		server.Logger.Error(err)
		os.Exit(2)
	}
}

func registerServices(server *goyave.Server) {
	server.Logger.Info("Registering services")

	session := session.GORM(server.DB(), nil)

	userRepo := repository.NewUser(server.DB())
	articleRepo := repository.NewArticle(server.DB())

	storageFS, err := (&osfs.FS{}).Sub(".storage/avatars")
	if err != nil {
		panic(errors.New(err))
	}
	storageService := storage.NewService(storageFS)

	server.RegisterService(storageService)
	server.RegisterService(user.NewService(session, server.Logger, userRepo, storageService))
	server.RegisterService(article.NewService(session, articleRepo))
}
