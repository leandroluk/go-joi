package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestObjectSchema_Base(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{})

	_, errs1 := schema.Validate(map[string]any{})
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("not-object")
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{})
	val, errs := schema.Validate(nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestObjectSchema_Min(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{}).Unknown(true).Min(2)

	_, errs1 := schema.ValidateWithOpts(map[string]any{"a": 1, "b": 2}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(map[string]any{"a": 1}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Max(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{}).Unknown(true).Max(1)

	_, errs1 := schema.ValidateWithOpts(map[string]any{"a": 1}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(map[string]any{"a": 1, "b": 2}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Length(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{}).Unknown(true).Length(2)

	_, errs1 := schema.ValidateWithOpts(map[string]any{"a": 1, "b": 2}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(map[string]any{"a": 1}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Fields(t *testing.T) {
	fields := map[string]joi.Schema{
		"a": joi.Any[joi.Schema]().Required(),
		"b": joi.Any[joi.Schema](),
	}
	schema := joi.Object(fields)

	_, errs1 := schema.ValidateWithOpts(map[string]any{"a": "ok"}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts(map[string]any{}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestObjectSchema_Unknown(t *testing.T) {
	fields := map[string]joi.Schema{
		"a": joi.Any[joi.Schema](),
	}
	schema := joi.Object(fields).Unknown(false)

	_, errs1 := schema.ValidateWithOpts(map[string]any{"a": 1, "x": 2}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs1)

	schema2 := joi.Object(fields).Unknown(true)
	_, errs2 := schema2.ValidateWithOpts(map[string]any{"a": 1, "x": 2}, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs2)
}

func TestObjectSchema_ChildErrors_WhenKeyExists(t *testing.T) {
	fields := map[string]joi.Schema{
		"a": joi.String().Min(3),
	}
	schema := joi.Object(fields)

	_, errs := schema.ValidateWithOpts(map[string]any{"a": "ab"}, joi.ValidateOptions{Path: joi.Ptr("obj")})
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0].String(), `validation error at "obj.a"`)
}

func TestObjectSchema_Rules_NonMapBranches(t *testing.T) {
	schema := joi.Object(map[string]joi.Schema{}).
		Min(2).
		Max(3).
		Length(1)

	val, errs := schema.ValidateWithOpts("not-an-object", joi.ValidateOptions{Path: joi.Ptr("obj")})

	assert.Equal(t, "not-an-object", val)
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0].String(), "must be an object")
}
