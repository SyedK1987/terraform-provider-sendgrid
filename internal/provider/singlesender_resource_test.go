package sendgrid

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccsinglesenderResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "sendgrid_single_sender" "test" {
	nickname = "sk"
	from_email = "sikaleel87@gmail.com"
	from_name = "Syed Kaleel"
	reply_to = "sikaleel87@gmail.com"
	reply_to_name = "Syed Kaleel"
	address = "1234 Fake St"
	address2 = "Apt 123"
	city = "San Francisco"
	state = "CA"
	zip = "95369"
	country = "US"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of items
					// Verify first order item
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "nickname", "sk"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "from_email", "sikaleel87@gmail.com"),
					// Verify first coffee item has Computed attributes filled.
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "from_name", "Syed Kaleel"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "reply_to", "sikaleel87@gmail.com"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "reply_to_name", "Syed Kaleel"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "address", "1234 Fake St"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "address2", "Apt 123"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "city", "San Francisco"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "state", "CA"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "zip", "95639"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "country", "US"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "verified", "false"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "locked", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("sendgrid_single_sender.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "sendgrid_single_sender.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The last_updated attribute does not exist in the sendgrid
				// API, therefore there is no value for it during import.
				//				ImportStateVerifyIgnore: []string{"last_updated"},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "sendgrid_single_sender" "test" {
	nickname = "sk"
	from_email = "sikaleel87@gmail.com"
	from_name = "Syed Kaleel"
	reply_to = "sikaleel87@gmail.com"
	reply_to_name = "Syed Kaleel"
	address = "1234 Fake St"
	address2 = "Apt 123"
	city = "San Francisco"
	state = "CA"
	zip = "95369"
	country = "US"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify first order item updated
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "nickname", "sk"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "from_email", "sikaleel87@gmail.com"),
					// Verify first coffee item has Computed attributes filled.
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "from_name", "Syed Kaleel"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "reply_to", "sikaleel87@gmail.com"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "reply_to_name", "Syed Kaleel"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "address", "1234 Fake St"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "address2", "Apt 123"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "city", "San Francisco"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "state", "CA"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "zip", "95639"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "country", "US"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "verified", "false"),
					resource.TestCheckResourceAttr("sendgrid_single_sender.test", "locked", "false"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
