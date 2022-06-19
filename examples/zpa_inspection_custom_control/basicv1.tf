resource "zpa_inspection_custom_controls" "test10" {
  name           = "Test200"
  description    = "Test200"
  action         = "PASS"
  default_action = "PASS"
  paranoia_level = "1"
  severity       = "CRITICAL"
  type           = "REQUEST"
  # type = "RESPONSE"
  # associated_inspection_profile_names {
  #     id = [data.zpa_inspection_profile.example.id, data.zpa_inspection_profile.example2.id]
  # }
  rules {
    names = [""]
    type  = "REQUEST_BODY"
    conditions {
      lhs = "SIZE"
      op  = "GE"
      rhs = "1000"
    }
  }
  rules {
    names = [""]
    type  = "REQUEST_METHOD"
    conditions {
      lhs = "VALUE"
      op  = "RX"
      rhs = "GET"
    }
  }
  rules {
    names = [""]
    type  = "REQUEST_METHOD"
    conditions {
      lhs = "VALUE"
      op  = "RX"
      rhs = "HEAD"
    }
  }
  rules {
    names = [""]
    type  = "REQUEST_METHOD"
    conditions {
      lhs = "VALUE"
      op  = "RX"
      rhs = "POST"
    }
  }
  rules {
    names = [""]
    type  = "REQUEST_METHOD"
    conditions {
      lhs = "VALUE"
      op  = "RX"
      rhs = "OPTIONS"
    }
  }
  rules {
    names = [""]
    type  = "REQUEST_METHOD"
    conditions {
      lhs = "VALUE"
      op  = "RX"
      rhs = "PUT"
    }
  }
  rules {
    names = [""]
    type  = "REQUEST_METHOD"
    conditions {
      lhs = "VALUE"
      op  = "RX"
      rhs = "DELETE"
    }
  }
  rules {
    names = [""]
    type  = "REQUEST_METHOD"
    conditions {
      lhs = "VALUE"
      op  = "RX"
      rhs = "TRACE"
    }
  }
}


