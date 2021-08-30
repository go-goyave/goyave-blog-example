package validation

import "goyave.dev/goyave/v4/validation"

func init() {
	validation.AddRule("password", &validation.RuleDefinition{
		Function:           validatePassword,
		RequiredParameters: 0,
	})
}

func isLowerCaseLetter(r rune) bool {
	return r >= 'a' && r <= 'z'
}

func isUpperCaseLetter(r rune) bool {
	return r >= 'A' && r <= 'Z'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isSpecialChar(r rune) bool {
	return !isLowerCaseLetter(r) && !isUpperCaseLetter(r) && !isDigit(r)
}

// validatePassword takes an input and checks if it fulfills password strength criteria:
// - at least one uppercase and one lowercase letter
// - at least one digit
// - at least one special character (! @ # ? ] etc., any utf-8 character that is not a letter or a digit)
func validatePassword(ctx *validation.Context) bool {
	str, ok := ctx.Value.(string)

	if ok {
		lower := false
		upper := false
		digit := false
		special := false

		for _, r := range str {
			switch {
			case isLowerCaseLetter(r):
				lower = true
			case isUpperCaseLetter(r):
				upper = true
			case isDigit(r):
				digit = true
			case isSpecialChar(r):
				special = true
			}
		}
		return lower && upper && special && digit
	}

	return false // Cannot validate this field
}
