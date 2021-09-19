data "zpa_saml_attribute" "email_user_sso" {
    name = "Email_User SSO"
}

output "get_zpa_saml_attribute" {
  value = data.zpa_saml_attribute.email_user_sso.id
}