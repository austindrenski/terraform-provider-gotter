// Package utils defines OpenTelemetry helpers.
package utils // import "go.austindrenski.io/gotter/utils"

import (
	"context"
	"log"
	"sync"

	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv/v1.37.0"
)

//goland:noinspection GoNameStartsWithPackageName,GoSnakeCaseUsage
var (
	OTEL_VCS_CHANGE_ID           = ""
	OTEL_VCS_OWNER_NAME          = ""
	OTEL_VCS_REF_BASE_NAME       = ""
	OTEL_VCS_REF_BASE_REVISION   = ""
	OTEL_VCS_REF_BASE_TYPE       = ""
	OTEL_VCS_REF_HEAD_NAME       = ""
	OTEL_VCS_REF_HEAD_REVISION   = ""
	OTEL_VCS_REF_HEAD_TYPE       = ""
	OTEL_VCS_REPOSITORY_NAME     = ""
	OTEL_VCS_REPOSITORY_URL_FULL = ""
)

var (
	withFallbackLogExporter = autoexport.WithFallbackLogExporter(func(ctx context.Context) (sdklog.Exporter, error) {
		return nil, nil
	})
	withFallbackMetricReader = autoexport.WithFallbackMetricReader(func(ctx context.Context) (metric.Reader, error) {
		return nil, nil
	})
	withFallbackSpanExporter = autoexport.WithFallbackSpanExporter(func(ctx context.Context) (trace.SpanExporter, error) {
		return nil, nil
	})
)

func Start(ctx context.Context) func(context.Context) {
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	res := resource.Default()
	if r, err := resource.Merge(res, resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.VCSChangeID(OTEL_VCS_CHANGE_ID),
		semconv.VCSOwnerName(OTEL_VCS_OWNER_NAME),
		semconv.VCSRefBaseName(OTEL_VCS_REF_BASE_NAME),
		semconv.VCSRefBaseRevision(OTEL_VCS_REF_BASE_REVISION),
		semconv.VCSRefBaseTypeKey.String(OTEL_VCS_REF_BASE_TYPE),
		semconv.VCSRefHeadName(OTEL_VCS_REF_HEAD_NAME),
		semconv.VCSRefHeadRevision(OTEL_VCS_REF_HEAD_REVISION),
		semconv.VCSRefHeadTypeKey.String(OTEL_VCS_REF_HEAD_TYPE),
		semconv.VCSRepositoryName(OTEL_VCS_REPOSITORY_NAME),
		semconv.VCSRepositoryURLFull(OTEL_VCS_REPOSITORY_URL_FULL))); err == nil {
		res = r
	}

	if exp, err := autoexport.NewLogExporter(ctx, withFallbackLogExporter); err != nil {
		log.Fatal(err)
	} else {
		global.SetLoggerProvider(sdklog.NewLoggerProvider(sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)), sdklog.WithResource(res)))
	}

	if exp, err := autoexport.NewMetricReader(ctx, withFallbackMetricReader); err != nil {
		log.Fatal(err)
	} else {
		otel.SetMeterProvider(metric.NewMeterProvider(metric.WithReader(exp), metric.WithResource(res)))
	}

	if exp, err := autoexport.NewSpanExporter(ctx, withFallbackSpanExporter); err != nil {
		log.Fatal(err)
	} else {
		otel.SetTracerProvider(trace.NewTracerProvider(trace.WithBatcher(exp), trace.WithResource(res)))
	}

	return func(ctx context.Context) {
		wg := sync.WaitGroup{}

		wg.Go(func() {
			if p, ok := global.GetLoggerProvider().(*sdklog.LoggerProvider); ok {
				if err := p.Shutdown(ctx); err != nil {
					log.Print(err)
				}
			}
		})

		wg.Go(func() {
			if p, ok := otel.GetMeterProvider().(*metric.MeterProvider); ok {
				if err := p.Shutdown(ctx); err != nil {
					log.Print(err)
				}
			}
		})

		wg.Go(func() {
			if p, ok := otel.GetTracerProvider().(*trace.TracerProvider); ok {
				if err := p.Shutdown(ctx); err != nil {
					log.Print(err)
				}
			}
		})

		wg.Wait()
	}
}
