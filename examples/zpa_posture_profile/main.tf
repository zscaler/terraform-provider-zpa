terraform {
    required_providers {
        zpa = {
            version = "1.0.0"
            source = "zscaler.com/zpa/zpa"
        }
    }
}

provider "zpa" {}


// Testing Data Source Posture Profile
// data "zpa_posture_profile" "example" {
//  name = "sgio-windows-domain"
// }

// output "all_posture_profile" {
//   value = data.zpa_posture_profile.example
// }

data "zpa_posture_profile" "crwd_zta_score_40" {
 name = "CrowdStrike_ZPA_ZTA_40"
}

output "all_posture_profile" {
  value = data.zpa_posture_profile.crwd_zta_score_40.id
}