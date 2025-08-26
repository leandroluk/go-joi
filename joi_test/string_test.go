package joi_test

import (
	"regexp"
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestStringSchema_Base(t *testing.T) {
	schema := joi.String()

	_, errs1 := schema.Validate("ok")
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate(123)
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.String()
	val, errs := schema.Validate(nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Min(t *testing.T) {
	schema := joi.String().Min(3)

	_, errs1 := schema.ValidateWithOpts("abcd", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("ab", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Min_NonStringWithNil(t *testing.T) {
	schema := joi.String().Min(3)
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Max(t *testing.T) {
	schema := joi.String().Max(3)

	_, errs1 := schema.ValidateWithOpts("abc", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("abcd", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Max_NonStringWithNil(t *testing.T) {
	schema := joi.String().Max(3)
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Length(t *testing.T) {
	schema := joi.String().Length(3)

	_, errs1 := schema.ValidateWithOpts("abc", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("abcd", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Length_NonStringWithNil(t *testing.T) {
	schema := joi.String().Length(3)
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Regex(t *testing.T) {
	schema := joi.String().Regex(regexp.MustCompile(`^[a-z]+$`))

	_, errs1 := schema.ValidateWithOpts("abc", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("123", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestStringSchema_Regex_NonStringWithNil(t *testing.T) {
	re := regexp.MustCompile(`^ok$`)
	schema := joi.String().Regex(re)
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Email_Valid(t *testing.T) {
	schema := joi.String().Email()

	_, errs := schema.ValidateWithOpts("user@example.com", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
}

func TestStringSchema_Email_Invalid(t *testing.T) {
	schema := joi.String().Email()

	_, errs := schema.ValidateWithOpts("not-an-email", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
}

func TestStringSchema_Email_NonStringWithNil(t *testing.T) {
	schema := joi.String().Email()

	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Email_EmptyString(t *testing.T) {
	schema := joi.String().Email()

	val, errs := schema.ValidateWithOpts("", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Equal(t, "", val)
}

func TestStringSchema_Trim(t *testing.T) {
	schema := joi.String().Trim()

	val, errs := schema.ValidateWithOpts("  abc  ", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Equal(t, "abc", val)
}

func TestStringSchema_Trim_NonStringWithNil(t *testing.T) {
	schema := joi.String().Trim()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Lower(t *testing.T) {
	schema := joi.String().Lower()

	val, errs1 := schema.ValidateWithOpts("ABC", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs1)
	assert.Equal(t, "abc", val)

	val, errs2 := schema.ValidateWithOpts("ABC", joi.ValidateOptions{Path: joi.Ptr("field")}) // revalida com mesmo input?
	_ = val
	_ = errs2

	val, errs3 := schema.ValidateWithOpts("abc", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs3)
	assert.Equal(t, "abc", val)
}

func TestStringSchema_Lower_NonStringWithNil(t *testing.T) {
	schema := joi.String().Lower()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestStringSchema_Upper(t *testing.T) {
	schema := joi.String().Upper()

	val, errs1 := schema.ValidateWithOpts("abc", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs1)
	assert.Equal(t, "ABC", val)

	val, errs2 := schema.ValidateWithOpts("ABC", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs2)
	assert.Equal(t, "ABC", val)
}

func TestStringSchema_Upper_NonStringWithNil(t *testing.T) {
	schema := joi.String().Upper()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}
