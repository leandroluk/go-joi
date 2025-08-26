package joi_test

import (
	"testing"
	"time"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestDateSchema_Base(t *testing.T) {
	schema := joi.Date()

	// time.Time direto
	_, errs := schema.Validate("field", time.Now())
	assert.Empty(t, errs)

	// ISO 8601 string
	_, errs = schema.Validate("field", "2025-08-25T12:00:00Z")
	assert.Empty(t, errs)

	// Unix timestamp (int64)
	_, errs = schema.Validate("field", int64(1735185600))
	assert.Empty(t, errs)

	// Unix timestamp (float64, como JSON number)
	_, errs = schema.Validate("field", float64(1735185600))
	assert.Empty(t, errs)

	// Valor inválido
	_, errs = schema.Validate("field", "not-a-date")
	assert.NotEmpty(t, errs)
}

func TestDateSchema_Base_AllowsNil(t *testing.T) {
	schema := joi.Date()

	val, errs := schema.Validate("field", nil)
	assert.Empty(t, errs)
	assert.Nil(t, val)
}

func TestDateSchema_Min(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Min(limit)

	// Igual ao limite → ok
	_, errs := schema.Validate("field", "2025-08-25T12:00:00Z")
	assert.Empty(t, errs)

	// Depois do limite → ok
	_, errs = schema.Validate("field", "2025-08-26T12:00:00Z")
	assert.Empty(t, errs)

	// Antes do limite → erro
	_, errs = schema.Validate("field", "2025-08-24T12:00:00Z")
	assert.NotEmpty(t, errs)
}

func TestDateSchema_Min_NonDateInput(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Min(limit)

	val, errs := schema.Validate("field", "not-a-date")
	assert.NotEmpty(t, errs)
	assert.Equal(t, "not-a-date", val)
}

func TestDateSchema_Max(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Max(limit)

	// Igual ao limite → ok
	_, errs := schema.Validate("field", "2025-08-25T12:00:00Z")
	assert.Empty(t, errs)

	// Antes do limite → ok
	_, errs = schema.Validate("field", "2025-08-24T12:00:00Z")
	assert.Empty(t, errs)

	// Depois do limite → erro
	_, errs = schema.Validate("field", "2025-08-26T12:00:00Z")
	assert.NotEmpty(t, errs)
}

func TestDateSchema_Max_NonDateInput(t *testing.T) {
	limit := time.Date(2025, 8, 25, 12, 0, 0, 0, time.UTC)
	schema := joi.Date().Max(limit)

	val, errs := schema.Validate("field", 12345)
	assert.NotEmpty(t, errs)
	assert.Equal(t, 12345, val)
}
