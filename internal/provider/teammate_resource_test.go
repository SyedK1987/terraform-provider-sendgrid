package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccteammateResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_teammate" "name" {
					email = "yourname@example.com"
					is_admin = false
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "username", "sk.test"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "first_name", "sk"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "last_name", "Test"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "is_admin", "false"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "is_read_only", "false"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "expiration_date", "none"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "user_type", "teammate"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_teammate.test", "token"),
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
				resource "sendgrid_teammate" "name" {
					email = "yourname@example.com"
					is_admin = false
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "username", "sk.test"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "first_name", "sk"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "last_name", "Test"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "is_admin", "false"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "is_read_only", "false"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "expiration_date", "none"),
					resource.TestCheckResourceAttr("sendgrid_teammate.test", "user_type", "teammate"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_teammate.test", "token"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
