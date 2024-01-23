package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcclinkvalidateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_linkbrand_validate" "name" {
					id = sendgrid_domain_authentication.name.id
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendgrid_linkbrand_validate.test", "id"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand_validate.test", "valid", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendgrid_linkbrand_validate.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the sendgrid
				// API, therefore there is no value for it during import.
				//				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_linkbrand_validate" "name" {
					id = sendgrid_domain_authentication.name.id
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_linkbrand_validate.test", "valid", "true"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_linkbrand_validate.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
