package user

import "github.com/System-Glitch/goyave/v3/validation"

var (
	// InsertRequest validates Post requests for users
	InsertRequest validation.RuleSet = validation.RuleSet{
		"email":    {"required", "string", "email", "between:3,100", "unique:users"},
		"username": {"required", "string", "between:3,100", "unique:users"},
		"image":    {"nullable", "file", "image", "max:2048", "count:1"},
		"password": {"required", "string", "between:6,100", "password"},
	}

	// UpdateRequest validates Patch requests for users
	UpdateRequest validation.RuleSet = validation.RuleSet{
		"email":    {"string", "email", "between:3,100", "unique:users"},
		"username": {"string", "between:3,100", "unique:users"},
		"image":    {"nullable", "file", "image", "max:2048", "count:1"},
		"password": {"string", "between:6,100", "password"},
	}

	// LoginRequest validates user login requests
	LoginRequest validation.RuleSet = validation.RuleSet{
		"email":    {"required", "string", "email"},
		"password": {"required", "string"},
	}
)
