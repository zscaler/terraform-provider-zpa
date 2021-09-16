data "zpa_segment_group" "all" { 
  name = "Browser Access Apps"
}

output "segment_group" {
    value = data.zpa_segment_group.all.id
}
