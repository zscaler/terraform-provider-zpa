resource "zpa_user_portal_link" "this" {
  name        = "server1.example.com"
  description = "server1.example.com"
  enabled     = true
  link        = "server1.example.com"
  icon_text   = ""
  protocol    = "https://"
  user_portals {
    id = [zpa_user_portal_controller.this.id]
  }
}
