resource "zpa_microtenant_controller" "this" {
  name = "Microtenant_A"
  description = "Microtenant_A"
  enabled = true
  criteria_attribute = "AuthDomain"
  criteria_attribute_values = ["acme.com"]
}

// To output specific Microtenant user information,
// the following output configuration is required.
output "zpa_microtenant_controller1" {
  value = [for u in zpa_microtenant_controller.this.user : {
    microtenant_id = u.microtenant_id
    username       = u.username
    password       = u.password
  }]
}