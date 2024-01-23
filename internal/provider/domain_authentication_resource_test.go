package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccdomainauthResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_domain_authentication" "name" {
  
					domain = "example.com"
					custom_dkim = "tss" # only 3 characters are allowed. it can be any 3 characters.
					ips = [
					  ""
					]
					custom_spf = false
					default = false # if you want to make this domain as default domain then change this value to true.
					auto_security = false
					valid = false # once domain created updated the dns/route 53 records and then run terraform apply again by changing this value true.
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "subdomain", ""),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "custom_spf", "false"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "default", "false"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "auto_Security", "false"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "valid", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "id"),
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
				resource "sendgrid_domain_authentication" "name" {
  
					domain = "example.com"
					custom_dkim = "tss" # only 3 characters are allowed. it can be any 3 characters.
					ips = [
					  ""
					]
					custom_spf = false
					default = false # if you want to make this domain as default domain then change this value to true.
					auto_security = false
					valid = false # once domain created updated the dns/route 53 records and then run terraform apply again by changing this value true.
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "subdomain", ""),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "custom_spf", "false"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "default", "false"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "auto_Security", "false"),
					resource.TestCheckResourceAttr("sendgrid_domain_authentication.test", "valid", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_domain_authentication.test", "id"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
