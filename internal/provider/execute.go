package provider

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"go.austindrenski.io/gotter/templates"
)

var (
	_ function.Function                 = (*execute)(nil)
	_ function.StringParameterValidator = (*execute)(nil)
)

type execute struct {
	file bool
	name string
}

func (f execute) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	templateParameter := function.StringParameter{
		AllowNullValue: false,
		Validators: []function.StringParameterValidator{
			execute{},
		},
	}

	if f.file {
		templateParameter.Description = "The text template file"
		templateParameter.Name = "file"
	} else {
		templateParameter.Description = "The text template"
		templateParameter.Name = "text"
	}

	resp.Definition = function.Definition{
		Description: fmt.Sprintf("Executes a Go text/template from `%s` using the provided `data`", templateParameter.GetName()),
		Parameters: []function.Parameter{
			templateParameter,
			function.DynamicParameter{
				AllowNullValue: true,
				Description:    "The data passed to the template",
				Name:           "data",
			},
		},
		Return:  function.StringReturn{},
		Summary: fmt.Sprintf("Executes a Go text/template from `%s` using the provided `data`", templateParameter.GetName()),
	}
}

func (f execute) Metadata(_ context.Context, _ function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = f.name
}

func (f execute) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var text string
	var data types.Dynamic

	if err := req.Arguments.Get(ctx, &text, &data); err != nil {
		resp.Error = function.ConcatFuncErrors(err)
		return
	}

	t, err := f.parse(ctx, text)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	b := strings.Builder{}
	if err := templates.Execute(ctx, t, unwrap(data), &b); err != nil {
		resp.Error = function.NewFuncError(err.Error())
		return
	}

	if err := resp.Result.Set(ctx, b.String()); err != nil {
		resp.Error = err
		return
	}
}

func (f execute) ValidateParameterString(ctx context.Context, req function.StringParameterValidatorRequest, resp *function.StringParameterValidatorResponse) {
	v := req.Value.ValueString()

	if f.file {
		if stat, err := os.Stat(v); err != nil {
			resp.Error = function.NewArgumentFuncError(req.ArgumentPosition, err.Error())
			return
		} else if stat.IsDir() {
			resp.Error = function.NewArgumentFuncError(req.ArgumentPosition, fmt.Sprintf("%q is a directory", v))
			return
		} else if stat.Size() == 0 {
			resp.Error = function.NewArgumentFuncError(req.ArgumentPosition, fmt.Sprintf("%q is empty", v))
			return
		} else if _, err := templates.ParseFile(ctx, v, templates.WithFuncs(templates.Functions)); err != nil {
			resp.Error = function.NewArgumentFuncError(req.ArgumentPosition, err.Error())
			return
		}
	} else if _, err := templates.Parse(ctx, "", v, templates.WithFuncs(templates.Functions)); err != nil {
		resp.Error = function.NewArgumentFuncError(req.ArgumentPosition, err.Error())
	}
}

func (f execute) parse(ctx context.Context, text string) (*template.Template, error) {
	if f.file {
		return templates.ParseFile(ctx, text, templates.WithFuncs(templates.Functions))
	} else {
		return templates.Parse(ctx, "", text, templates.WithFuncs(templates.Functions))
	}
}

func unwrap(d types.Dynamic) any {
	switch d := d.UnderlyingValue().(type) {
	case types.Bool:
		return d.ValueBool()
	case types.Number:
		return d.ValueBigFloat()
	case types.List:
		return d.Elements()
	case types.Map:
		return d.Elements()
	case types.Object:
		return d.Attributes()
	case types.Set:
		return d.Elements()
	case types.String:
		return d.ValueString()
	case types.Tuple:
		return d.Elements()
	default:
		return nil
	}
}
