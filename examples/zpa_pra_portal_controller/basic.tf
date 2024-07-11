data "zpa_ba_certificate" "this" {
 name = "portal.acme.com"
}

resource "zpa_pra_portal_controller" "this" {
  name = "portal.acme.com"
  description = "portal.acme.com"
  enabled = true
  domain = "portal.acme.com"
  certificate_id = data.zpa_ba_certificate.this.id
  user_notification = "Created with Terraform"
  user_notification_enabled = true
}