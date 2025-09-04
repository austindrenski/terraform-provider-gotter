package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"go.austindrenski.io/gotter/utils"
	"go.austindrenski.io/terraform-provider-gotter/internal/provider"

	"go.opentelemetry.io/otel"
)

// scopeName is the instrumentation scope name.
const scopeName = "go.austindrenski.io/terraform-provider-gotter"

var version = "dev"

func main() {
	ctx := context.Background()

	end := utils.Start(ctx)
	defer end(ctx)

	ctx, span := otel.Tracer(scopeName).Start(ctx, "main")
	defer span.End()

	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/austindrenski/gotter",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), provider.New(version), opts); err != nil {
		span.RecordError(err)
		log.Fatal(err)
	}
}
