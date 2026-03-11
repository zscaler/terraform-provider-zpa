resource "zpa_tag_namespace" "this" {
  name        = "Example Namespace"
  description = "An example tag namespace"
  enabled     = true
}

resource "zpa_tag_key" "this" {
  name         = "Environment"
  description  = "Environment tag key"
  enabled      = true
  namespace_id = zpa_tag_namespace.this.id

  tag_values {
    name = "Production"
  }

  tag_values {
    name = "Staging"
  }
}
