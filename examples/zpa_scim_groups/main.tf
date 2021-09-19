data "zpa_scim_groups" "engineering" {
    name = "Engineering"
    idp_name = "idp_name"
}

output "get_zpa_scim_groups" {
  value = data.zpa_scim_groups.engineering
}