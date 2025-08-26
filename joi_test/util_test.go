package joi_test

import (
	"testing"
	"time"

	"github.com/leandroluk/go-joi/joi"
	"github.com/stretchr/testify/assert"
)

func TestRenderTemplate(t *testing.T) {
	msg := joi.RenderTemplate(
		"{{#label}} at {{#path}} with {{#value}} / lim={{#limit}}",
		map[string]any{
			"label": "Field",
			"path":  "obj.a",
			"value": 42,
			"limit": 3,
		},
	)
	assert.Equal(t, "Field at obj.a with 42 / lim=3", msg)
}

func TestCoalesce(t *testing.T) {
	// string
	assert.Equal(t, "first", joi.Coalesce("", "first", "second"))
	// int
	assert.Equal(t, 10, joi.Coalesce(0, 0, 10, 20))
	// all zero -> zero
	assert.Equal(t, "", joi.Coalesce[string]())
}

func TestWhenNil(t *testing.T) {
	err := &joi.ValidationError{Path: "x", Msg: "nil!"}
	assert.Nil(t, joi.WhenNil(123, err))
	assert.Equal(t, err, joi.WhenNil(nil, err))
}

func TestPickSchemaMsg(t *testing.T) {
	assert.Equal(t, "override", joi.PickSchemaMsg("default", "override"))
	assert.Equal(t, "default", joi.PickSchemaMsg("default"))
	assert.Equal(t, "default", joi.PickSchemaMsg("default", "")) // empty override ignored
}

func TestValueInList(t *testing.T) {
	list := []any{"a", 1, true, []int{1, 2}}
	assert.True(t, joi.ValueInList("a", list))
	assert.True(t, joi.ValueInList(1, list))
	assert.True(t, joi.ValueInList(true, list))
	assert.True(t, joi.ValueInList([]int{1, 2}, list)) // DeepEqual
	assert.False(t, joi.ValueInList("missing", list))
}

func TestParseDate(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)

	// time.Time
	tm, ok := joi.ParseDate(now)
	assert.True(t, ok)
	assert.True(t, tm.Equal(now))

	// RFC3339 string (using PARSE_LAYOUT)
	s := now.Format(joi.PARSE_LAYOUT)
	tm, ok = joi.ParseDate(s)
	assert.True(t, ok)
	assert.True(t, tm.Equal(now))

	// unix int64
	sec := now.Unix()
	tm, ok = joi.ParseDate(sec)
	assert.True(t, ok)
	assert.Equal(t, sec, tm.Unix())

	// unix float64 (default for numbers in JSON)
	tm, ok = joi.ParseDate(float64(sec))
	assert.True(t, ok)
	assert.Equal(t, sec, tm.Unix())

	// invalid
	_, ok = joi.ParseDate("not-a-date")
	assert.False(t, ok)
	_, ok = joi.ParseDate(struct{}{})
	assert.False(t, ok)
}

func TestRunValidation_FullFlow(t *testing.T) {
	var rules []joi.Rule

	rules = append(rules, joi.Rule{
		Name: "transform",
		Msg:  "",
		Args: nil,
		Fn: func(r joi.Rule, path string, value any) (any, *joi.ValidationError) {
			return "B", nil
		},
	})

	rules = append(rules, joi.Rule{
		Name: "error_with_rule_msg",
		Msg:  "X {{#label}} {{#path}} {{#value}} lim={{#limit}}",
		Args: map[string]any{"limit": 3},
		Fn: func(r joi.Rule, path string, value any) (any, *joi.ValidationError) {
			return nil, &joi.ValidationError{Path: path, Msg: ""}
		},
	})

	rules = append(rules, joi.Rule{
		Name: "error_with_err_msg_and_transform",
		Msg:  "",
		Args: nil,
		Fn: func(r joi.Rule, path string, value any) (any, *joi.ValidationError) {
			return "C", &joi.ValidationError{Path: path, Msg: "E {{#label}} {{#path}} {{#value}}"}
		},
	})

	start := "A"
	got, errs := joi.RunValidation(rules, "MyLabel", "obj.a", start)

	assert.Equal(t, "C", got)

	assert.Len(t, errs, 2)

	assert.Contains(t, errs[0].String(), `validation error at "obj.a"`)
	assert.Contains(t, errs[0].String(), "X MyLabel obj.a B lim=3")

	assert.Contains(t, errs[1].String(), `validation error at "obj.a"`)
	assert.Contains(t, errs[1].String(), "E MyLabel obj.a B")
}
