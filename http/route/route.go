package route

import (
	"github.com/System-Glitch/goyave-blog-example/database/model"
	"github.com/System-Glitch/goyave-blog-example/http/controller/article"
	"github.com/System-Glitch/goyave-blog-example/http/controller/user"
	"github.com/System-Glitch/goyave-blog-example/http/middleware"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/cors"
	gmiddleware "github.com/System-Glitch/goyave/v3/middleware"
	"github.com/System-Glitch/goyave/v3/middleware/ratelimiter"
)

// Register all the application routes. This is the main route registrer.
func Register(router *goyave.Router) {

	router.CORS(cors.Default())
	router.Middleware(ratelimiter.New(model.RateLimiterFunc))
	router.Middleware(gmiddleware.DisallowNonValidatedFields)

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

	authRouter := userRouter.Subrouter("")
	authRouter.Middleware(authenticator)
	authRouter.Get("/", user.Show)
	authRouter.Patch("/", user.Update).Validate(user.UpdateRequest)

}

func registerArticleRoutes(parent *goyave.Router, authenticator goyave.Middleware) {

	articleRouter := parent.Subrouter("/article")
	articleRouter.Get("/", article.Index).Validate(article.IndexRequest)
	articleRouter.Get("/{slug}", article.Show)

	authRouter := articleRouter.Subrouter("")
	authRouter.Middleware(authenticator)
	authRouter.Post("/", article.Store).Validate(article.InsertRequest)

	ownedRouter := authRouter.Subrouter("")
	ownerMiddleware := middleware.Owner("id", "author_id", &model.Article{})
	ownedRouter.Middleware(ownerMiddleware)
	ownedRouter.Patch("/{id:[0-9]+}", article.Update).Validate(article.UpdateRequest)
	ownedRouter.Patch("/{slug}", article.Update).Validate(article.UpdateRequest)
	ownedRouter.Delete("/{id:[0-9]+}", article.Destroy)
	ownedRouter.Delete("/{slug}", article.Destroy)

}
