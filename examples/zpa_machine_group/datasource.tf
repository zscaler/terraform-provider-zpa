data "zpa_machine_group" "example" {
  name = "Example-MGR01"
}

output "all_machine_group" {
  value = data.zpa_machine_group.example
}