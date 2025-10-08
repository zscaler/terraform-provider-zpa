## Example Usage - With Customer Own Certificate


data "zpa_ba_certificate" "this" {
  name = "example.acme.com"
}

resource "zpa_user_portal_controller" "this" {
  name                      = "UserPortal01"
  description               = "UserPortal01"
  enabled                   = true
  user_notification         = "User_Portal_Terraform_01"
  user_notification_enabled = true
  certificate_id            = data.zpa_ba_certificate.this.id
  domain                    = "portal01"
}


## Example Usage - With Zscaler Managed Certificate

resource "zpa_user_portal_controller" "this" {
  name                      = "UserPortal01"
  description               = "UserPortal01"
  enabled                   = true
  user_notification         = "UserPortal01"
  user_notification_enabled = true
  certificate_id            = ""
  ext_domain_translation    = "acme.io"
  ext_label                 = "portal01"
  ext_domain_name           = "acme-io.b.zscalerportal.net"
  ext_domain                = "acme.io"
  domain                    = "portal01-acme-io.b.zscalerportal.net"
}
