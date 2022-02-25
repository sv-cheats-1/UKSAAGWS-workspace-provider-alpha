terraform {
  required_providers {
    wpa = {
      source  = "sv-cheats-1/wpa"
      version = "0.0.1"
    }
  }
}

provider wpa {
  customer_id             = "C02x7xx6x"
  impersonated_user_email = "your.email@email.com"
  credentials             = "credentials.json"
}

data "wpa_email_gateway" "default" {
  domain_name = "your.domain"
}

output "xxx" {
   value = data.wpa_email_gateway.default.smart_host
}

output "yyy" {
   value = data.wpa_email_gateway.default.smtp_mode
}

output "zzz" {
  value = data.wpa_email_gateway.default
}
