terraform {
  required_providers {
    wpa = {
      source  = "sv-cheats-1/wpa"
      version = "0.0.1"
    }
  }
}

provider wpa {
  customer_id             = "C01rqf6p6" //"C02j7nz6z"
  impersonated_user_email = "admin@sellyoursoul.today"
  credentials             = "credentials.json"
}


resource "wpa_email_gateway" "default" {
  domain_name = "sellyoursoul.today"
}

#output "xxx" {
#   value = data.wpa_email_gateway.default.smart_host
#}
#
#output "yyy" {
#   value = data.wpa_email_gateway.default.smtp_mode
#}
#
#output "zzz" {
#  value = data.wpa_email_gateway.default
#}

