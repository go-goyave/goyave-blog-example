package route

import (
	"goyave-blog-example/http/controller/hello"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/cors"
)

// Routing is an essential part of any Goyave application.
// Routes definition is the action of associating a URI, sometimes having
// parameters, with a handler which will process the request and respond to it.

// Routes are defined in routes registrer functions.
// The main route registrer is passed to "goyave.Start()" and is executed
// automatically with a newly created root-level router.

// Register all the application routes. This is the main route registrer.
func Register(router *goyave.Router) {

	// Applying default CORS settings (allow all methods and all origins)
	// Learn more about CORS options here: https://system-glitch.github.io/goyave/guide/advanced/cors.html
	router.CORS(cors.Default())

	// Register your routes here

	// Route without validation
	router.Get("/hello", hello.SayHi)

	// Route with validation
	router.Post("/echo", hello.Echo).Validate(hello.EchoRequest)
}
