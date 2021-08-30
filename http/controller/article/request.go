package article

import "goyave.dev/goyave/v4/validation"

var (
	// InsertRequest validates Post requests for articles
	InsertRequest validation.RuleSet = validation.RuleSet{
		"title":    validation.List{"required", "string", "max:200"},
		"contents": validation.List{"required", "string"},
	}

	// UpdateRequest validates Patch requests for articles
	UpdateRequest validation.RuleSet = validation.RuleSet{
		"title":    validation.List{"string", "max:200"},
		"contents": validation.List{"string"},
	}

	// IndexRequest validates query parameters for paginating articles
	IndexRequest validation.RuleSet = validation.RuleSet{
		"page":     validation.List{"integer", "min:1"},
		"pageSize": validation.List{"integer", "between:10,100"},
		"search":   validation.List{"string", "max:200"},
	}
)
