package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestBooleanSchema_Base(t *testing.T) {
	schema := joi.Boolean()

	_, errs1 := schema.Validate(true)
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate(false)
	assert.Empty(t, errs2)

	_, errs3 := schema.Validate("not-bool")
	assert.NotEmpty(t, errs3)
}

func TestBooleanSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.Boolean()

	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestBooleanSchema_True(t *testing.T) {
	schema := joi.Boolean().True()

	_, errs1 := schema.ValidateWithOpts(true, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(false, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestBooleanSchema_True_NonBoolInput(t *testing.T) {
	schema := joi.Boolean().True()

	val, errs := schema.ValidateWithOpts("not-bool", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
	assert.Equal(t, "not-bool", val)
}

func TestBooleanSchema_False(t *testing.T) {
	schema := joi.Boolean().False()

	_, errs1 := schema.ValidateWithOpts(false, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(true, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestBooleanSchema_False_NonBoolInput(t *testing.T) {
	schema := joi.Boolean().False()

	val, errs := schema.ValidateWithOpts(123, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
	assert.Equal(t, 123, val)
}

func TestBooleanSchema_Truthy(t *testing.T) {
	schema := joi.Boolean().Truthy(1, "yes")

	_, errs1 := schema.ValidateWithOpts(1, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("yes", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs2)

	_, errs3 := schema.ValidateWithOpts(0, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs3)
}

func TestBooleanSchema_Falsy(t *testing.T) {
	schema := joi.Boolean().Falsy(0, "no")

	_, errs1 := schema.ValidateWithOpts(0, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("no", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs2)

	_, errs3 := schema.ValidateWithOpts(1, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs3)
}
