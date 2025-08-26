package joi

import "fmt"

// --- schema ---

type ValidateOptions struct {
	Path *string
}

type Schema interface {
	Validate(value any) (any, []ValidationError)
	ValidateWithOpts(value any, opts ValidateOptions) (any, []ValidationError)
}

// --- validation ---

type ValidationError struct {
	Path string
	Msg  string
}

func (e ValidationError) String() string {
	return fmt.Sprintf("validation error at %q: %s", e.Path, e.Msg)
}

// --- rule ---

type Rule struct {
	Name string
	Args map[string]any
	Msg  string
	Fn   func(r Rule, path string, value any) (any, *ValidationError)
}
