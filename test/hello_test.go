package test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/System-Glitch/goyave-blog-example/http/route"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/validation"
)

// Goyave provides an API to ease the unit and functional testing of your application.
// This API is an extension of testify (https://github.com/stretchr/testify).
// "goyave.TestSuite" inherits from testify's "suite.Suite", and sets up the environment for you.
//
// That means:
// - "GOYAVE_ENV" environment variable is set to `test` and restored to its original value when the suite is done.
// - All tests are run using your project's root as working directory.
//   This directory is determined by the presence of a `go.mod` file.
// - Config and language files are loaded before the tests start. As the environment is set to `test`,
//   you need a `config.test.json` in the root directory of your project.
//
// This setup is done by the function `goyave.RunTest`, so you shouldn't run your test suites using testify's `suite.Run()` function.
//
// Learn more about testing here: https://system-glitch.github.io/goyave/guide/advanced/testing.html

type HelloTestSuite struct { // Create a test suite for the Hello controller
	goyave.TestSuite
}

func (suite *HelloTestSuite) TestHello() {
	suite.RunServer(route.Register, func() {
		resp, err := suite.Get("/hello", nil)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(200, resp.StatusCode)
			suite.Equal("Hi!", string(suite.GetBody(resp)))
		}
	})
}

func (suite *HelloTestSuite) TestEcho() {
	suite.RunServer(route.Register, func() {
		resp, err := suite.Post("/echo", nil, nil)
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil { // Expect validation errors (field "text" is required)
			defer resp.Body.Close()
			suite.Equal(422, resp.StatusCode)
			json := map[string]validation.Errors{}
			err := suite.GetJSONBody(resp, &json)
			suite.Nil(err)
			if err == nil {
				textErrors, ok := json["validationError"]["text"]
				suite.True(ok)
				suite.Equal(2, len(textErrors))
			}
		}

		headers := map[string]string{"Content-Type": "application/json"}
		body, _ := json.Marshal(map[string]interface{}{"text": "hello world"})
		resp, err = suite.Post("/echo", headers, bytes.NewReader(body))
		suite.Nil(err)
		suite.NotNil(resp)
		if resp != nil {
			defer resp.Body.Close()
			suite.Equal(200, resp.StatusCode)
			suite.Equal("hello world", string(suite.GetBody(resp)))
		}
	})
}

func TestHelloSuite(t *testing.T) { // Run the test suite
	goyave.RunTest(t, new(HelloTestSuite))
}
