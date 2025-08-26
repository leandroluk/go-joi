package joi

import (
	"net/mail"
	"regexp"
	"strings"
)

// --- messages ---

type StringMsg string

var (
	StringMsgBase   StringMsg = "string_base"
	StringMsgMin    StringMsg = "string_min"
	StringMsgMax    StringMsg = "string_max"
	StringMsgLength StringMsg = "string_length"
	StringMsgRegex  StringMsg = "string_regex"
	StringMsgTrim   StringMsg = "string_trim"
	StringMsgLower  StringMsg = "string_lower"
	StringMsgUpper  StringMsg = "string_upper"
	StringMsgEmail  StringMsg = "string.email"
)

var StringMsgMap = map[StringMsg]string{
	StringMsgBase:   "{{#label}} must be a string",
	StringMsgMin:    "{{#label}} length must be at least {{#limit}} characters long",
	StringMsgMax:    "{{#label}} length must be less than or equal to {{#limit}} characters long",
	StringMsgLength: "{{#label}} length must be {{#limit}} characters long",
	StringMsgRegex:  "{{#label}} with value {{#value}} fails to match the required pattern",
	StringMsgTrim:   "{{#label}} must be a trimmed string",
	StringMsgLower:  "{{#label}} must be a lowercase string",
	StringMsgUpper:  "{{#label}} must be an uppercase string",
	StringMsgEmail:  "{{#label}} must be a valid email",
}

// --- structs ---

type StringSchema struct {
	*AnySchema[*StringSchema]
}

// --- methods ---

func (s *StringSchema) Min(limit int, msg ...string) *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgMin),
		Msg:  PickSchemaMsg(StringMsgMap[StringMsgMin], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			if len(str) < limit {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *StringSchema) Max(limit int, msg ...string) *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgMax),
		Msg:  PickSchemaMsg(StringMsgMap[StringMsgMax], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			if len(str) > limit {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *StringSchema) Length(limit int, msg ...string) *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgLength),
		Msg:  PickSchemaMsg(StringMsgMap[StringMsgLength], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			if len(str) != limit {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *StringSchema) Regex(re *regexp.Regexp, msg ...string) *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgRegex),
		Msg:  PickSchemaMsg(StringMsgMap[StringMsgRegex], msg...),
		Args: map[string]any{"pattern": re.String()},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			if !re.MatchString(str) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *StringSchema) Email(msg ...string) *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgEmail),
		Msg:  PickSchemaMsg(StringMsgMap[StringMsgEmail], msg...),
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			if str == "" {
				// Deixa "required" para outra regra; email vazio passa aqui.
				return value, nil
			}
			addr, err := mail.ParseAddress(str)
			if err != nil || addr.Address != str {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *StringSchema) Trim() *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgTrim),
		Msg:  "",
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			return strings.TrimSpace(str), nil
		},
	})
	return s
}

func (s *StringSchema) Lower() *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgLower),
		Msg:  StringMsgMap[StringMsgLower],
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			if str != strings.ToLower(str) {
				return strings.ToLower(str), &ValidationError{Path: path, Msg: r.Msg}
			}
			return str, nil
		},
	})
	return s
}

func (s *StringSchema) Upper() *StringSchema {
	s.rules = append(s.rules, Rule{
		Name: string(StringMsgUpper),
		Msg:  StringMsgMap[StringMsgUpper],
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			str, ok := value.(string)
			if !ok {
				return value, nil
			}
			if str != strings.ToUpper(str) {
				return strings.ToUpper(str), &ValidationError{Path: path, Msg: r.Msg}
			}
			return str, nil
		},
	})
	return s
}

// --- constructor ---

func String(msg ...string) *StringSchema {
	s := &StringSchema{}
	s.AnySchema = &AnySchema[*StringSchema]{
		self:  s,
		label: "value",
		rules: []Rule{{
			Name: string(StringMsgBase),
			Msg:  PickSchemaMsg(StringMsgMap[StringMsgBase], msg...),
			Fn: func(r Rule, path string, value any) (any, *ValidationError) {
				if value == nil {
					return value, nil
				}
				if _, ok := value.(string); !ok {
					return value, &ValidationError{Path: path, Msg: r.Msg}
				}
				return value, nil
			},
		}},
	}
	return s
}
