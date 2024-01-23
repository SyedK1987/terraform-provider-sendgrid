package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccdomainauthsubuserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_domainauth_add_subuser" "asub" {
					id = sendgrid_domain_authentication.name.id
					username = "sk.dev"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_domainauth_add_subuser.asub", "domain", "example.com"),
					resource.TestCheckResourceAttr("sendgrid_domainauth_add_subuser.asub", "username", "sk.dev"),
					resource.TestCheckResourceAttrSet("sendgrid_domainauth_add_subuser.asub", "subusers"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendgrid_domainauth_add_subuser.asub",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the HashiCups
				// API, therefore there is no value for it during import.
				//				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_domainauth_add_subuser" "asub" {
					id = sendgrid_domain_authentication.name.id
					username = "sk.dev"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_domainauth_add_subuser.asub", "domain", "example.com"),
					resource.TestCheckResourceAttr("sendgrid_domainauth_add_subuser.asub", "username", "sk.dev"),
					resource.TestCheckResourceAttrSet("sendgrid_domainauth_add_subuser.asub", "subusers"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
