package article

import "goyave.dev/goyave/v3/validation"

var (
	// InsertRequest validates Post requests for articles
	InsertRequest validation.RuleSet = validation.RuleSet{
		"title":    {"required", "string", "max:200"},
		"contents": {"required", "string"},
	}

	// UpdateRequest validates Patch requests for articles
	UpdateRequest validation.RuleSet = validation.RuleSet{
		"title":    {"string", "max:200"},
		"contents": {"string"},
	}

	// IndexRequest validates query parameters for paginating articles
	IndexRequest validation.RuleSet = validation.RuleSet{
		"page":     {"integer", "min:1"},
		"pageSize": {"integer", "between:10,100"},
		"search":   {"string", "max:200"},
	}
)
