data "zpa_scim_attribute_header" "email_value" {
    name = "emails.value"
    idp_name = "idp_name"
}

output "get_zpa_scim_attribute_header" {
  value = data.zpa_scim_attribute_header.email_value
}