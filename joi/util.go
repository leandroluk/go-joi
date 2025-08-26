package joi

import (
	"fmt"
	"maps"
	"reflect"
	"strings"
	"time"
)

func RunValidation(rules []Rule, label, path string, value any) (any, []ValidationError) {
	var errs []ValidationError
	current := value

	for _, r := range rules {
		newVal, err := r.Fn(r, path, current)
		if err != nil {
			msg := Coalesce(r.Msg, err.Msg)
			ctx := map[string]any{"label": label, "path": path, "value": current}
			maps.Copy(ctx, r.Args)
			err.Msg = RenderTemplate(msg, ctx)
			errs = append(errs, *err)
		}
		if newVal != nil {
			current = newVal
		}
	}

	return current, errs
}

func RenderTemplate(template string, context map[string]any) string {
	out := template
	for key, value := range context {
		placeholder := "{{#" + key + "}}"
		out = strings.ReplaceAll(out, placeholder, fmt.Sprintf("%v", value))
	}
	return out
}

func Coalesce[T comparable](vals ...T) T {
	var zero T
	for _, v := range vals {
		if v != zero {
			return v
		}
	}
	return zero
}

func WhenNil(v any, err *ValidationError) *ValidationError {
	if v == nil {
		return err
	}
	return nil
}

func PickSchemaMsg(defaultMsg string, override ...string) string {
	if len(override) > 0 && override[0] != "" {
		return override[0]
	}
	return defaultMsg
}

func ValueInList(value any, list []any) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

func ParseDate(value any) (time.Time, bool) {
	switch v := value.(type) {
	case time.Time:
		return v, true
	case string:
		// try RFC3339 parse
		if t, err := time.Parse(PARSE_LAYOUT, v); err == nil {
			return t, true
		}
		return time.Time{}, false
	case int64:
		return time.Unix(v, 0), true
	case float64:
		// JSON numbers caem como float64
		return time.Unix(int64(v), 0), true
	default:
		return time.Time{}, false
	}
}

func Ptr[T any](v T) *T {
	return &v
}
