package joi_test

import (
	"regexp"
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestArraySchema_Base(t *testing.T) {
	schema := joi.Array()

	_, errs1 := schema.Validate([]any{"a", "b"})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("not-an-array")
	assert.NotEmpty(t, errs2)
}

func TestArraySchema_Min(t *testing.T) {
	schema := joi.Array().Min(2)

	_, errs1 := schema.ValidateWithOpts([]any{"a", "b"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts([]any{"a"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestArraySchema_Min_NonArrayInput(t *testing.T) {
	schema := joi.Array().Min(2)
	_, errs := schema.ValidateWithOpts("not-an-array", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
}

func TestArraySchema_Max(t *testing.T) {
	schema := joi.Array().Max(2)

	_, errs1 := schema.ValidateWithOpts([]any{"a", "b"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts([]any{"a", "b", "c"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestArraySchema_Max_NonArrayInput(t *testing.T) {
	schema := joi.Array().Max(2)
	_, errs := schema.ValidateWithOpts(123, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
}

func TestArraySchema_Length(t *testing.T) {
	schema := joi.Array().Length(2)

	_, errs1 := schema.ValidateWithOpts([]any{"a", "b"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts([]any{"a"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)

	_, errs3 := schema.ValidateWithOpts([]any{"a", "b", "c"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs3)
}

func TestArraySchema_Length_NonArrayInput(t *testing.T) {
	schema := joi.Array().Length(2)
	_, errs := schema.ValidateWithOpts(map[string]any{"x": 1}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
}

func TestArraySchema_Items(t *testing.T) {
	itemSchema := joi.Any[joi.Schema]().Valid([]any{"ok"})
	schema := joi.Array().Items(itemSchema)

	_, errs1 := schema.ValidateWithOpts([]any{"ok", "ok"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts([]any{"ok", "bad"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestArraySchema_DefaultAppliedOnNil(t *testing.T) {
	schema := joi.Array().Default([]any{"x"})
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Equal(t, []any{"x"}, val)
}

func TestArraySchema_ItemsValidationOnProvidedArray(t *testing.T) {
	re := regexp.MustCompile(`^ok$`)
	item := joi.String().Regex(re)
	schema := joi.Array().Items(item)

	got, errs := schema.ValidateWithOpts([]any{"ok", "bad"}, joi.ValidateOptions{Path: joi.Ptr("field")})

	assert.NotEmpty(t, errs)
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].String(), `validation error at "field[1]"`)
	assert.Contains(t, errs[0].String(), "fails to match the required pattern")

	arr, ok := got.([]any)
	assert.True(t, ok, "expected []any return")
	assert.Equal(t, []any{"ok", "bad"}, arr)
}

func TestArraySchema_DefaultNonArrayTriggersBaseError(t *testing.T) {
	schema := joi.Array().Default("not-an-array")
	got, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
	assert.Equal(t, "not-an-array", got)
}

func TestArraySchema_DefaultBranchIsHit(t *testing.T) {
	as := joi.Array()
	as.AnySchema.Default([]any{"__d__"})

	got, errs := as.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})

	arr, ok := got.([]any)

	assert.Empty(t, errs)
	assert.True(t, ok)
	assert.Equal(t, []any{"__d__"}, arr)
}

func TestArraySchema_BaseAllowsNil(t *testing.T) {
	schema := joi.Array()
	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}
