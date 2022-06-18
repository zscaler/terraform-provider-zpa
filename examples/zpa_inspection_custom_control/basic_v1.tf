resource "zpa_inspection_custom_controls" "tf-test01" {
  name           = "TF-Test01"
  description    = "TF-Test01"
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
    type  = "RESPONSE_BODY"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
}



