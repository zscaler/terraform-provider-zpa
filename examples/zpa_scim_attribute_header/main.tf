data "zpa_scim_attribute_header" "email_value" {
    name = "emails.value"
    idp_name = "IdP_Name"
}

data "zpa_scim_attribute_header" "givenName" {
  name     = "name.givenName"
  idp_name = "IdP_Name"
}

data "zpa_scim_attribute_header" "familyName" {
  name     = "name.familyName"
  idp_name = "IdP_Name"
}