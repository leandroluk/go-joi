// joi_test/array_test.go

package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestArraySchema_Base(t *testing.T) {
	schema := joi.Array()

	_, errs1 := schema.Validate("field", []any{"a", "b"})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "not-an-array")
	assert.NotEmpty(t, errs2)
}

func TestArraySchema_Min(t *testing.T) {
	schema := joi.Array().Min(2)

	_, errs1 := schema.Validate("field", []any{"a", "b"})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", []any{"a"})
	assert.NotEmpty(t, errs2)
}

func TestArraySchema_Max(t *testing.T) {
	schema := joi.Array().Max(2)

	_, errs1 := schema.Validate("field", []any{"a", "b"})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", []any{"a", "b", "c"})
	assert.NotEmpty(t, errs2)
}

func TestArraySchema_Length(t *testing.T) {
	schema := joi.Array().Length(2)

	_, errs1 := schema.Validate("field", []any{"a", "b"})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", []any{"a"})
	assert.NotEmpty(t, errs2)

	_, errs3 := schema.Validate("field", []any{"a", "b", "c"})
	assert.NotEmpty(t, errs3)
}

func TestArraySchema_Items(t *testing.T) {
	itemSchema := joi.Any[joi.Schema]().Valid([]any{"ok"})
	schema := joi.Array().Items(itemSchema)

	_, errs1 := schema.Validate("field", []any{"ok", "ok"})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", []any{"ok", "bad"})
	assert.NotEmpty(t, errs2)
}
