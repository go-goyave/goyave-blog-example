package user

import "goyave.dev/goyave/v4/validation"

var (
	// InsertRequest validates Post requests for users
	InsertRequest validation.RuleSet = validation.RuleSet{
		"email":    validation.List{"required", "string", "email", "between:3,100", "unique:users"},
		"username": validation.List{"required", "string", "between:3,100", "unique:users"},
		"image":    validation.List{"nullable", "file", "image", "max:2048", "count:1"},
		"password": validation.List{"required", "string", "between:6,100", "password"},
	}

	// UpdateRequest validates Patch requests for users
	UpdateRequest validation.RuleSet = validation.RuleSet{
		"email":    validation.List{"string", "email", "between:3,100", "unique:users"},
		"username": validation.List{"string", "between:3,100", "unique:users"},
		"image":    validation.List{"nullable", "file", "image", "max:2048", "count:1"},
		"password": validation.List{"string", "between:6,100", "password"},
	}

	// LoginRequest validates user login requests
	LoginRequest validation.RuleSet = validation.RuleSet{
		"email":    validation.List{"required", "string", "email"},
		"password": validation.List{"required", "string"},
	}
)
