package joi

// --- messages ---

type NumberMsg string

var (
	NumberMsgBase     NumberMsg = "number_base"
	NumberMsgMin      NumberMsg = "number_min"
	NumberMsgMax      NumberMsg = "number_max"
	NumberMsgInteger  NumberMsg = "number_integer"
	NumberMsgPositive NumberMsg = "number_positive"
	NumberMsgNegative NumberMsg = "number_negative"
)

var NumberMsgMap = map[NumberMsg]string{
	NumberMsgBase:     "{{#label}} must be a number",
	NumberMsgMin:      "{{#label}} must be larger than or equal to {{#limit}}",
	NumberMsgMax:      "{{#label}} must be less than or equal to {{#limit}}",
	NumberMsgInteger:  "{{#label}} must be an integer",
	NumberMsgPositive: "{{#label}} must be a positive number",
	NumberMsgNegative: "{{#label}} must be a negative number",
}

// --- structs ---

type NumberSchema struct {
	*AnySchema[*NumberSchema]
}

// --- methods ---

func (s *NumberSchema) Min(limit float64, msg ...string) *NumberSchema {
	s.rules = append(s.rules, Rule{
		Name: string(NumberMsgMin),
		Msg:  PickSchemaMsg(NumberMsgMap[NumberMsgMin], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			num, ok := value.(float64)
			if !ok {
				return value, nil // number_base cuida disso
			}
			if num < limit {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *NumberSchema) Max(limit float64, msg ...string) *NumberSchema {
	s.rules = append(s.rules, Rule{
		Name: string(NumberMsgMax),
		Msg:  PickSchemaMsg(NumberMsgMap[NumberMsgMax], msg...),
		Args: map[string]any{"limit": limit},
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			num, ok := value.(float64)
			if !ok {
				return value, nil
			}
			if num > limit {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *NumberSchema) Integer(msg ...string) *NumberSchema {
	s.rules = append(s.rules, Rule{
		Name: string(NumberMsgInteger),
		Msg:  PickSchemaMsg(NumberMsgMap[NumberMsgInteger], msg...),
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			num, ok := value.(float64)
			if !ok {
				return value, nil
			}
			if num != float64(int64(num)) {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return int64(num), nil
		},
	})
	return s
}

func (s *NumberSchema) Positive(msg ...string) *NumberSchema {
	s.rules = append(s.rules, Rule{
		Name: string(NumberMsgPositive),
		Msg:  PickSchemaMsg(NumberMsgMap[NumberMsgPositive], msg...),
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			num, ok := value.(float64)
			if !ok {
				return value, nil
			}
			if num <= 0 {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

func (s *NumberSchema) Negative(msg ...string) *NumberSchema {
	s.rules = append(s.rules, Rule{
		Name: string(NumberMsgNegative),
		Msg:  PickSchemaMsg(NumberMsgMap[NumberMsgNegative], msg...),
		Fn: func(r Rule, path string, value any) (any, *ValidationError) {
			num, ok := value.(float64)
			if !ok {
				return value, nil
			}
			if num >= 0 {
				return value, &ValidationError{Path: path, Msg: r.Msg}
			}
			return value, nil
		},
	})
	return s
}

// --- constructor ---

func Number(msg ...string) *NumberSchema {
	s := &NumberSchema{}
	s.AnySchema = &AnySchema[*NumberSchema]{
		self:  s,
		label: "value",
		rules: []Rule{{
			Name: string(NumberMsgBase),
			Msg:  PickSchemaMsg(NumberMsgMap[NumberMsgBase], msg...),
			Fn: func(r Rule, path string, value any) (any, *ValidationError) {
				if value == nil {
					return value, nil
				}
				if _, ok := value.(float64); !ok {
					return value, &ValidationError{Path: path, Msg: r.Msg}
				}
				return value, nil
			},
		}},
	}
	return s
}
