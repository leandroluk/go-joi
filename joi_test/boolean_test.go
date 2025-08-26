package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestBooleanSchema_Base(t *testing.T) {
	schema := joi.Boolean()

	_, errs1 := schema.Validate("field", true)
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", false)
	assert.Empty(t, errs2)

	_, errs3 := schema.Validate("field", "not-bool")
	assert.NotEmpty(t, errs3)
}

func TestBooleanSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.Boolean()

	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestBooleanSchema_True(t *testing.T) {
	schema := joi.Boolean().True()

	_, errs1 := schema.Validate("field", true)
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", false)
	assert.NotEmpty(t, errs2)
}

func TestBooleanSchema_True_NonBoolInput(t *testing.T) {
	schema := joi.Boolean().True()

	val, errs := schema.Validate("field", "not-bool")
	assert.NotEmpty(t, errs)
	assert.Equal(t, "not-bool", val)
}

func TestBooleanSchema_False(t *testing.T) {
	schema := joi.Boolean().False()

	_, errs1 := schema.Validate("field", false)
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", true)
	assert.NotEmpty(t, errs2)
}

func TestBooleanSchema_False_NonBoolInput(t *testing.T) {
	schema := joi.Boolean().False()

	val, errs := schema.Validate("field", 123)
	assert.NotEmpty(t, errs)
	assert.Equal(t, 123, val)
}

func TestBooleanSchema_Truthy(t *testing.T) {
	schema := joi.Boolean().Truthy(1, "yes")

	_, errs1 := schema.Validate("field", 1)
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "yes")
	assert.Empty(t, errs2)

	_, errs3 := schema.Validate("field", 0)
	assert.NotEmpty(t, errs3)
}

func TestBooleanSchema_Falsy(t *testing.T) {
	schema := joi.Boolean().Falsy(0, "no")

	_, errs1 := schema.Validate("field", 0)
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "no")
	assert.Empty(t, errs2)

	_, errs3 := schema.Validate("field", 1)
	assert.NotEmpty(t, errs3)
}
