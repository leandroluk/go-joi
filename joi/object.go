package joi

// --- messages ---

type ObjectMsg string

var (
	ObjectMsgBase    ObjectMsg = "object_base"
	ObjectMsgMin     ObjectMsg = "object_min"
	ObjectMsgMax     ObjectMsg = "object_max"
	ObjectMsgLength  ObjectMsg = "object_length"
	ObjectMsgUnknown ObjectMsg = "object_unknown"
)

var ObjectMsgMap = map[ObjectMsg]string{
	ObjectMsgBase:    "{{#label}} must be an object",
	ObjectMsgMin:     "{{#label}} must have at least {{#limit}} keys",
	ObjectMsgMax:     "{{#label}} must have less than or equal to {{#limit}} keys",
	ObjectMsgLength:  "{{#label}} must have {{#limit}} keys",
	ObjectMsgUnknown: "{{#label}} contains unknown key '{{#key}}'",
}

// --- structs ---

type ObjectSchema struct {
	*AnySchema[*ObjectSchema]
	fields  map[string]Schema
	unknown bool
}

// --- methods ---

func (s *ObjectSchema) Unknown(allow bool) *ObjectSchema {
	s.unknown = allow
	return s
}

func (s *ObjectSchema) Min(limit int, msg ...string) *ObjectSchema {
	s.rules = append(s.rules, Rule{
		Name: string(ObjectMsgMin),
		Msg:  PickSchemaMsg(ObjectMsgMap[ObjectMsgMin], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			m, ok := value.(map[string]any)
			if !ok {
				return value, nil
			}
			if len(m) < r.Args["limit"].(int) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *ObjectSchema) Max(limit int, msg ...string) *ObjectSchema {
	s.rules = append(s.rules, Rule{
		Name: string(ObjectMsgMax),
		Msg:  PickSchemaMsg(ObjectMsgMap[ObjectMsgMax], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			m, ok := value.(map[string]any)
			if !ok {
				return value, nil
			}
			if len(m) > r.Args["limit"].(int) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *ObjectSchema) Length(limit int, msg ...string) *ObjectSchema {
	s.rules = append(s.rules, Rule{
		Name: string(ObjectMsgLength),
		Msg:  PickSchemaMsg(ObjectMsgMap[ObjectMsgLength], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			m, ok := value.(map[string]any)
			if !ok {
				return value, nil
			}
			if len(m) != r.Args["limit"].(int) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *ObjectSchema) Validate(value any) (any, []ValidationError) {
	return s.ValidateWithOpts(value, ValidateOptions{})
}

func (s *ObjectSchema) ValidateWithOpts(value any, opts ValidateOptions) (any, []ValidationError) {
	var path string
	if opts.Path != nil {
		path = *opts.Path
	}

	val, errs := RunValidation(s.rules, Coalesce(path, "value"), path, value)

	if val == nil {
		return nil, errs
	}

	m, ok := val.(map[string]any)
	if !ok {
		return val, errs
	}

	var childErrs []ValidationError
	parsed := make(map[string]any)

	for k, schema := range s.fields {
		if v, exists := m[k]; exists {
			// valida campo existente
			childPath := path + "." + k
			parsedVal, ce := schema.ValidateWithOpts(v, ValidateOptions{Path: &childPath})
			if len(ce) > 0 {
				childErrs = append(childErrs, ce...)
			}
			parsed[k] = parsedVal
		} else {
			// campo ausente → valida contra nil (pra Required() funcionar)
			childPath := path + "." + k
			_, ce := schema.ValidateWithOpts(nil, ValidateOptions{Path: &childPath})
			if len(ce) > 0 {
				childErrs = append(childErrs, ce...)
			}
		}
	}

	if !s.unknown {
		for k := range m {
			if _, ok := s.fields[k]; !ok {
				childErrs = append(childErrs, ValidationError{
					Path: path + "." + k,
					Msg:  RenderTemplate(ObjectMsgMap[ObjectMsgUnknown], map[string]any{"label": path, "key": k}),
				})
			}
		}
	} else {
		for k, v := range m {
			if _, ok := s.fields[k]; !ok {
				parsed[k] = v
			}
		}
	}

	if len(childErrs) > 0 {
		errs = append(errs, childErrs...)
	}
	return parsed, errs
}

// --- constructor ---

func Object(fields map[string]Schema, msg ...string) *ObjectSchema {
	s := &ObjectSchema{fields: fields, unknown: false}
	s.AnySchema = &AnySchema[*ObjectSchema]{
		self:  s,
		label: "value",
		rules: []Rule{{
			Name: string(ObjectMsgBase),
			Msg:  PickSchemaMsg(ObjectMsgMap[ObjectMsgBase], msg...),
			Fn: func(r Rule, path string, value any) (any, *ValidationError) {
				if value == nil {
					return value, nil
				}
				if _, ok := value.(map[string]any); !ok {
					return value, &ValidationError{Path: path, Msg: r.Msg}
				}
				return value, nil
			},
		}},
	}
	return s
}
