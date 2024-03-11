terraform {
  required_providers {
    zpa = {
      version = "3.2.0"
      source  = "zscaler.com/zpa/zpa"
    }
  }
}

provider "zpa" {}

# data "zpa_saml_attribute" "email_user_sso" {
#     name = "Email_BD_Okta_Users"
#     idp_name = "BD_Okta_Users"
# }

# output "get_zpa_saml_attribute" {
#   value = data.zpa_saml_attribute.email_user_sso
# }

data "zpa_scim_groups" "a000" {
    name = "A000"
    idp_name = "BD_Okta_Users"
}

output "zpa_scim_groups" {
  value = data.zpa_scim_groups.a000
}