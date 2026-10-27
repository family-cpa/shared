package directive

import (
	"context"
	"errors"
	"reflect"
	"slices"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/family-cpa/shared/pkg/gqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/vektah/gqlparser/v2/ast"
)

var engine = validator.New(validator.WithRequiredStructEnabled())

func Validate(ctx context.Context, obj any, next graphql.Resolver, tag string) (any, error) {
	val, err := next(ctx)
	if err != nil {
		return nil, err
	}
	val = trimString(val)
	inner, ok := concrete(val)
	if !ok || blank(inner) {
		if hasRequired(tag) {
			return nil, invalid(ctx, map[string]string{fieldName(ctx): "is required"})
		}
		return val, nil
	}
	if err := engine.Var(inner, tag); err != nil {
		return nil, invalid(ctx, map[string]string{fieldName(ctx): messageFor(err, inner)})
	}
	return val, nil
}

func invalid(ctx context.Context, fields map[string]string) error {
	ext := make(map[string]any, len(fields))
	for k, v := range fields {
		ext[k] = v
	}
	return gqlerr.Validation(ctx, ext)
}

func concrete(val any) (any, bool) {
	if val == nil {
		return nil, false
	}
	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return nil, false
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		return nil, false
	}
	return rv.Interface(), true
}

func blank(val any) bool {
	s, ok := val.(string)
	return ok && s == ""
}

func hasRequired(tag string) bool {
	for part := range strings.SplitSeq(tag, ",") {
		if strings.TrimSpace(part) == "required" {
			return true
		}
	}
	return false
}

func fieldName(ctx context.Context) string {
	path := graphql.GetPath(ctx)
	for _, p := range slices.Backward(path) {
		if name, ok := p.(ast.PathName); ok && name != "" {
			return string(name)
		}
	}
	return "value"
}

func messageFor(err error, val any) string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) || len(verrs) == 0 {
		return "is invalid"
	}
	e := verrs[0]
	_, isString := val.(string)
	switch e.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "max":
		if isString {
			return "must be at most " + e.Param() + " characters"
		}
		return "must be at most " + e.Param()
	case "min":
		if isString {
			return "must be at least " + e.Param() + " characters"
		}
		return "must be at least " + e.Param()
	default:
		return "is invalid"
	}
}

func trimString(val any) any {
	switch v := val.(type) {
	case string:
		return strings.TrimSpace(v)
	case *string:
		if v == nil {
			return v
		}
		s := strings.TrimSpace(*v)
		return &s
	default:
		return val
	}
}
