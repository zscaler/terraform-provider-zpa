resource "zpa_user_portal_aup" "this" {
  name        = "Org_AUP01"
  description = "Org_AUP01"
  enabled     = true
  aup         = "Org_AUP01"
  email       = "company@acme.com"
  phone_num   = "+1 123-1458"
}
