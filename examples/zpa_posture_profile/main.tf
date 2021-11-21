data "zpa_posture_profile" "crwd_zta_score_40" {
 name = "CrowdStrike_ZPA_ZTA_40"
}

output "all_posture_profile" {
  value = data.zpa_posture_profile.crwd_zta_score_40
}