// Retrieve "Default" customer version profile
data "zpa_customer_version_profile" "default"{
    name = "Default"
}

output "zpa_customer_version_profile_default" {
    value = data.zpa_customer_version_profile.default
}

// Retrieve "Previous Default" customer version profile
data "zpa_customer_version_profile" "previous_default"{
    name = "Previous Default"
}

output "zpa_customer_version_profile_previous_default" {
    value = data.zpa_customer_version_profile.previous_default
}

// Retrieve "New Release" customer version profile
data "zpa_customer_version_profile" "new_release"{
    name = "New Release"
}

output "zpa_customer_version_profile_new_release" {
    value = data.zpa_customer_version_profile.new_release
}