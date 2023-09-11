---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: cloud_browser_isolation_external_profile"
description: |-
  Creates and manages Cloud Browser Isolation External Profile.
---

# Resource: zpa_cloud_browser_isolation_external_profile

The **zpa_cloud_browser_isolation_external_profile** resource creates a Cloud Browser Isolation external profile. This resource can then be used in as part of `zpa_policy_isolation_rule` when the `action` attribute is set to `ISOLATE`.

## Example Usage

```hcl
# Retrieve CBI Banner ID
data "zpa_cloud_browser_isolation_banner" "this" {
  name = "Default"
}

# Retrieve Primary CBI Region ID
data "zpa_cloud_browser_isolation_region" "singapore" {
  name = "Singapore"
}

# Retrieve Secondary CBI Region ID
data "zpa_cloud_browser_isolation_region" "frankfurt" {
  name = "Frankfurt"
}

# Retrieve CBI Certificate ID
data "zpa_cloud_browser_isolation_certificate" "this" {
    name = "Zscaler Root Certificate"
}

resource "zpa_cloud_browser_isolation_external_profile" "this" {
    name = "CBI_Profile_Example"
    description = "CBI_Profile_Example"
    banner_id = data.zpa_cloud_browser_isolation_banner.this.id
    region_ids = [data.zpa_cloud_browser_isolation_region.singapore.id]
    certificate_ids = [data.zpa_cloud_browser_isolation_certificate.this.id]
    user_experience {
      session_persistence = true
      browser_in_browser = true
    }
    security_controls {
      copy_paste = "all"
      upload_download = "all"
      document_viewer = true
      local_render = true
      allow_printing = true
      restrict_keystrokes = false
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the CBI banner to be exported.
* `banner_id` - (Required) The ID of the CBI banner to be exported.
* `certificate_ids` - (Optional) The CBI security controls enabled for the profile
  * `id:` - (Optional) The ID of the CBI Certificate to be associated with the profile.

* `region_ids` - (Optional) The CBI region
  * `id:` - (Optional) The ID of CBI region where the profile must be deployed. At least 2 regions are required.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (Optional) - The description of the CBI profile

* `security_controls` - (Optional) The CBI security controls enabled for the profile
  * `copy_paste:` - (Optional) Enable or disable copy & paste for local computer to isolation. Supported values are: `none` or `all`
  * `document_viewer:` - (Optional) Enable or disable to view Microsoft Office files in isolation.
  * `local_render:` - (Optional) Enables non-isolated hyperlinks to be opened on the user's native browser.
  * `upload_download:` - (Optional) Enable or disable file transfer from local computer to isolation. Supported values are: `none` or `all`
  * `allow_printing:` - (Optional) Enables the user to print web pages and documents rendered within the isolation browser. Supported values are: `true` or `false`
  * `restrict_keystrokes:` - (Optional) Prevents keyboard and text input to isolated web pages. Supported values are: `true` or `false`

* `user_experience` - The CBI security controls enabled for the profile
  * `session_persistence:` - (Optional) Save user cookies between sessions. If disabled, all cookies will be discarded when isolation session ends. Supported values are: `true` or `false`
  * `browser_in_browser:` - (Optional) Enable or disable browser-in-browser or native browser experience. Supported values are: `true` or `false`
