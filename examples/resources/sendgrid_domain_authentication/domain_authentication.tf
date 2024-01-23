resource "sendgrid_domain_authentication" "name" {
  
  domain = "example.com"
  environment = "nonprod" # prod or nonprod
  ips = [
    ""
  ]
  custom_spf = false
  default = false # if you want to make this domain as default domain then change this value to true.
}