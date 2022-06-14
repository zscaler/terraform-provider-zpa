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
    names = ["test", "test1", "test2"]
    type  = "REQUEST_HEADERS"
    conditions {
      lhs = "SIZE"
      op  = "EQ"
      rhs = "1000"
    }
  }
  rules {
    names = ["test", "test1", "test2"]
    type  = "REQUEST_COOKIES"
    conditions {
      lhs = "SIZE"
      op  = "LE"
      rhs = "1000"
    }
  }
  rules {
    type = "REQUEST_URI"
    conditions {
      lhs = "VALUE"
      op  = "CONTAINS"
      rhs = "tf-test"
    }
  }
  rules {
    type = "QUERY_STRING"
    conditions {
      lhs = "VALUE"
      op  = "STARTS_WITH"
      rhs = "tf-test"
    }
  }
}



