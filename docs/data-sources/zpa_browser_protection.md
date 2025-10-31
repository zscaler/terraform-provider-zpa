---
page_title: "zpa_browser_protection Data Source - terraform-provider-zpa"
subcategory: "Browser Protection Profile"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-browser-protection-profiles
  API documentation https://help.zscaler.com/zpa/about-browser-protection-profiles
  Get information about location resource
---

# zpa_browser_protection (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-browser-protection-profiles)
* [API documentation](https://help.zscaler.com/zpa/about-browser-protection-profiles)

Use the **zpa_browser_protection** data source to get information about managed browser protection profiles within the Zscaler Private Access cloud. This data source can be used when configuring `zpa_policy_browser_protection_rule`.

## Example Usage

```terraform
data "zpa_browser_protection" "this" {
  name = "Profile01"
}
```

## Schema

### Required

At least one of the following arguments must be provided:

* `name` - (String) The name of the browser protection profile to be retrieved. If not provided, the default profile will be returned.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The unique identifier of the browser protection profile.
* `name` - (String) The name of the browser protection profile.
* `description` - (String) The description of the browser protection profile.
* `default_csp` - (Boolean) Indicates whether this is the default Content Security Policy profile.
* `creation_time` - (String) The timestamp when the browser protection profile was created.
* `modified_by` - (String) The identifier of the user who last modified the browser protection profile.
* `modified_time` - (String) The timestamp when the browser protection profile was last modified.
* `criteria_flags_mask` - (String) The criteria flags mask used for browser protection evaluation.
* `criteria` - (List) Browser protection criteria configuration block. Contains:
  * `finger_print_criteria` - (List) Fingerprint criteria configuration block. Contains:
    * `collect_location` - (Boolean) Indicates whether location data should be collected for fingerprinting.
    * `fingerprint_timeout` - (String) The timeout value for fingerprint collection (in seconds).
    * `browser` - (List) Browser-specific fingerprinting criteria. Contains:
      * `browser_eng` - (Boolean) Collect browser engine information.
      * `browser_eng_ver` - (Boolean) Collect browser engine version.
      * `browser_name` - (Boolean) Collect browser name.
      * `browser_version` - (Boolean) Collect browser version.
      * `canvas` - (Boolean) Collect canvas fingerprinting data.
      * `flash_ver` - (Boolean) Collect Flash version information.
      * `fp_usr_agent_str` - (Boolean) Collect user agent string for fingerprinting.
      * `is_cookie` - (Boolean) Collect cookie information.
      * `is_local_storage` - (Boolean) Collect local storage information.
      * `is_sess_storage` - (Boolean) Collect session storage information.
      * `ja3` - (Boolean) Collect JA3 fingerprint information.
      * `mime` - (Boolean) Collect MIME type information.
      * `plugin` - (Boolean) Collect browser plugin information.
      * `silverlight_ver` - (Boolean) Collect Silverlight version information.
    * `location` - (List) Location-based fingerprinting criteria. Contains:
      * `lat` - (Boolean) Collect latitude information.
      * `lon` - (Boolean) Collect longitude information.
    * `system` - (List) System-level fingerprinting criteria. Contains:
      * `avail_screen_resolution` - (Boolean) Collect available screen resolution.
      * `cpu_arch` - (Boolean) Collect CPU architecture information.
      * `curr_screen_resolution` - (Boolean) Collect current screen resolution.
      * `font` - (Boolean) Collect font information.
      * `java_ver` - (Boolean) Collect Java version information.
      * `mobile_dev_type` - (Boolean) Collect mobile device type information.
      * `monitor_mobile` - (Boolean) Monitor mobile device characteristics.
      * `os_name` - (Boolean) Collect operating system name.
      * `os_version` - (Boolean) Collect operating system version.
      * `sys_lang` - (Boolean) Collect system language information.
      * `tz` - (Boolean) Collect timezone information.
      * `usr_lang` - (Boolean) Collect user language information.