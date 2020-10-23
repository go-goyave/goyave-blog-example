package route

import (
	"github.com/System-Glitch/goyave-blog-example/database/model"
	"github.com/System-Glitch/goyave-blog-example/http/controller/user"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/cors"
	"github.com/System-Glitch/goyave/v3/middleware"
)

// Register all the application routes. This is the main route registrer.
func Register(router *goyave.Router) {

	router.CORS(cors.Default())
	router.Middleware(middleware.DisallowNonValidatedFields)

	registerUserRoutes(router)
}

func registerUserRoutes(parent *goyave.Router) {
	jwtController := auth.NewJWTController(&model.User{})
	jwtController.UsernameField = "email"

	userRouter := parent.Subrouter("/user")
	userRouter.Post("/login", jwtController.Login).Validate(user.LoginRequest)
	userRouter.Post("/", user.Register).Validate(user.InsertRequest)
	userRouter.Get("/{id:[0-9]+}/image", user.Image)

	authenticator := auth.Middleware(&model.User{}, &auth.JWTAuthenticator{})
	authRouter := userRouter.Subrouter("")
	authRouter.Middleware(authenticator)
	authRouter.Get("/", user.Show)
	authRouter.Patch("/", user.Update).Validate(user.UpdateRequest)

}
