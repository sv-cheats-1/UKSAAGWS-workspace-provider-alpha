terraform {
  required_providers {
    wpa = {
      source  = "sv-cheats-1/wpa"
      version = "0.0.1"
    }
  }
}

provider wpa {
  # TODO implement optional provider configuration (check hashicups)
  #domain_name = "" # optional
}

data "wpa_email_gateway" "default" {
  domain_name = "38xck8d5j"
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
