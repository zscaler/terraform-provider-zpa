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
}

resource "zpa_tag_group" "this" {
  name        = "Example Tag Group"
  description = "An example tag group"

  tags {
    namespace_id = zpa_tag_namespace.this.id
    tag_key_id   = zpa_tag_key.this.id
    tag_value_id = zpa_tag_key.this.tag_values[0].id
  }
}
