terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}


data "zpa_machine_group" "all" {
  name = "MGR01"
}

output "all_machine_group" {
  value = data.zpa_machine_group.all.id
}