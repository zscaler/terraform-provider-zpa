// Retrieve Policy Types
data "zpa_policy_type" "access_policy" {
    policy_type = "ACCESS_POLICY"
}

data "zpa_application_segment" "this"{
  name = "Example"
}

data "zpa_idp_controller" "user_idp_name" {
 name = "User_IDP_Name"
}

data "zpa_scim_groups" "engineering" {
  name     = "Engineering"
  idp_name = "User_IDP_Name"
}
resource "zpa_policy_access_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "ALLOW"
  rule_order                     = 2
  operator = "AND"
  policy_set_id = data.zpa_policy_type.access_policy.id

  conditions {
    negated = false
    operator = "OR"
    operands {
      name =  "Example"
      object_type = "APP"
      lhs = "id"
      rhs = data.zpa_application_segment.this.id
    }
  }

  conditions {
     negated = false
     operator = "OR"
    operands {
      object_type = "IDP"
      lhs = "id"
      rhs = data.zpa_idp_controller.user_idp_name.id
    }
    operands {
      object_type = "SCIM_GROUP"
      lhs = data.zpa_idp_controller.user_idp_name.id
      rhs = data.zpa_scim_groups.engineering.id
      idp_id = data.zpa_idp_controller.user_idp_name.id
    }
  }
}
