package joi

// --- messages ---

type BooleanMsg string

var (
	BooleanMsgBase   BooleanMsg = "boolean_base"
	BooleanMsgTrue   BooleanMsg = "boolean_true"
	BooleanMsgFalse  BooleanMsg = "boolean_false"
	BooleanMsgTruthy BooleanMsg = "boolean_truthy"
	BooleanMsgFalsy  BooleanMsg = "boolean_falsy"
)

var BooleanMsgMap = map[BooleanMsg]string{
	BooleanMsgBase:   "{{#label}} must be a boolean",
	BooleanMsgTrue:   "{{#label}} must be true",
	BooleanMsgFalse:  "{{#label}} must be false",
	BooleanMsgTruthy: "{{#label}} must be a truthy value",
	BooleanMsgFalsy:  "{{#label}} must be a falsy value",
}

// --- structs ---

type BooleanSchema struct {
	*AnySchema[*BooleanSchema]
}

var _ Schema = (*BooleanSchema)(nil)

// --- methods ---

func (s *BooleanSchema) True(msg ...string) *BooleanSchema {
	s.rules = append(s.rules, Rule{
		Name: string(BooleanMsgTrue),
		Msg:  PickSchemaMsg(BooleanMsgMap[BooleanMsgTrue], msg...),
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			b, ok := value.(bool)
			if !ok {
				return value, nil
			}
			if !b {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *BooleanSchema) False(msg ...string) *BooleanSchema {
	s.rules = append(s.rules, Rule{
		Name: string(BooleanMsgFalse),
		Msg:  PickSchemaMsg(BooleanMsgMap[BooleanMsgFalse], msg...),
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			b, ok := value.(bool)
			if !ok {
				return value, nil
			}
			if b {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *BooleanSchema) Truthy(values ...any) *BooleanSchema {
	rule := Rule{
		Name: string(BooleanMsgTruthy),
		Msg:  PickSchemaMsg(BooleanMsgMap[BooleanMsgTruthy]),
		Args: map[string]any{"truthy": values},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if ValueInList(value, r.Args["truthy"].([]any)) {
				return true, nil // convert to bool
			}
			return value, nil
		},
	}
	// insert before base
	s.rules = append([]Rule{rule}, s.rules...)
	return s
}

func (s *BooleanSchema) Falsy(values ...any) *BooleanSchema {
	rule := Rule{
		Name: string(BooleanMsgFalsy),
		Msg:  PickSchemaMsg(BooleanMsgMap[BooleanMsgFalsy]),
		Args: map[string]any{"falsy": values},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if ValueInList(value, r.Args["falsy"].([]any)) {
				return false, nil // convert to bool
			}
			return value, nil
		},
	}
	// insert before base
	s.rules = append([]Rule{rule}, s.rules...)
	return s
}

// --- constructor ---

func Boolean(msg ...string) *BooleanSchema {
	s := &BooleanSchema{}
	s.AnySchema = &AnySchema[*BooleanSchema]{
		self:  s,
		label: "value",
		rules: []Rule{{
			Name: string(BooleanMsgBase),
			Msg:  PickSchemaMsg(BooleanMsgMap[BooleanMsgBase], msg...),
			Fn: func(r Rule, path string, value any) (any, *ValidationError) {
				if value == nil {
					return value, nil
				}
				if _, ok := value.(bool); !ok {
					return value, &ValidationError{Path: path, Msg: r.Msg}
				}
				return value, nil
			},
		}},
	}
	return s
}
