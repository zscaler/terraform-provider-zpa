data "zpa_tag_group" "this" {
  name = "Example Tag Group"
}

output "zpa_tag_group_id" {
  value = data.zpa_tag_group.this.id
}
