package joi_test

import (
	"testing"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestAnySchema_Required(t *testing.T) {
	schema := joi.Any[joi.Schema]().Required()

	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{"nil value", nil, true},
		{"non-nil value", "ok", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errs := schema.ValidateWithOpts(tt.input, joi.ValidateOptions{Path: joi.Ptr("field")})
			assert.Equal(t, tt.wantErr, len(errs) > 0)
		})
	}
}

func TestAnySchema_Valid(t *testing.T) {
	schema := joi.Any[joi.Schema]().Valid([]any{"a", "b"})

	_, errs1 := schema.ValidateWithOpts("a", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("c", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs2)
}

func TestAnySchema_Invalid(t *testing.T) {
	schema := joi.Any[joi.Schema]().Invalid([]any{"x", "y"})

	_, errs1 := schema.ValidateWithOpts("x", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("z", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs2)
}

func TestAnySchema_Custom(t *testing.T) {
	schema := joi.Any[joi.Schema]().Custom(func(path string, v any) *joi.ValidationError {
		if v == "bad" {
			return &joi.ValidationError{Path: path, Msg: "nope"}
		}
		return nil
	})

	_, errs1 := schema.ValidateWithOpts("bad", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs1)

	_, errs2 := schema.ValidateWithOpts("good", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs2)
}

func TestAnySchema_Label(t *testing.T) {
	schema := joi.Any[joi.Schema]().Label("myField").Required()

	_, errs := schema.Validate(nil)
	assert.NotEmpty(t, errs)
	assert.Contains(t, errs[0].String(), "myField")
}

func TestAnySchema_Default(t *testing.T) {
	schema := joi.Any[joi.Schema]().Default("x")

	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Equal(t, "x", val)
}

func TestAnySchema_CustomNil(t *testing.T) {
	schema := joi.Any[joi.Schema]().Custom(nil)

	_, errs := schema.Validate("anything")
	assert.Empty(t, errs)
}

func TestAnySchema_InvalidNilValue(t *testing.T) {
	schema := joi.Any[joi.Schema]().Invalid([]any{"bad"})

	_, errs := schema.Validate(nil)
	assert.Empty(t, errs) // covers nil branch
}

func TestAnySchema_ValidNilValue(t *testing.T) {
	schema := joi.Any[joi.Schema]().Valid([]any{"ok"})

	_, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs) // covers nil branch
}
