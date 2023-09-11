// Retrieve the CBI Banner to be associated with the external profile
data "zpa_cloud_browser_isolation_banner" "this" {
  name = "Default"
}

// Retrieve the CBI Region ID where the profile will be created (At least 2 required)
data "zpa_cloud_browser_isolation_region" "singapore" {
  name = "Singapore"
}

data "zpa_cloud_browser_isolation_region" "frankfurt" {
  name = "Frankfurt"
}

// Retrieve the CBI Certificate ID to be associated with the external profile
data "zpa_cloud_browser_isolation_certificate" "this" {
  name = "Zscaler Root Certificate"
}

resource "zpa_cloud_browser_isolation_external_profile" "this" {
  name            = "CBI Profile"
  description     = "CBI Profile"
  banner_id       = data.zpa_cloud_browser_isolation_banner.this.id
  region_ids      = [data.zpa_cloud_browser_isolation_region.singapore.id, data.zpa_cloud_browser_isolation_region.frankfurt.id]
  certificate_ids = [data.zpa_cloud_browser_isolation_certificate.this.id]
  user_experience {
    session_persistence = true
    browser_in_browser  = true
  }
  security_controls {
    copy_paste          = "all"
    upload_download     = "all"
    document_viewer     = true
    local_render        = true
    allow_printing      = true
    restrict_keystrokes = false
  }
}