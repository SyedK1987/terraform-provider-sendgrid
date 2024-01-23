package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccapikeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_api_key" "name" {
					name = "test"
					permission = "bill" # "full" or "bill" or "custom"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "name", "test"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "id"),
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "permission"),
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "api_key"),
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "api_key_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendgrid_domain_authentication.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the sendgrid
				// API, therefore there is no value for it during import.
				//				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_api_key" "name" {
					name = "test"
					permission = "bill" # "full" or "bill" or "custom"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "name", "test"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "id"),
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "permission"),
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "api_key"),
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "api_key_id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
