// joi_test/number_test.go

package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestNumberSchema_Base(t *testing.T) {
	schema := joi.Number()

	_, errs1 := schema.Validate("field", float64(10))
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "not-number")
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Min(t *testing.T) {
	schema := joi.Number().Min(5)

	_, errs1 := schema.Validate("field", float64(10))
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", float64(3))
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Max(t *testing.T) {
	schema := joi.Number().Max(5)

	_, errs1 := schema.Validate("field", float64(3))
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", float64(10))
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Integer(t *testing.T) {
	schema := joi.Number().Integer()

	_, errs1 := schema.Validate("field", float64(10))
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", float64(10.5))
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Positive(t *testing.T) {
	schema := joi.Number().Positive()

	_, errs1 := schema.Validate("field", float64(5))
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", float64(-3))
	assert.NotEmpty(t, errs2)
}

func TestNumberSchema_Negative(t *testing.T) {
	schema := joi.Number().Negative()

	_, errs1 := schema.Validate("field", float64(-5))
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", float64(3))
	assert.NotEmpty(t, errs2)
}
