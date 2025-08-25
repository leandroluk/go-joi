package joi

// --- messages ---

type AnyMsg string

var (
	AnyMsgDefault  AnyMsg = "any_default"
	AnyMsgCustom   AnyMsg = "any_custom"
	AnyMsgRequired AnyMsg = "any_required"
	AnyMsgInvalid  AnyMsg = "any_invalid"
	AnyMsgValid    AnyMsg = "any_valid"
)

var AnyMsgMap = map[AnyMsg]string{
	AnyMsgDefault:  "{{#label}} threw an error when running default method",
	AnyMsgCustom:   "{{#label}} failed custom validation because {{#message}}",
	AnyMsgRequired: "{{#label}} is required",
	AnyMsgInvalid:  "{{#label}} contains an invalid value",
	AnyMsgValid:    "{{#label}} must be one of {{#valid}}",
}

// --- structs ---

type DefaultValue struct {
	value any
}

type AnySchema[T any] struct {
	label        string
	rules        []Rule
	defaultValue *DefaultValue
	self         T
}

// --- methods ---

func (s *AnySchema[T]) Label(label string) T {
	s.label = label
	return s.self
}

func (s *AnySchema[T]) Default(value any) T {
	s.defaultValue = &DefaultValue{value: value}
	return s.self
}

func (s *AnySchema[T]) Custom(fn func(path string, value any) *ValidationError, msg ...string) T {
	name := string(AnyMsgCustom)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  pickMsg(AnyMsgMap[AnyMsgCustom], msg...),
		Args: map[string]any{"message": ""},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if fn == nil {
				return value, nil
			}
			if err := fn(path, value); err != nil {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s.self
}

func (s *AnySchema[T]) Required(msg ...string) T {
	name := string(AnyMsgRequired)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  pickMsg(AnyMsgMap[AnyMsgRequired], msg...),
		Args: map[string]any{},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			return value, whenNil(value, &ValidationError{Path: path, Msg: r.Msg})
		},
	})
	return s.self
}

func (s *AnySchema[T]) Invalid(disallowed []any, msg ...string) T {
	name := string(AnyMsgInvalid)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  pickMsg(AnyMsgMap[AnyMsgInvalid], msg...),
		Args: map[string]any{"invalid": disallowed},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if value == nil {
				return value, nil
			}
			list, _ := r.Args["invalid"].([]any)
			if valueInList(value, list) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s.self
}

func (s *AnySchema[T]) Valid(allowed []any, msg ...string) T {
	name := string(AnyMsgValid)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  pickMsg(AnyMsgMap[AnyMsgValid], msg...),
		Args: map[string]any{"valid": allowed},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if value == nil {
				return value, nil
			}
			list, _ := r.Args["valid"].([]any)
			if valueInList(value, list) {
				return value, nil
			}
			return value, &ValidationError{Path: path, Msg: r.Msg}
		},
	})
	return s.self
}

func (s *AnySchema[T]) Validate(path string, value any) (any, []ValidationError) {
	if value == nil && s.defaultValue != nil {
		value = s.defaultValue.value
	}
	return runValidation(s.rules, coalesce(s.label, path, "value"), path, value)
}

// --- constructor ---

func Any[T Schema]() *AnySchema[T] {
	s := &AnySchema[T]{label: "value", rules: make([]Rule, 0)}
	s.self = any(s).(T)
	return s
}
