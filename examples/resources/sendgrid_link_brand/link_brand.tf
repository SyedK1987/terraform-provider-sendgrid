resource "sendgrid_linkbrand" "example" {
  domain = "example.com"
  subdomain = "url09em21"
  default = false # if you want to make this domain as default domain then change this value to true.
}