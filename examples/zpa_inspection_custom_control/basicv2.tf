# Validating RESPONSE Type

resource "zpa_inspection_custom_controls" "test10" {
  name = "Test200"
  description = "Test200"
  action = "PASS"
  default_action = "PASS"
  paranoia_level = "1"
  severity = "CRITICAL"
  type = "RESPONSE"
  # associated_inspection_profile_names {
  #     id = [data.zpa_inspection_profile.example.id, data.zpa_inspection_profile.example2.id]
  # }
  rules = [
  {
    names = [""]
    type  = "RESPONSE_HEADERS"
    conditions = {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  },
  {
    names = [""]
    type  = "RESPONSE_HEADERS"
    conditions = {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
]
}