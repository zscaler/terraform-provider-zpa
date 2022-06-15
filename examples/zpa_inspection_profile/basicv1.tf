resource "zpa_inspection_profile" "example" {
  name                        = "Example"
  description                 = "Example"
  paranoia_level              = "2"
  predefined_controls_version = "OWASP_CRS/3.3.0"
  incarnation_number          = "6"
  controls_info {
    control_type = "PREDEFINED"
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
    id     = "72057594037930388"
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

data "zpa_inspection_all_predefined_controls" "default_predefined_controls" {
  version    = "OWASP_CRS/3.3.0"
  group_name = "preprocessors"
}

output "zpa_inspection_profile" {
  value = resource.zpa_inspection_profile.example
}

output "zpa_inspection_all_predefined_controls" {
  value = data.zpa_inspection_all_predefined_controls.default_predefined_controls.list
}

