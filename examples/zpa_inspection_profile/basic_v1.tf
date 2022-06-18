data "zpa_inspection_all_predefined_controls" "default_predefined_controls" {
  version    = "OWASP_CRS/3.3.0"
  group_name = "preprocessors"
}

data "zpa_inspection_predefined_controls" "this" {
  name    = "Failed to parse request body"
  version = "OWASP_CRS/3.3.0"
}

resource "zpa_inspection_profile" "example" {
  name                        = "Example"
  description                 = "Example"
  paranoia_level              = "2"
  predefined_controls_version = "OWASP_CRS/3.3.0"
  incarnation_number          = "6"
  controls_info {
    control_type = "PREDEFINED"
  }
  custom_controls {
    id     = zpa_inspection_custom_controls.this.id
    action = "BLOCK"
  }
  dynamic "predefined_controls" {
    for_each = data.zpa_inspection_all_predefined_controls.default_predefined_controls.list
    content {
      id           = predefined_controls.value.id
      action       = predefined_controls.value.action == "" ? predefined_controls.value.default_action : predefined_controls.value.action
      action_value = predefined_controls.value.action_value
    }
  }
  predefined_controls {
    id     = data.zpa_inspection_predefined_controls.this.id
    action = "BLOCK"
  }
  global_control_actions = [
    "PREDEFINED:PASS",
    "CUSTOM:NONE",
    "OVERRIDE_ACTION:COMMON"
  ]
  common_global_override_actions_config = {
    "PREDEF_CNTRL_GLOBAL_ACTION" : "PASS",
    "IS_OVERRIDE_ACTION_COMMON" : "TRUE"
  }
}

resource "zpa_inspection_custom_controls" "this" {
  name           = "Example"
  description    = "Example"
  action         = "PASS"
  default_action = "PASS"
  paranoia_level = "1"
  severity       = "CRITICAL"
  type           = "RESPONSE"
  rules {
    names = ["test1", "test2", "test3"]
    type  = "RESPONSE_HEADERS"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
  rules {
    type = "RESPONSE_BODY"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
}



