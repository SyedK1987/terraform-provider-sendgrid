package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccsubuserResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_subuser" "test" {
					email      = "sk@example.com"
					username   = "sk.test"
					ips = [
					  "" # your domain ip. you can get this from sendgrid dashboard.
					]
					password   = "C3|zh!%SR],jgD5d"
					disabled   = false # if you want to disable this subuser then change this value to true.
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "username", "sk.test"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "password", "C3|zh!%SR],jgD5d"),
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
				resource "sendgrid_subuser" "test" {
					email      = "sk@example.com"
					username   = "chid.test"
					ips = [
					  "" # your domain ip. you can get this from sendgrid dashboard.
					]
					password   = "C3|zh!%SR],jgD5d"
					disabled   = false # if you want to disable this subuser then change this value to true.
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "email", "sk@example.com"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "username", "sk.test"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "password", "C3|zh!%SR],jgD5d"),
					resource.TestCheckResourceAttr("sendgrid_subuser.test", "disabled", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_subuser.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
