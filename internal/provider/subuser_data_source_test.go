package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccsubuserDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				data "sendgrid_subuser" "test" {
					username   = "chid.test"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "username", "sk.test"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "disabled", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_subuser.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendgrid_subuser.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the sendgrid
				// API, therefore there is no value for it during import.
				//				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				data "sendgrid_subuser" "test" {
					username   = "chid.test"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "username", "sk.test"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "disabled", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_subuser.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
