data "zpa_idp_controller" "example" {
 name = "IDP-Name"
}

output "idp_controller" {
    value = data.zpa_idp_controller.example
}