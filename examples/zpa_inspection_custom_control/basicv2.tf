# Validating RESPONSE Type

provider "zpa" {}
resource "zpa_inspection_custom_controls" "test300" {
  name           = "Test3000"
  description    = "Test3000"
  action         = "PASS"
  default_action = "PASS"
  paranoia_level = "1"
  severity       = "CRITICAL"
  type           = "RESPONSE"
  associated_inspection_profile_names {
    id = [data.zpa_inspection_profile.example.id]
  }
  rules {
    names = ["test"]
    type  = "RESPONSE_HEADERS"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
  rules {
    names = []
    type  = "RESPONSE_BODY"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
}

output "zpa_inspection_custom_controls" {
  value = zpa_inspection_custom_controls.test300
}

data "zpa_inspection_profile" "example" {
  name = "Test100"
}
