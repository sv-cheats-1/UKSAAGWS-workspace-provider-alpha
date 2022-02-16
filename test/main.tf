provider wpa {
  customer_id = "" # optional
}

data "wpa_mail_routing_settings" "default" {
   customer_id = "38xck8d5j"
}

output "xxx" {
   value = data.wpa_mail_routing_settings.default.smart_host
}

output "yyy" {
   value = data.wpa_mail_routing_settings.default.smtp_mode
}
