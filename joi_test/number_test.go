package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestNumberSchema_Base(t *testing.T) {
	schema := joi.Number()

	_, errs1 := schema.ValidateWithOpts(float64(10), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("not-number", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.Number()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestNumberSchema_Min(t *testing.T) {
	schema := joi.Number().Min(5)

	_, errs1 := schema.ValidateWithOpts(float64(10), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(float64(3), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Min_NonNumberInput_WithNil(t *testing.T) {
	schema := joi.Number().Min(5)
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestNumberSchema_Max(t *testing.T) {
	schema := joi.Number().Max(5)

	_, errs1 := schema.ValidateWithOpts(float64(3), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(float64(10), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Max_NonNumberInput_WithNil(t *testing.T) {
	schema := joi.Number().Max(5)
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestNumberSchema_Integer(t *testing.T) {
	schema := joi.Number().Integer()

	_, errs1 := schema.ValidateWithOpts(float64(10), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(float64(10.5), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Integer_NonNumberInput_WithNil(t *testing.T) {
	schema := joi.Number().Integer()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestNumberSchema_Positive(t *testing.T) {
	schema := joi.Number().Positive()

	_, errs1 := schema.ValidateWithOpts(float64(5), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(float64(-3), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Positive_NonNumberInput_WithNil(t *testing.T) {
	schema := joi.Number().Positive()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestNumberSchema_Negative(t *testing.T) {
	schema := joi.Number().Negative()

	_, errs1 := schema.ValidateWithOpts(float64(-5), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(float64(3), joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Negative_NonNumberInput_WithNil(t *testing.T) {
	schema := joi.Number().Negative()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}
