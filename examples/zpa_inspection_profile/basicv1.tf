resource "zpa_inspection_profile" "example" {
  name = "Example"
  description = "Example"
  paranoia_level = "2"
  predefined_controls_version = "OWASP_CRS/3.3.0"
  incarnation_number = "6"
  custom_controls {
      id = [ "216196257331305413" ]
  }
  predefined_controls {
      id = [ "72057594037930391"]
  }
  controls_info {
    control_type = "PREDEFINED"
  }
  global_control_actions = [
          "PREDEFINED:PASS",
          "CUSTOM:NONE",
          "OVERRIDE_ACTION:COMMON"
  ]
  common_global_override_actions_config = {
          "PREDEF_CNTRL_GLOBAL_ACTION": "PASS",
          "IS_OVERRIDE_ACTION_COMMON": "TRUE"
  }
}