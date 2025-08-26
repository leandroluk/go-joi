package joi_test

import (
	"regexp"
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestStringSchema_Base(t *testing.T) {
	schema := joi.String()

	_, errs1 := schema.Validate("field", "ok")
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", 123)
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.String()
	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Min(t *testing.T) {
	schema := joi.String().Min(3)

	_, errs1 := schema.Validate("field", "abcd")
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "ab")
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Min_NonStringWithNil(t *testing.T) {
	schema := joi.String().Min(3)
	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Max(t *testing.T) {
	schema := joi.String().Max(3)

	_, errs1 := schema.Validate("field", "abc")
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "abcd")
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Max_NonStringWithNil(t *testing.T) {
	schema := joi.String().Max(3)
	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Regex(t *testing.T) {
	schema := joi.String().Regex(regexp.MustCompile(`^[a-z]+$`))

	_, errs1 := schema.Validate("field", "abc")
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "123")
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Regex_NonStringWithNil(t *testing.T) {
	re := regexp.MustCompile(`^ok$`)
	schema := joi.String().Regex(re)
	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Trim(t *testing.T) {
	schema := joi.String().Trim()

	val, errs := schema.Validate("field", "  abc  ")
	assert.Empty(t, errs)
	assert.Equal(t, "abc", val)
}

func TestStringSchema_Trim_NonStringWithNil(t *testing.T) {
	schema := joi.String().Trim()
	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Lowercase(t *testing.T) {
	schema := joi.String().Lowercase()

	val, errs1 := schema.Validate("field", "ABC")
	assert.NotEmpty(t, errs1)
	assert.Equal(t, "abc", val)

	val, errs2 := schema.Validate("field", "abc")
	assert.Empty(t, errs2)
	assert.Equal(t, "abc", val)
}

func TestStringSchema_Lowercase_NonStringWithNil(t *testing.T) {
	schema := joi.String().Lowercase()
	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Uppercase(t *testing.T) {
	schema := joi.String().Uppercase()

	val, errs1 := schema.Validate("field", "abc")
	assert.NotEmpty(t, errs1)
	assert.Equal(t, "ABC", val)

	val, errs2 := schema.Validate("field", "ABC")
	assert.Empty(t, errs2)
	assert.Equal(t, "ABC", val)
}

func TestStringSchema_Uppercase_NonStringWithNil(t *testing.T) {
	schema := joi.String().Uppercase()
	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}
