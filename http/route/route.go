package route

import (
	"github.com/System-Glitch/goyave-blog-example/database/model"
	"github.com/System-Glitch/goyave-blog-example/http/controller/hello"
	"github.com/System-Glitch/goyave-blog-example/http/controller/user"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/cors"
	"github.com/System-Glitch/goyave/v3/middleware"
)

// Register all the application routes. This is the main route registrer.
func Register(router *goyave.Router) {

	// Applying default CORS settings (allow all methods and all origins)
	// Learn more about CORS options here: https://system-glitch.github.io/goyave/guide/advanced/cors.html
	router.CORS(cors.Default())
	router.Middleware(middleware.DisallowNonValidatedFields)

	// Register your routes here

	// Route without validation
	router.Get("/hello", hello.SayHi)

	// Route with validation
	router.Post("/echo", hello.Echo).Validate(hello.EchoRequest)

	registerUserRoutes(router)
}

func registerUserRoutes(parent *goyave.Router) {
	userRouter := parent.Subrouter("/user")
	userRouter.Post("/login", auth.NewJWTController(&model.User{}).Login).Validate(user.LoginRequest)
	userRouter.Post("/", user.Register).Validate(user.InsertRequest)
	userRouter.Get("/{id:[0-9+]}/image", user.Image)

	authenticator := auth.Middleware(&model.User{}, &auth.JWTAuthenticator{})
	authRouter := userRouter.Subrouter("")
	authRouter.Middleware(authenticator)
	authRouter.Get("/", user.Show)
	authRouter.Patch("/", user.Update).Validate(user.UpdateRequest)

}
