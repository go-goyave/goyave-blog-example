package user

import "github.com/System-Glitch/goyave/v3/validation"

var (
	// InsertRequest validates Post requests for users
	InsertRequest validation.RuleSet = validation.RuleSet{
		"email":    {"required", "string", "email", "between:3,100", "unique:users"},
		"username": {"required", "string", "between:3,100", "unique:users"},
		"image":    {"nullable", "file", "image", "max:2048", "count:1"},
		"password": {"required", "string", "between:6,100"}, // TODO implement password validation
	}

	// UpdateRequest validates Put requests for users
	UpdateRequest validation.RuleSet = validation.RuleSet{
		"email":    {"nullable", "string", "email", "between:3,100", "unique:users"},
		"username": {"nullable", "string", "between:3,100", "unique:users"},
		"image":    {"nullable", "file", "image", "max:2048", "count:1"},
		"password": {"nullable", "string", "between:6,100"},
	}

	// LoginRequest validates user login requests
	LoginRequest validation.RuleSet = validation.RuleSet{
		"email":    {"required", "string", "email"},
		"password": {"required", "string"},
	}
)
