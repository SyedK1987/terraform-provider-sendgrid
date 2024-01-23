resource "sendgrid_domainauth_add_subuser" "asub" {
  id = sendgrid_domain_authentication.name.id
  username = "subusername"
}