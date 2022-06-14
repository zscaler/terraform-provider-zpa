data "zpa_inspection_profile" "this" {
  name = "ZTX-Inspection-Profile"
}

resource "zpa_inspection_custom_controls" "tf-test01" {
  name           = "TF-Test01"
  description    = "TF-Test01"
  action         = "PASS"
  default_action = "PASS"
  paranoia_level = "1"
  severity       = "CRITICAL"
  type           = "REQUEST"
  associated_inspection_profile_names {
    id = [data.zpa_inspection_profile.this.id]
  }
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
    type  = "RESPONSE_BODY"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
}



