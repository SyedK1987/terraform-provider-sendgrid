resource "sendgrid_api_key" "name" {
  name = "test"
  scopes = [
    "mail.send",
    "alerts.create",
    "alerts.read"
  ]
}