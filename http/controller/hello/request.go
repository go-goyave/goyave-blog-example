package hello

import "github.com/System-Glitch/goyave/v3/validation"

// Goyave provides a powerful, yet easy way to validate all incoming data, no matter
// its type or its format, thanks to a large number of validation rules.

// Incoming requests are validated using rules set, which associate rules
// with each expected field in the request.

// Learn more about validation here: https://system-glitch.github.io/goyave/guide/basics/validation.html

// This is the validation rules for the "/echo" route, which is simply
// writing the input as a response.
var (
	EchoRequest validation.RuleSet = validation.RuleSet{
		"text": {"required", "string", "between:3,50"},
	}
)
