package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.ProviderWithFunctions = (*gotterProvider)(nil)

type gotterProvider struct {
	name    string
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return gotterProvider{
			name:    "gotter",
			version: version,
		}
	}
}

func (p gotterProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
}

func (p gotterProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

func (p gotterProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		func() function.Function {
			return execute{
				file: false,
				name: "execute",
			}
		},
		func() function.Function {
			return execute{
				file: true,
				name: "execute_file",
			}
		},
	}
}

func (p gotterProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = p.name
	resp.Version = p.version
}

func (p gotterProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p gotterProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A provider for Go text/template processing.",
	}
}
