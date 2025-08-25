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
		Msg:  pickMsg(ObjectMsgMap[ObjectMsgMin], msg...),
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
		Msg:  pickMsg(ObjectMsgMap[ObjectMsgMax], msg...),
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
		Msg:  pickMsg(ObjectMsgMap[ObjectMsgLength], msg...),
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

// --- validation ---

func (s *ObjectSchema) Validate(path string, value any) (any, []ValidationError) {
	if value == nil {
		return nil, nil
	}
	m, ok := value.(map[string]any)
	if !ok {
		return value, []ValidationError{{Path: path, Msg: ObjectMsgMap[ObjectMsgBase]}}
	}

	var errs []ValidationError
	parsed := make(map[string]any)

	val, valErrs := runValidation(s.rules, coalesce(path, "value"), path, value)
	if len(valErrs) > 0 {
		errs = append(errs, valErrs...)
	}
	if obj, ok := val.(map[string]any); ok {
		m = obj
	}

	for k, schema := range s.fields {
		if v, exists := m[k]; exists {
			childPath := path + "." + k
			parsedVal, childErrs := schema.Validate(childPath, v)
			if len(childErrs) > 0 {
				errs = append(errs, childErrs...)
			}
			parsed[k] = parsedVal
		} else {
			_, childErrs := schema.Validate(path+"."+k, nil)
			if len(childErrs) > 0 {
				errs = append(errs, childErrs...)
			}
		}
	}

	if !s.unknown {
		for k := range m {
			if _, ok := s.fields[k]; !ok {
				errs = append(errs, ValidationError{
					Path: path + "." + k,
					Msg:  renderTemplate(ObjectMsgMap[ObjectMsgUnknown], map[string]any{"label": path, "key": k}),
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
			Msg:  pickMsg(ObjectMsgMap[ObjectMsgBase], msg...),
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
