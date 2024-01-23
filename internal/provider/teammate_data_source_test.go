package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccteammateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				data "sendgrid_teammate" "test" {
					email = "yourname@example.com"
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "username", "sk.test"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendgrid_teammate.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the sendgrid
				// API, therefore there is no value for it during import.
				//				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				data "sendgrid_teammate" "test" {
					email = "yourname@example.com"
					is_admin = false
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "username", "sk.test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
