package joi

import (
	"strconv"
)

// --- messages ---

type ArrayMsg string

var (
	ArrayMsgBase   ArrayMsg = "array_base"
	ArrayMsgMin    ArrayMsg = "array_min"
	ArrayMsgMax    ArrayMsg = "array_max"
	ArrayMsgLength ArrayMsg = "array_length"
)

var ArrayMsgMap = map[ArrayMsg]string{
	ArrayMsgBase:   "{{#label}} must be an array",
	ArrayMsgMin:    "{{#label}} must contain at least {{#limit}} items",
	ArrayMsgMax:    "{{#label}} must contain less than or equal to {{#limit}} items",
	ArrayMsgLength: "{{#label}} must contain {{#limit}} items",
}

// --- structs ---

type ArraySchema struct {
	*AnySchema[*ArraySchema]
	itemsSchema Schema // schema dos itens (pode ser nil = aceita qualquer coisa)
}

// --- methods ---

func (s *ArraySchema) Items(schema Schema) *ArraySchema {
	s.itemsSchema = schema
	return s
}

func (s *ArraySchema) Min(limit int, msg ...string) *ArraySchema {
	s.rules = append(s.rules, Rule{
		Name: string(ArrayMsgMin),
		Msg:  PickSchemaMsg(ArrayMsgMap[ArrayMsgMin], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			arr, ok := value.([]any)
			if !ok {
				return value, nil // base cuida
			}
			if len(arr) < r.Args["limit"].(int) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *ArraySchema) Max(limit int, msg ...string) *ArraySchema {
	s.rules = append(s.rules, Rule{
		Name: string(ArrayMsgMax),
		Msg:  PickSchemaMsg(ArrayMsgMap[ArrayMsgMax], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			arr, ok := value.([]any)
			if !ok {
				return value, nil
			}
			if len(arr) > r.Args["limit"].(int) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *ArraySchema) Length(limit int, msg ...string) *ArraySchema {
	s.rules = append(s.rules, Rule{
		Name: string(ArrayMsgLength),
		Msg:  PickSchemaMsg(ArrayMsgMap[ArrayMsgLength], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			arr, ok := value.([]any)
			if !ok {
				return value, nil
			}
			if len(arr) != r.Args["limit"].(int) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *ArraySchema) Validate(path string, value any) (any, []ValidationError) {
	if value == nil && s.defaultValue != nil {
		value = s.defaultValue.value
	}

	val, errs := RunValidation(s.rules, Coalesce(s.label, path, "value"), path, value)

	arr, ok := val.([]any)
	if !ok {
		return val, errs
	}

	if s.itemsSchema != nil {
		newArr := make([]any, len(arr))
		for i, v := range arr {
			itemPath := path + "[" + strconv.Itoa(i) + "]"
			parsed, itemErrs := s.itemsSchema.Validate(itemPath, v)
			if len(itemErrs) > 0 {
				errs = append(errs, itemErrs...)
			}
			newArr[i] = parsed
		}
		return newArr, errs
	}

	return arr, errs
}

// --- constructor ---

func Array(msg ...string) *ArraySchema {
	base := Rule{
		Name: string(ArrayMsgBase),
		Msg:  PickSchemaMsg(ArrayMsgMap[ArrayMsgBase], msg...),
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if value == nil {
				return value, nil
			}
			if _, ok := value.([]any); !ok {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	}
	s := &ArraySchema{}
	s.AnySchema = &AnySchema[*ArraySchema]{
		self:  s,
		label: "value",
		rules: []Rule{base},
	}
	return s
}
