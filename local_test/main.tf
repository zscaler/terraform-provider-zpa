terraform {
  required_providers {
    zpa = {
      version = "2.0.5"
      source  = "zscaler.com/zpa/zpa"
    }
  }
  required_version = ">= 0.13"
}

provider "zpa" {}

data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

locals {
  app_group = [for rhs_value in tolist(data.zpa_segment_group.zpa_segment_group) : { rhs = rhs_value.id } ]
}

resource "zpa_policy_access_rule" "all_other_services" {
  name                          = "All Other Services"
  description                   = "All Other Services"
  action                        = "ALLOW"
  rule_order                     = 2
  operator = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id

  conditions {
    negated = false
    operator = "OR"
    dynamic "operands" {
      for_each = local.app_group
      content {
      name =  "All Other Services"
      object_type = "APP_GROUP"
      lhs         = "id"
      rhs = operands.value.rhs
      }
    operands {
      name =  "All Other Services"
      object_type = "APP"
      lhs = "id"
      rhs = "216196257331292173"
    }
    operands {
      name =  "All Other Services"
      object_type = "APP"
      lhs = "id"
      rhs = "216196257331292105"
    }
  }
}


// data "zpa_idp_controller" "example" {
//  name = "SGIO-User-Okta"
// }

// output "idp_controller" {
//     value = data.zpa_idp_controller.example
// }


// resource "zpa_idp_controller" "okta_test" {
//   name = "Okta-Test"
//   description = "Okta-Test"
//   enabled = true
//   disable_saml_based_policy = false
//   reauth_on_user_update = false
//   domain_list = ["216196257331281920.zpa-customer.com"]
//   sso_type = ["USER"]
//   use_custom_sp_metadata = true
//   sign_saml_request = "1"
//   auto_provision = "0"
//   enable_scim_based_policy = false
//   idp_entity_id = "http://www.okta.com/exk4h5tzk0cjN50wm4x7"
//   login_url = "https://dev-151399.okta.com/app/zscaler_private_access/exk4h5tzk0cjN50wm4x7/sso/saml"
//   user_sp_signing_cert_id = "0"
//   // zpa_saml_request = "1"
// }

// output "idp_controller" {
//     value = zpa_idp_controller.okta_test
// }
