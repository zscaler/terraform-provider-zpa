data "zpa_tag_namespace" "this" {
  name = "Example Namespace"
}

output "zpa_tag_namespace_id" {
  value = data.zpa_tag_namespace.this.id
}
