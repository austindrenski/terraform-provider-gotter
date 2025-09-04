package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"

	"go.austindrenski.io/gotter/templates"
)

var _ function.Function = execute{}

type execute struct{}

func (f execute) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Description: "Executes a Go text/template from `file` using the provided `data`",
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue: false,
				Description:    "The text template",
				Name:           "text",
				Validators: []function.StringParameterValidator{
					execute{},
				},
			},
			function.DynamicParameter{
				AllowNullValue: true,
				Description:    "The data object passed to the template",
				Name:           "data",
			},
		},
		Return:  function.StringReturn{},
		Summary: "Executes a Go text/template from `file` using the provided `data`",
	}
}

func (f execute) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "execute"
}

func (f execute) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var text string
	var data any

	if err := req.Arguments.Get(ctx, &text, &data); err != nil {
		resp.Error = function.ConcatFuncErrors(err)
		return
	}

	t, err := templates.Parse(ctx, "", text, templates.WithFuncs(templates.Functions))
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	b := strings.Builder{}
	if err := templates.Execute(ctx, t, data, &b); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	if err := resp.Result.Set(ctx, b.String()); err != nil {
		resp.Error = err
		return
	}
}

func (f execute) ValidateParameterString(ctx context.Context, req function.StringParameterValidatorRequest, resp *function.StringParameterValidatorResponse) {
	if _, err := templates.Parse(ctx, "", req.Value.ValueString(), templates.WithFuncs(templates.Functions)); err != nil {
		resp.Error = function.NewArgumentFuncError(req.ArgumentPosition, err.Error())
	}
}
