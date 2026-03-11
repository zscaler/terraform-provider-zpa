data "zpa_tag_namespace" "this" {
  name = "Example Namespace"
}

data "zpa_tag_key" "this" {
  name         = "Environment"
  namespace_id = data.zpa_tag_namespace.this.id
}

output "zpa_tag_key_id" {
  value = data.zpa_tag_key.this.id
}
