---
subcategory: "Cloud Browser Isolation"
layout: "zscaler"
page_title: "ZPA: cloud_browser_isolation_external_profile"
description: |-
  Get information about Cloud Browser Isolation External Profile.
---

# Data Source: zpa_cloud_browser_isolation_external_profile

Use the **zpa_cloud_browser_isolation_external_profile** data source to get information about Cloud Browser Isolation external profile. This data source information can then be used in as part of `zpa_policy_isolation_rule` when the `action` attribute is set to `ISOLATE`.

## Example Usage

```hcl
# Retrieve CBI External Profile
data "zpa_cloud_browser_isolation_external_profile" "this" {
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the CBI banner to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - The description of the CBI profile
* `is_default` - (bool) - Indicates if the CBI profile is the default one.
* `href` - (string)
* `regions` - (string) List of regions where multi-region deployment is enabled
  * `id:` - (string) Region ID where the profile is applied to
  * `name:` - (string) Region name where the profile is applied to

* `security_controls` - The CBI security controls enabled for the profile
  * `copy_paste:` - (string) Enable or disable copy & paste for local computer to isolation
  * `document_viewer:` - (bool) Enable or disable to view Microsoft Office files in isolation.
  * `local_render:` - (bool) Enables non-isolated hyperlinks to be opened on the user's native browser.
  * `upload_download:` - (string) Enable or disable file transfer from local computer to isolation
  * `allow_printing:` - (bool) Enables the user to print web pages and documents rendered within the isolation browser.
  * `restrict_keystrokes:` - (bool) Prevents keyboard and text input to isolated web pages.

* `user_experience` - The CBI security controls enabled for the profile
  * `session_persistence:` - (bool) Save user cookies between sessions. If disabled, all cookies will be discarded when isolation session ends.
  * `browser_in_browser:` - (bool) Enable or disable browser-in-browser or native browser experience
