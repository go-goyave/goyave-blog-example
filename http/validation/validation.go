package validation

// If none of the available validation rules satisfy your needs, you can implement custom validation rules.
// https://system-glitch.github.io/goyave/guide/basics/validation.html#custom-rules

// import "github.com/System-Glitch/goyave/v3/validation"

// func validateCustomFormat(field string, value interface{}, parameters []string, form map[string]interface{}) bool {
// 	str, ok := value.(string)

// 	if ok { // The data under validation is a string
// 		return regexp.MustCompile(parameters[0]).MatchString(str)
// 	}

// 	return false // Cannot validate this field
// }

func init() {
	// Register your custom validation rules here.
	// validation.AddRule("custom_format", &validation.RuleDefinition{
	// 	Function:           validateCustomFormat,
	// 	RequiredParameters: 1, // Ensure the rule has at least one parameter
	// })
}
