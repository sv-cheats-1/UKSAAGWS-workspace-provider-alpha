terraform {
  required_providers {
    wpa = {
      source  = "sv-cheats-1/wpa"
      version = "0.0.1"
    }
  }
}

provider wpa {
  customer_id             = "yourcustomerID"
  impersonated_user_email = "admin@yourdomain.today"
  credentials             = "credentials.json"
}


data "wpa_email_gateway" "default" {
  domain_name = "yourdomain.today"
}

output "smart_host" {
  value = data.wpa_email_gateway.default.smart_host
}

output "smtp_mode" {
  value = data.wpa_email_gateway.default.smtp_mode
}



