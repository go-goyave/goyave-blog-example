package hello

import (
	"net/http"

	"github.com/System-Glitch/goyave/v3"
)

// Controllers are files containing a collection of Handlers related to a specific feature.
// Each feature should have its own package.
//
// Controller handlers contain the business logic of your application.
// They should be concise and focused on what matters for this particular feature in your application.
// Learn more about controllers here: https://system-glitch.github.io/goyave/guide/basics/controllers.html

// ----------------------------------------------------------------------

// SayHi is a controller handler writing "Hi!" as a response.
//
// The Response object is used to write your response.
// https://system-glitch.github.io/goyave/guide/basics/responses.html
//
// The Request object contains all the information about the incoming request, including it's parsed body,
// query params and route parameters.
// https://system-glitch.github.io/goyave/guide/basics/requests.html
func SayHi(response *goyave.Response, request *goyave.Request) {
	response.String(http.StatusOK, "Hi!")
}

// Echo is a controller handler writing the input field "text" as a response.
// This route is validated. See "http/request/echorequest/echo.go" for more details.
func Echo(response *goyave.Response, request *goyave.Request) {
	response.String(http.StatusOK, request.String("text"))
}
