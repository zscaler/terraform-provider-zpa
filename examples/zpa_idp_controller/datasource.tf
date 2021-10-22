data "zpa_idp_controller" "example" {
 name = "IdP-User-Name"
}

output "idp_controller" {
    value = data.zpa_idp_controller.example
}