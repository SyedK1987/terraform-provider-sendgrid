resource "sendgrid_single_sender" "name" {
  
  nickname = "tst"
  from_email = "yourname@example.com"
  from_name = "Your Name"
  reply_to = "replyto@example.com"
  reply_to_name = "Reply to name"
  address = "1234 Fake St"
  address2 = "Apt 123"
  city = "San Francisco"
  state = "CA"
  zip = "95369"
  country = "US"
}