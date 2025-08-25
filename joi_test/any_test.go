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
			_, errs := schema.Validate("field", tt.input)
			assert.Equal(t, tt.wantErr, len(errs) > 0)
		})
	}
}

func TestAnySchema_Valid(t *testing.T) {
	schema := joi.Any[joi.Schema]().Valid([]any{"a", "b"})

	_, errs1 := schema.Validate("field", "a")
	assert.Empty(t, errs1)

	_, errs2 := schema.Validate("field", "c")
	assert.NotEmpty(t, errs2)
}

func TestAnySchema_Invalid(t *testing.T) {
	schema := joi.Any[joi.Schema]().Invalid([]any{"x", "y"})

	_, errs1 := schema.Validate("field", "x")
	assert.NotEmpty(t, errs1)

	_, errs2 := schema.Validate("field", "z")
	assert.Empty(t, errs2)
}

func TestAnySchema_Custom(t *testing.T) {
	schema := joi.Any[joi.Schema]().Custom(func(path string, v any) *joi.ValidationError {
		if v == "bad" {
			return &joi.ValidationError{Path: path, Msg: "nope"}
		}
		return nil
	})

	_, errs1 := schema.Validate("field", "bad")
	assert.NotEmpty(t, errs1)

	_, errs2 := schema.Validate("field", "good")
	assert.Empty(t, errs2)
}
