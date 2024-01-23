package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcclinkbrandingDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
				data "sendgrid_linkbrand" "test" {
					id = 123456789
				  }
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendgrid_linkbrand.test", "id"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "subdomain", "url091234"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "default", "false"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "valid", "true"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "legacy", "false"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "user_id", "1234567"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "username", "testuser"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "domain_cname", "{\"data\":\"sendgrid.net\",\"valid\":true,\"host\":\"url091234.example.com\",\"types\":\"cname\"}"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "owner_cname", "{\"data\":\"sendgrid.net\",\"valid\":true,\"host\":\"url091234.example.com\",\"types\":\"cname\"}"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendgrid_linkbrand.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the sendgrid
				// API, therefore there is no value for it during import.
				//				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "sendgrid_link_branding" "name" {
					domain = "example.com"
  					subdomain = "url09em21"
  					default = false
				  }
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("sendgrid_linkbrand.test", "id"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "domain", "example.com"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "subdomain", "url091234"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "default", "false"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "valid", "true"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "legacy", "false"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "user_id", "1234567"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "username", "testuser"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "domain_cname", "{\"data\":\"sendgrid.net\",\"valid\":true,\"host\":\"url091234.example.com\",\"types\":\"cname\"}"),
					resource.TestCheckResourceAttr("sendgrid_linkbrand.test", "owner_cname", "{\"data\":\"sendgrid.net\",\"valid\":true,\"host\":\"url091234.example.com\",\"types\":\"cname\"}"),
				),
			},
		},
	})
}
