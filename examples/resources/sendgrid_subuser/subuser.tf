resource "sendgrid_subuser" "name" {
  email      = "yourname@example.com"
  username   = "name.test"
  ips = [
    "" # your domain ip. you can get this from sendgrid dashboard.
  ]
  password   = "yourpassword"
  disabled   = false # if you want to disable this subuser then change this value to true.
}