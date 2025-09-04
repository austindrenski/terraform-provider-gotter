package templates

import (
	"context"
	"text/template"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Parse(ctx context.Context, name string, text string, opts ...func(context.Context, *template.Template) *template.Template) (*template.Template, error) {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "parse")
	defer span.End()

	span.SetAttributes(
		attribute.String("gotter.template.name", name),
		attribute.String("gotter.template.text", text))

	t := template.New(name)

	for _, opt := range opts {
		t = opt(ctx, t)
	}

	if t, err := t.Parse(text); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse template")
		return nil, err
	} else {
		return t, err
	}
}

func WithFuncs(f func(ctx context.Context) template.FuncMap) func(context.Context, *template.Template) *template.Template {
	return func(ctx context.Context, t *template.Template) *template.Template {
		return t.Funcs(f(ctx))
	}
}
