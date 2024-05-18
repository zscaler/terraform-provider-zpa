# Retrieves ALL application segments by type
data "zpa_application_segment_by_type" "this" {
    application_type = "BROWSER_ACCESS"
}

data "zpa_application_segment_by_type" "this" {
    application_type = "INSPECT"
}

data "zpa_application_segment_by_type" "this" {
    application_type = "SECURE_REMOTE_ACCESS"
}

# Retrieves ALL application segment names by type
data "zpa_application_segment_by_type" "this" {
    application_type = "BROWSER_ACCESS"
    name = "ba_app01"
}

data "zpa_application_segment_by_type" "this" {
    application_type = "INSPECT"
    name = "inspect_app01"
}

data "zpa_application_segment_by_type" "this" {
    application_type = "SECURE_REMOTE_ACCESS"
    name = "pra_app01"
}