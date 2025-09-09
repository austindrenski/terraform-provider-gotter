package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestGotterProvider(t *testing.T) {
	for name, test := range map[string]struct {
		check knownvalue.Check
		data  any
		text  string
	}{
		"empty_with_data": {
			check: knownvalue.StringExact(""),
			data:  `{ some_data = true }`,
			text:  ``,
		},
		"empty_with_empty": {
			check: knownvalue.StringExact(""),
			data:  `{}`,
			text:  ``,
		},
		"empty_with_null": {
			check: knownvalue.StringExact(""),
			data:  `null`,
			text:  ``,
		},
		"plaintext_with_data": {
			check: knownvalue.StringExact("test"),
			data:  `{ some_data = true }`,
			text:  `test`,
		},
		"plaintext_with_empty": {
			check: knownvalue.StringExact("test"),
			data:  `{}`,
			text:  `test`,
		},
		"plaintext_with_null": {
			check: knownvalue.StringExact("test"),
			data:  `null`,
			text:  `test`,
		},
		"identity_with_data": {
			check: knownvalue.StringExact("map[some_data:true]"),
			data:  `{ some_data = true }`,
			text:  `{{ print . }}`,
		},
		"identity_with_empty": {
			check: knownvalue.StringExact("map[]"),
			data:  `{}`,
			text:  `{{ print . }}`,
		},
		"identity_with_null": {
			check: knownvalue.StringExact("<nil>"),
			data:  `null`,
			text:  `{{ print . }}`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
					"gotter": providerserver.NewProtocol6WithError(New("dev")()),
				},
				Steps: []resource.TestStep{
					{
						Config: fmt.Sprintf(`output "test" { value = provider::gotter::execute(%q, %s) }`, test.text, test.data),
						ConfigStateChecks: []statecheck.StateCheck{
							statecheck.ExpectKnownOutputValue("test", test.check),
						},
					},
				},
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.SkipBelow(tfversion.Version1_8_0),
				},
			})
		})
	}
}
