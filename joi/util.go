package joi

import (
	"fmt"
	"maps"
	"reflect"
	"strings"
)

func runValidation(rules []Rule, label, path string, value any) (any, []ValidationError) {
	var errs []ValidationError
	current := value

	for _, r := range rules {
		newVal, err := r.Fn(r, path, current)
		if err != nil {
			msg := coalesce(r.Msg, err.Msg)
			ctx := map[string]any{"label": label, "path": path, "value": current}
			maps.Copy(ctx, r.Args)
			err.Msg = renderTemplate(msg, ctx)
			errs = append(errs, *err)
		}
		if newVal != nil {
			current = newVal
		}
	}

	return current, errs
}

func renderTemplate(template string, context map[string]any) string {
	out := template
	for key, value := range context {
		placeholder := "{{#" + key + "}}"
		out = strings.ReplaceAll(out, placeholder, fmt.Sprintf("%v", value))
	}
	return out
}

func coalesce[T comparable](vals ...T) T {
	var zero T
	for _, v := range vals {
		if v != zero {
			return v
		}
	}
	return zero
}

func whenNil(v any, err *ValidationError) *ValidationError {
	if v == nil {
		return err
	}
	return nil
}

func pickMsg(defaultMsg string, override ...string) string {
	if len(override) > 0 && override[0] != "" {
		return override[0]
	}
	return defaultMsg
}

func valueInList(value any, list []any) bool {
	for _, v := range list {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

func newSchema[T any](base Rule) *AnySchema[T] {
	s := &AnySchema[T]{label: "value", rules: []Rule{base}}
	s.self = any(s).(T)
	return s
}
