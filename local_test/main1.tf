
terraform {
  required_providers {
    zpa = {
      version = "1.0.0"
      source  = "zscaler.com/zpa/zpa"
    }
  }
  required_version = ">= 0.13"
}

provider "zpa" {}

data "zpa_global_policy_isolation" "test"{
  policy_type = "ISOLATION_POLICY"
}

resource "zpa_policy_isolation_rule" "Test1" {
  name          = "Test1"
  description   = "Test1"
  action        = "ISOLATE"
  rule_order    = 2
  operator      = "AND"
  policy_set_id = data.zpa_global_policy_isolation.test.id
  zpn_cbi_profile_id =  "216196257331286656"
  conditions {
    negated  = false
    operator = "OR"
    operands {
      name        = "zpn_client_type_exporter"
      object_type = "CLIENT_TYPE"
      lhs         = "id"
      rhs         = "zpn_client_type_exporter"
    }
  }
}