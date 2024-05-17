package route

import (
	"github.com/go-goyave/goyave-blog-example/http/controller/article"
	"github.com/go-goyave/goyave-blog-example/http/controller/user"
	"github.com/go-goyave/goyave-blog-example/service"
	userservice "github.com/go-goyave/goyave-blog-example/service/user"
	"goyave.dev/goyave/v5"
	"goyave.dev/goyave/v5/auth"
	"goyave.dev/goyave/v5/cors"
	"goyave.dev/goyave/v5/log"
	"goyave.dev/goyave/v5/middleware/parse"
)

func Register(server *goyave.Server, router *goyave.Router) {
	router.CORS(cors.Default())
	router.GlobalMiddleware(log.CombinedLogMiddleware())

	userService := server.Service(service.User).(*userservice.Service)
	authenticator := auth.NewJWTAuthenticator(userService)
	authMiddleware := auth.Middleware(authenticator)
	router.GlobalMiddleware(authMiddleware)

	router.GlobalMiddleware(&parse.Middleware{})

	loginController := auth.NewJWTController(userService, "Password")
	loginController.UsernameRequestField = "email"
	router.Controller(loginController)
	router.Controller(user.NewController())
	router.Controller(article.NewController())
}
