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
	path         string
	label        string
	rules        []Rule
	defaultValue *DefaultValue
	self         T
}

var _ Schema = (*AnySchema[AnySchema[any]])(nil)

// --- methods ---

func (s *AnySchema[T]) Label(label string) *AnySchema[T] {
	s.label = label
	return s
}

func (s *AnySchema[T]) Default(value any) *AnySchema[T] {
	s.defaultValue = &DefaultValue{value: value}
	return s
}

func (s *AnySchema[T]) Custom(fn func(path string, value any) *ValidationError, msg ...string) *AnySchema[T] {
	name := string(AnyMsgCustom)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  PickSchemaMsg(AnyMsgMap[AnyMsgCustom], msg...),
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
	return s
}

func (s *AnySchema[T]) Required(msg ...string) *AnySchema[T] {
	name := string(AnyMsgRequired)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  PickSchemaMsg(AnyMsgMap[AnyMsgRequired], msg...),
		Args: map[string]any{},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			return value, WhenNil(value, &ValidationError{Path: path, Msg: r.Msg})
		},
	})
	return s
}

func (s *AnySchema[T]) Invalid(disallowed []any, msg ...string) *AnySchema[T] {
	name := string(AnyMsgInvalid)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  PickSchemaMsg(AnyMsgMap[AnyMsgInvalid], msg...),
		Args: map[string]any{"invalid": disallowed},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if value == nil {
				return value, nil
			}
			list, _ := r.Args["invalid"].([]any)
			if ValueInList(value, list) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *AnySchema[T]) Valid(allowed []any, msg ...string) *AnySchema[T] {
	name := string(AnyMsgValid)
	s.rules = append(s.rules, Rule{
		Name: name,
		Msg:  PickSchemaMsg(AnyMsgMap[AnyMsgValid], msg...),
		Args: map[string]any{"valid": allowed},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if value == nil {
				return value, nil
			}
			list, _ := r.Args["valid"].([]any)
			if ValueInList(value, list) {
				return value, nil
			}
			return value, &ValidationError{Path: path, Msg: r.Msg}
		},
	})
	return s
}

func (s *AnySchema[T]) Validate(value any) (any, []ValidationError) {
	return s.ValidateWithOpts(value, ValidateOptions{})
}

func (s *AnySchema[T]) ValidateWithOpts(value any, opts ValidateOptions) (any, []ValidationError) {
	if value == nil && s.defaultValue != nil {
		value = s.defaultValue.value
	}

	if opts.Path != nil {
		s.path = *opts.Path
	}

	return RunValidation(s.rules, Coalesce(s.label, s.path, "value"), s.path, value)
}

// --- constructor ---

func Any[T Schema]() *AnySchema[T] {
	s := &AnySchema[T]{label: "value", rules: make([]Rule, 0)}
	s.self = any(s).(T)
	return s
}
