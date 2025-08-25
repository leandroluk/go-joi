// joi_test/object_test.go

package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestObjectSchema_Base(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{})

	_, errs1 := schema.Validate("field", map[string]any{})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "not-object")
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Min(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{}).Unknown(true).Min(2)

	_, errs1 := schema.Validate("field", map[string]any{"a": 1, "b": 2})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", map[string]any{"a": 1})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Max(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{}).Unknown(true).Max(1)

	_, errs1 := schema.Validate("field", map[string]any{"a": 1})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", map[string]any{"a": 1, "b": 2})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Length(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{}).Unknown(true).Length(2)

	_, errs1 := schema.Validate("field", map[string]any{"a": 1, "b": 2})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", map[string]any{"a": 1})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Fields(t *testing.T) {
	fields := map[string]joi.Schema{
		"a": joi.Any[joi.Schema]().Required(),
		"b": joi.Any[joi.Schema](),
	}
	schema := joi.Object(fields)

	_, errs1 := schema.Validate("field", map[string]any{"a": "ok"})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", map[string]any{})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Unknown(t *testing.T) {
	fields := map[string]joi.Schema{
		"a": joi.Any[joi.Schema](),
	}
	schema := joi.Object(fields).Unknown(false)

	_, errs1 := schema.Validate("field", map[string]any{"a": 1, "x": 2})
	assert.NotEmpty(t, errs1)

	schema2 := joi.Object(fields).Unknown(true)
	_, errs2 := schema2.Validate("field", map[string]any{"a": 1, "x": 2})
	assert.Empty(t, errs2)
}
