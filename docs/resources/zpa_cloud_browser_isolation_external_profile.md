---
page_title: "zpa_cloud_browser_isolation_external_profile Resource - terraform-provider-zpa"
subcategory: "Cloud Browser Isolation"
description: |-
  Official documentation https://help.zscaler.com/isolation/about-custom-root-certificates-cloud-browser-isolation
  Creates and manages Cloud Browser Isolation External Profile.
---

# zpa_cloud_browser_isolation_external_profile (Resource)

* [Official documentation](https://help.zscaler.com/isolation/about-custom-root-certificates-cloud-browser-isolation)

The **zpa_cloud_browser_isolation_external_profile** resource creates a Cloud Browser Isolation external profile. This resource can then be used in as part of `zpa_policy_isolation_rule` when the `action` attribute is set to `ISOLATE`.

## Example Usage

```terraform
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
      forward_to_zia {
        enabled         = true
        organization_id = "***********"
        cloud_name      = "<cloud_name>"
        pac_file_url    = "https://pac.<cloud_name>/<cloud_name>/proxy.pac"
      }
      browser_in_browser     = true
      persist_isolation_bar  = true
      translate              = true
      session_persistence    = true
    }
    security_controls {
      copy_paste          = "all"
      upload_download     = "upstream"
      document_viewer     = true
      local_render        = true
      allow_printing      = true
      restrict_keystrokes = true
      flattened_pdf       = true
      deep_link {
        enabled           = true
        applications      = ["test1", "test"]
      }
      watermark {
        enabled          = true
        show_user_id     = true
        show_timestamp   = true
        show_message     = true
        message          = "Zscaler CBI"
      }
    }
    debug_mode {
      allowed             = true
      file_password       = "***********"
    }
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (Required) The name of the CBI banner to be exported.
- `banner_id` - (Required) The ID of the CBI banner to be exported.
- `certificate_ids` - (Optional) The CBI security controls enabled for the profile
  - `id:` - (Optional) The ID of the CBI Certificate to be associated with the profile.

- `region_ids` - (Optional) The CBI region
  - `id:` - (Optional) The ID of CBI region where the profile must be deployed. At least 2 regions are required.

### Optional

In addition to all arguments above, the following attributes are exported:

- `description` - (Optional) - The description of the CBI profile
- `flattened_pdf` - (Optional) - Enable to allow downloading of flattened files from isolation container to your local computer.

    **NOTE** `flattened_pdf` must be set to `false` when `upload_download` is set to `all`

- `security_controls` - (Optional) The CBI security controls enabled for the profile
  - `copy_paste:` - (Optional) Enable or disable copy & paste for local computer to isolation. Supported values are: `none` or `all`
  - `document_viewer:` - (Optional) Enable or disable to view Microsoft Office files in isolation.
  - `local_render:` - (Optional) Enables non-isolated hyperlinks to be opened on the user's native browser.
  - `upload_download` - (Optional) Enable or disable file transfer from local computer to isolation. Supported values are: `none`, `all`, `upstream`

    **NOTE** `upload_download` must be set to `none` or `upstream` when `flattened_pdf` is set to `true`

  - `allow_printing` - (Optional) Enables the user to print web pages and documents rendered within the isolation browser. Supported values are: `true` or `false`

  - `restrict_keystrokes:` - (Optional) Prevents keyboard and text input to isolated web pages. Supported values are: `true` or `false`
  - `deep_link:` - (Optional) Enter applications that are allowed to launch outside of the Isolation session
    - `enabled:` - (Optional) Enable or disable to view Microsoft Office files in isolation.
    - `applications:` - (Optional) List of deep link applications

  - `watermark:` - (Optional) Enable to display a custom watermark on isolated web pages.
    - `enabled:` - (Optional) Enable to display a custom watermark on isolated web pages.
    - `show_user_id:` - (Optional) Display the user ID on watermark isolated web pages.
    - `show_timestamp:` - (Optional) Display the timestamp on watermark isolated web pages.
    - `show_message:` - (Optional) Enable custom message on watermark isolated web pages.
    - `message:` - (Optional) Display custom message on watermark isolated web pages.

- `user_experience` - The CBI security controls enabled for the profile
  - `session_persistence:` - (Optional) Save user cookies between sessions. If disabled, all cookies will be discarded when isolation session ends. Supported values are: `true` or `false`
  - `browser_in_browser:` - (Optional) Enable or disable browser-in-browser or native browser experience. Supported values are: `true` or `false`
  - `forward_to_zia:` - (Optional) Enable to forward non-ZPA Internet traffic via ZIA.
    - `enabled:` - (Optional) Enable to forward non-ZPA Internet traffic via ZIA.
    - `organization_id:` - (Optional) Use the ZIA organization ID from the Company Profile section.
    - `cloud_name:` - (Optional) The ZIA cloud name on which the organization exists i.e `zscalertwo`
    - `pac_file_url:` - (Optional) Enable to have the PAC file be configured on the Isolated browser to forward traffic via ZIA.

- `debug_mode`- (Optional) Enable to allow starting isolation sessions in debug mode to collect troubleshooting information.
  - `allowed:` - (Optional)  Enable to allow starting isolation sessions in debug mode to collect troubleshooting information.
  - `file_password:` - (Optional) Set an optional password to debug files when this mode is enabled.
