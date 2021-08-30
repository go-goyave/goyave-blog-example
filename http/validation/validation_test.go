package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"goyave.dev/goyave/v4/validation"
)

func TestIsLowerCaseLetter(t *testing.T) {
	assert.True(t, isLowerCaseLetter('a'))
	assert.True(t, isLowerCaseLetter('z'))
	assert.False(t, isLowerCaseLetter('A'))
	assert.False(t, isLowerCaseLetter('Z'))
	assert.False(t, isLowerCaseLetter('0'))
	assert.False(t, isLowerCaseLetter('9'))
	assert.False(t, isLowerCaseLetter(' '))
	assert.False(t, isLowerCaseLetter('*'))
	assert.False(t, isLowerCaseLetter('·'))
	assert.False(t, isLowerCaseLetter('👍'))
}

func TestIsUpperCaseLetter(t *testing.T) {
	assert.False(t, isUpperCaseLetter('a'))
	assert.False(t, isUpperCaseLetter('z'))
	assert.True(t, isUpperCaseLetter('A'))
	assert.True(t, isUpperCaseLetter('Z'))
	assert.False(t, isUpperCaseLetter('0'))
	assert.False(t, isUpperCaseLetter('9'))
	assert.False(t, isUpperCaseLetter(' '))
	assert.False(t, isUpperCaseLetter('*'))
	assert.False(t, isUpperCaseLetter('·'))
	assert.False(t, isUpperCaseLetter('👍'))
}

func TestIsDigit(t *testing.T) {
	assert.False(t, isDigit('a'))
	assert.False(t, isDigit('z'))
	assert.False(t, isDigit('A'))
	assert.False(t, isDigit('Z'))
	assert.True(t, isDigit('0'))
	assert.True(t, isDigit('9'))
	assert.False(t, isDigit(' '))
	assert.False(t, isDigit('*'))
	assert.False(t, isDigit('·'))
	assert.False(t, isDigit('👍'))
}

func TestIsSpecialChar(t *testing.T) {
	assert.False(t, isSpecialChar('a'))
	assert.False(t, isSpecialChar('z'))
	assert.False(t, isSpecialChar('A'))
	assert.False(t, isSpecialChar('Z'))
	assert.False(t, isSpecialChar('0'))
	assert.False(t, isSpecialChar('9'))
	assert.True(t, isSpecialChar(' '))
	assert.True(t, isSpecialChar('*'))
	assert.True(t, isSpecialChar('·'))
	assert.True(t, isSpecialChar('👍'))
}

func TestValidatePassword(t *testing.T) {
	assert.True(t, validatePassword(&validation.Context{Value: "pAssword.1"}))
	assert.False(t, validatePassword(&validation.Context{Value: "pAssword."}))
	assert.False(t, validatePassword(&validation.Context{Value: "pAssword1"}))
	assert.False(t, validatePassword(&validation.Context{Value: "password.1"}))
	assert.False(t, validatePassword(&validation.Context{Value: "PASSWORD.1"}))
	assert.False(t, validatePassword(&validation.Context{Value: 42}))
}
