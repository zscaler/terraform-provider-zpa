data "zpa_policy_type" "this" {
  policy_type = "REDIRECTION_POLICY"
}


data "zpa_service_edge_group" "this" {
 name = "Example"
}

resource "zpa_policy_redirection_rule" "this" {
  name                          = "Example"
  description                   = "Example"
  action                        = "REDIRECT_ALWAYS"
  operator = "AND"
  policy_set_id = data.zpa_policy_type.this.id

  conditions {
    negated = false
    operator = "OR"
	    operands {
            object_type = "CLIENT_TYPE"
            lhs         = "id"
            rhs         = "zpn_client_type_branch_connector"
	}
	    operands {
            object_type = "CLIENT_TYPE"
            lhs         = "id"
            rhs         = "zpn_client_type_edge_connector"
	}
  }
  service_edge_groups {
    id = [ data.zpa_service_edge_group.this.id ]
  }
}

// ZPA Private Service Edge groups must be empty when the Private Service Edge Selection Method is Default.