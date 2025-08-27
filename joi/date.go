package joi

import (
	"time"
)

// --- messages ---

type DateMsg string

var (
	DateMsgBase  DateMsg = "date_base"
	DateMsgMin   DateMsg = "date_min"
	DateMsgMax   DateMsg = "date_max"
	PARSE_LAYOUT         = time.RFC3339 // default ISO 8601
)

var DateMsgMap = map[DateMsg]string{
	DateMsgBase: "{{#label}} must be a valid date",
	DateMsgMin:  "{{#label}} must be larger than or equal to {{#limit}}",
	DateMsgMax:  "{{#label}} must be less than or equal to {{#limit}}",
}

// --- structs ---

type DateSchema struct {
	*AnySchema[*DateSchema]
}

var _ Schema = (*DateSchema)(nil)

// --- methods ---

func (s *DateSchema) Min(limit time.Time, msg ...string) *DateSchema {
	s.rules = append(s.rules, Rule{
		Name: string(DateMsgMin),
		Msg:  PickSchemaMsg(DateMsgMap[DateMsgMin], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			t, ok := ParseDate(value)
			if !ok {
				return value, nil // base handles error
			}
			if t.Before(r.Args["limit"].(time.Time)) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return t, nil
		},
	})
	return s
}

func (s *DateSchema) Max(limit time.Time, msg ...string) *DateSchema {
	s.rules = append(s.rules, Rule{
		Name: string(DateMsgMax),
		Msg:  PickSchemaMsg(DateMsgMap[DateMsgMax], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			t, ok := ParseDate(value)
			if !ok {
				return value, nil
			}
			if t.After(r.Args["limit"].(time.Time)) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return t, nil
		},
	})
	return s
}

// --- constructor ---

func Date(msg ...string) *DateSchema {
	s := &DateSchema{}
	s.AnySchema = &AnySchema[*DateSchema]{
		self:  s,
		label: "value",
		rules: []Rule{{
			Name: string(DateMsgBase),
			Msg:  PickSchemaMsg(DateMsgMap[DateMsgBase], msg...),
			Fn: func(r Rule, path string, value any) (any, *ValidationError) {
				if value == nil {
					return value, nil
				}
				if t, ok := ParseDate(value); ok {
					return t, nil
				}
				return value, &ValidationError{Path: path, Msg: r.Msg}
			},
		}},
	}
	return s
}
