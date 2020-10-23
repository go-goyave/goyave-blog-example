package article

import "github.com/System-Glitch/goyave/v3/validation"

var (
	InsertRequest validation.RuleSet = validation.RuleSet{
		"title":    {"required", "string", "max:200"},
		"contents": {"required", "string"},
	}

	UpdateRequest validation.RuleSet = validation.RuleSet{
		"title":    {"string", "max:200"},
		"contents": {"string"},
	}
)
