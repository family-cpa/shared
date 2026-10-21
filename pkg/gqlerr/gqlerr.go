package gqlerr

import (
	"context"
	"errors"
	"fmt"
	"maps"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type ErrorType string

func (t ErrorType) String() string {
	return string(t)
}

const (
	ErrorUnknown      ErrorType = ""
	ErrorInternal     ErrorType = "INTERNAL_SERVER_ERROR"
	ErrorForbidden    ErrorType = "FORBIDDEN"
	ErrorUnauthorized ErrorType = "UNAUTHENTICATED"
	ErrorNotFound     ErrorType = "NOT_FOUND"
	ErrorUnavailable  ErrorType = "UNAVAILABLE"
	ErrorConflict     ErrorType = "CONFLICT"
	ErrorValidation   ErrorType = "GRAPHQL_VALIDATION_FAILED"
)

type Options struct {
	Type       ErrorType
	Message    string
	Extensions map[string]any
}

func New(ctx context.Context, opt Options) *gqlerror.Error {
	err := &gqlerror.Error{
		Extensions: make(map[string]any),
	}
	if opt.Extensions != nil && len(opt.Extensions) > 0 {
		maps.Copy(err.Extensions, opt.Extensions)
	}
	if path := graphql.GetPath(ctx); len(path) > 0 {
		err.Path = path
	}
	if len(opt.Type) > 0 {
		err.Extensions["code"] = opt.Type.String()
		switch opt.Type {
		case ErrorInternal:
			err.Extensions["message"] = "Internal server error"
		case ErrorForbidden:
			err.Extensions["message"] = "You are not authorized to perform this action"
		case ErrorUnauthorized:
			err.Extensions["message"] = "You are not authenticated"
		case ErrorNotFound:
			err.Extensions["message"] = "Object does not exist"
		case ErrorUnavailable:
			err.Extensions["message"] = "This action is temporarily unavailable"
		case ErrorConflict:
			err.Extensions["message"] = "Conflict of action interdependencies"
		case ErrorValidation:
			err.Extensions["message"] = "Validation failed, the given input was invalid"
		case ErrorUnknown:
			err.Extensions["message"] = "Unknown error occurred"
		}
	}
	if len(opt.Message) > 0 {
		err.Message = opt.Message
	}
	return err
}

func Is(err error, t ErrorType) bool {
	if err == nil {
		return false
	}
	var gqlErr *gqlerror.Error
	if !errors.As(err, &gqlErr) {
		return false
	}
	if gqlErr.Extensions == nil {
		return false
	}
	code, ok := gqlErr.Extensions["code"].(string)
	if !ok {
		return false
	}
	return code == t.String()
}

func Internal(ctx context.Context) *gqlerror.Error {
	return New(ctx, Options{Type: ErrorInternal})
}

func Forbidden(ctx context.Context) *gqlerror.Error {
	return New(ctx, Options{Type: ErrorForbidden})
}

func Unauthorized(ctx context.Context) *gqlerror.Error {
	return New(ctx, Options{Type: ErrorUnauthorized})
}

func NotFound(ctx context.Context, id ...any) *gqlerror.Error {
	msg := ""
	if id != nil || len(id) > 0 {
		msg = fmt.Sprintf("Object with ID %v does not exist", id)
	}
	return New(ctx, Options{Type: ErrorNotFound, Message: msg})
}

func Unavailable(ctx context.Context) *gqlerror.Error {
	return New(ctx, Options{Type: ErrorUnavailable})
}

func Conflict(ctx context.Context) *gqlerror.Error {
	return New(ctx, Options{Type: ErrorConflict})
}

func Validation(ctx context.Context, fields map[string]any) *gqlerror.Error {
	return New(ctx, Options{Type: ErrorValidation, Extensions: map[string]any{"fields": fields}})
}
