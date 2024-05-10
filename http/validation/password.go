package validation

import "goyave.dev/goyave/v5/validation"

type PasswordValidator struct {
	validation.BaseValidator
}

func (v *PasswordValidator) Validate(ctx *validation.Context) bool {
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

func (v *PasswordValidator) Name() string { return "password" }

// Password the field under validation must contain at least one lower case,
// one upper case letter, one digit and one special character.
func Password() *PasswordValidator {
	return &PasswordValidator{}
}
