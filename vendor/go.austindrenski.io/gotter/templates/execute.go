package templates

import (
	"context"
	"fmt"
	"io"
	"text/template"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Execute renders a Go text template to the io.Writer
func Execute(ctx context.Context, t *template.Template, data any, w io.Writer) error {
	ctx, span := trace.SpanFromContext(ctx).TracerProvider().Tracer(scopeName).Start(ctx, "execute")
	defer span.End()

	span.SetAttributes(attribute.String("gotter.template.name", t.Name()))

	if err := t.Execute(w, data); err != nil {
		span.RecordError(err)
		span.SetAttributes(attribute.String("gotter.template.data", fmt.Sprintf("%v", data)))
		span.SetAttributes(attribute.String("gotter.template.text", fmt.Sprintf("%v", t)))
		span.SetStatus(codes.Error, "failed to execute template")
		return err
	}

	return nil
}
