variable "ips" {
  type =  list(string)
  default = [
    "xxx.xx.xxx.xxx/32",
    "xxx.xx.xxx.0/22"
  ]
}

# multiple IP can be managed using for_each. Please make sure your ip always ends with /32 if it's single ip or you can use complete cidr.
resource "sendgrid_ipwhitelist" "name" {
  for_each = toset(var.ips)
  ip = each.key
}