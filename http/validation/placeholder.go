package validation

// Validation messages can use placeholders to inject dynamic values in the validation error message. For example, in the rules.json language file:
//
//     "between.string": "The :field must be between :min and :max characters."
//
// Here, the :field placeholder will be replaced by the field name, :min with the first parameter and :max with the second parameter, effectively giving the following result:
//
//     The password must be between 6 and 32 characters.
//
// Placeholders are replacer functions. These functions should return the value to replace the placeholder with.
// Learn more about placeholders here: https://system-glitch.github.io/goyave/guide/basics/validation.html#placeholders

// import "github.com/System-Glitch/goyave/v3/validation"

// func simpleParameterPlaceholder(field string, rule string, parameters []string, language string) string {
// 	return parameters[0]
// }

func init() {
	// Register your custom placeholders here.
	// validation.SetPlaceholder(simpleParameterPlaceholder)
}
