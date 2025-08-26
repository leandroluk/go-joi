package joi_test

import (
	"testing"
	"time"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestDateSchema_Base(t *testing.T) {
	schema := joi.Date()

	// time.Time
	_, errs := schema.Validate(time.Now())
	assert.Empty(t, errs)

	// ISO 8601 string
	_, errs = schema.Validate("2025-08-25T12:00:00Z")
	assert.Empty(t, errs)

	// Unix timestamp (int64)
	_, errs = schema.Validate(int64(1735185600))
	assert.Empty(t, errs)

	// Unix timestamp (float64, como JSON number)
	_, errs = schema.Validate(float64(1735185600))
	assert.Empty(t, errs)

	// Valor inválido
	_, errs = schema.Validate("not-a-date")
	assert.NotEmpty(t, errs)
}

func TestDateSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.Date()

	val, errs := schema.ValidateWithOpts(nil, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestDateSchema_Min(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Min(limit)

	_, errs := schema.ValidateWithOpts("2025-08-25T12:00:00Z", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)

	_, errs = schema.ValidateWithOpts("2025-08-26T12:00:00Z", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)

	_, errs = schema.ValidateWithOpts("2025-08-24T12:00:00Z", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
}

func TestDateSchema_Min_NonDateInput(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Min(limit)

	val, errs := schema.ValidateWithOpts("not-a-date", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
	assert.Equal(t, "not-a-date", val)
}

func TestDateSchema_Max(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Max(limit)

	_, errs := schema.ValidateWithOpts("2025-08-25T12:00:00Z", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)

	_, errs = schema.ValidateWithOpts("2025-08-24T12:00:00Z", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.Empty(t, errs)

	_, errs = schema.ValidateWithOpts("2025-08-26T12:00:00Z", joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
}

func TestDateSchema_Max_NonDateInput(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Max(limit)

	val, errs := schema.ValidateWithOpts(12345, joi.ValidateOptions{Path: joi.Ptr("field")})
	assert.NotEmpty(t, errs)
	assert.Equal(t, 12345, val)
}
