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

// --- methods ---

func (s *BooleanSchema) True(msg ...string) *BooleanSchema {
	s.rules = append(s.rules, Rule{
		Name: string(BooleanMsgTrue),
		Msg:  pickMsg(BooleanMsgMap[BooleanMsgTrue], msg...),
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
		Msg:  pickMsg(BooleanMsgMap[BooleanMsgFalse], msg...),
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
	s.rules = append(s.rules, Rule{
		Name: string(BooleanMsgTruthy),
		Msg:  pickMsg(BooleanMsgMap[BooleanMsgTruthy]),
		Args: map[string]any{"truthy": values},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if valueInList(value, r.Args["truthy"].([]any)) {
				return value, nil
			}
			return value, &ValidationError{Path: path, Msg: r.Msg}
		},
	})
	return s
}

func (s *BooleanSchema) Falsy(values ...any) *BooleanSchema {
	s.rules = append(s.rules, Rule{
		Name: string(BooleanMsgFalsy),
		Msg:  pickMsg(BooleanMsgMap[BooleanMsgFalsy]),
		Args: map[string]any{"falsy": values},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			if valueInList(value, r.Args["falsy"].([]any)) {
				return value, nil
			}
			return value, &ValidationError{Path: path, Msg: r.Msg}
		},
	})
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
			Msg:  pickMsg(BooleanMsgMap[BooleanMsgBase], msg...),
			Fn: func(r Rule, path string, value any) (any, *ValidationError) {
				if value == nil {
					return value, nil
				}
				if _, ok := value.(bool); ok {
					return value, nil
				}
				return value, nil
			},
		}},
	}
	return s
}
