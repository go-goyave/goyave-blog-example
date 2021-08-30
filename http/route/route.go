package route

import (
	"github.com/go-goyave/goyave-blog-example/database/model"
	"github.com/go-goyave/goyave-blog-example/http/controller/article"
	"github.com/go-goyave/goyave-blog-example/http/controller/user"
	"github.com/go-goyave/goyave-blog-example/http/middleware"

	"goyave.dev/goyave/v4"
	"goyave.dev/goyave/v4/auth"
	"goyave.dev/goyave/v4/cors"
	"goyave.dev/goyave/v4/middleware/ratelimiter"
)

// Register all the application routes. This is the main route registrer.
func Register(router *goyave.Router) {

	router.CORS(cors.Default())
	router.Middleware(ratelimiter.New(model.RateLimiterFunc))

	authenticator := auth.Middleware(&model.User{}, &auth.JWTAuthenticator{})

	registerUserRoutes(router, authenticator)
	registerArticleRoutes(router, authenticator)
}

func registerUserRoutes(parent *goyave.Router, authenticator goyave.Middleware) {

	jwtController := auth.NewJWTController(&model.User{})
	jwtController.UsernameField = "email"
	userRouter := parent.Subrouter("/user")
	userRouter.Post("/login", jwtController.Login).Validate(user.LoginRequest)
	userRouter.Post("/", user.Register).Validate(user.InsertRequest)
	userRouter.Get("/{id:[0-9]+}/image", user.Image)

	authRouter := userRouter.Group()
	authRouter.Middleware(authenticator)
	authRouter.Get("/", user.Show)
	authRouter.Patch("/", user.Update).Validate(user.UpdateRequest)

}

func registerArticleRoutes(parent *goyave.Router, authenticator goyave.Middleware) {

	articleRouter := parent.Subrouter("/article")
	articleRouter.Get("/", article.Index).Validate(article.IndexRequest)
	articleRouter.Get("/{slug}", article.Show)

	authRouter := articleRouter.Group()
	authRouter.Middleware(authenticator)
	authRouter.Post("/", article.Store).Validate(article.InsertRequest)

	ownedRouter := authRouter.Group()
	ownerMiddleware := middleware.Owner("id", "author_id", &model.Article{})
	ownedRouter.Middleware(ownerMiddleware)
	ownedRouter.Patch("/{id:[0-9]+}", article.Update).Validate(article.UpdateRequest)
	ownedRouter.Patch("/{slug}", article.Update).Validate(article.UpdateRequest)
	ownedRouter.Delete("/{id:[0-9]+}", article.Destroy)
	ownedRouter.Delete("/{slug}", article.Destroy)

}
