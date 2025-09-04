package templates

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func Functions(ctx context.Context) template.FuncMap {
	return template.FuncMap{
		"json": func(source any) (string, error) {
			return json_(ctx, source)
		},
		"lower": func(source string) string {
			return lower(ctx, source)
		},
		"match": func(pattern string, s string) (bool, error) {
			return match(ctx, pattern, s)
		},
		"replace": func(pattern string, replacement string, source string) (string, error) {
			return replace(ctx, pattern, replacement, source)
		},
		"split": func(pattern string, source string) ([]string, error) {
			return split(ctx, pattern, source)
		},
		"split_n": func(pattern string, n int, source string) ([]string, error) {
			return split_n(ctx, pattern, n, source)
		},
		"title": func(source string) string {
			return title(ctx, source)
		},
		"truncate": func(n int, source string) string {
			return truncate(ctx, n, source)
		},
		"upper": func(source string) string {
			return upper(ctx, source)
		},
	}
}

func json_(ctx context.Context, data any) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "json")
	defer span.End()

	if m, err := json.Marshal(data); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("gotter.template.func.json.data", fmt.Sprintf("%v", data)))
		return "", err
	} else {
		return string(m), nil
	}
}

func lower(ctx context.Context, source string) string {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "lower")
	defer span.End()

	return strings.ToLower(source)
}

func match(ctx context.Context, pattern string, source string) (bool, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "match")
	defer span.End()

	if m, err := regexp.MatchString(pattern, source); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("gotter.template.func.match.pattern", pattern))
		span.SetAttributes(attribute.String("gotter.template.func.match.source", source))
		return false, err
	} else {
		return m, nil
	}
}

func replace(ctx context.Context, pattern string, replacement string, source string) (string, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "replace")
	defer span.End()

	if r, err := regexp.Compile(pattern); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("gotter.template.func.replace.pattern", pattern))
		span.SetAttributes(attribute.String("gotter.template.func.replace.replacement", replacement))
		span.SetAttributes(attribute.String("gotter.template.func.replace.source", source))
		return "", err
	} else {
		return r.ReplaceAllString(source, replacement), nil
	}
}

func split(ctx context.Context, pattern string, source string) ([]string, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "split")
	defer span.End()

	if r, err := regexp.Compile(pattern); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("gotter.template.func.split.pattern", pattern))
		span.SetAttributes(attribute.String("gotter.template.func.split.source", source))
		return nil, err
	} else {
		return r.Split(source, -1), nil
	}
}

//goland:noinspection GoSnakeCaseUsage
func split_n(ctx context.Context, pattern string, n int, source string) ([]string, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "split_n")
	defer span.End()

	if r, err := regexp.Compile(pattern); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.Int("gotter.template.func.split_n.n", n))
		span.SetAttributes(attribute.String("gotter.template.func.split_n.pattern", pattern))
		span.SetAttributes(attribute.String("gotter.template.func.split_n.source", source))
		return nil, err
	} else {
		return r.Split(source, n), nil
	}
}

func title(ctx context.Context, source string) string {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "title")
	defer span.End()

	return cases.Title(language.AmericanEnglish).String(source)
}

func truncate(ctx context.Context, n int, source string) string {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "truncate")
	defer span.End()

	if n < len(source) {
		return source[:n]
	} else {
		return source
	}
}

func upper(ctx context.Context, source string) string {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "upper")
	defer span.End()

	return strings.ToUpper(source)
}
