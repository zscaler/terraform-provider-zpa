---
page_title: "zpa_managed_browser_profile Data Source - terraform-provider-zpa"
subcategory: "Managed Browser Profile"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-browser-protection-profiles
  API documentation https://help.zscaler.com/zpa/about-browser-protection-profiles
  Get information about location resource
---

# zpa_managed_browser_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-browser-protection-profiles)
* [API documentation](https://help.zscaler.com/zpa/about-browser-protection-profiles)

Use the **zpa_managed_browser_profile** data source to get information about managed browser protection profiles within the Zscaler Private Access cloud. This data source can be used when configuring `zpa_policy_access_rule_v2` or `zpa_policy_isolation_rule_v2` where the `object_type` is `CHROME_POSTURE_PROFILE`

## Example Usage

```terraform
data "zpa_managed_browser_profile" "this" {
  name = "Profile01"
}
```

## Schema

### Required

At least one of the following arguments must be provided:

* `name` - (String) The name of the managed browser profile to be retrieved.
* `id` - (String) The ID of the managed browser profile to be retrieved.

### Optional

* `microtenant_id` - (String) The microtenant ID the middleware will use when calling the API endpoint.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (String) The unique identifier of the managed browser profile.
* `name` - (String) The name of the managed browser profile.
* `description` - (String) The description of the managed browser profile.
* `browser_type` - (String) The type of browser for this profile (e.g., "CHROME").
* `customer_id` - (String) The unique identifier of the customer associated with this profile.
* `microtenant_id` - (String) The unique identifier of the microtenant associated with this profile.
* `microtenant_name` - (String) The name of the microtenant associated with this profile.
* `creation_time` - (String) The timestamp when the managed browser profile was created.
* `modified_by` - (String) The identifier of the user who last modified the managed browser profile.
* `modified_time` - (String) The timestamp when the managed browser profile was last modified.
* `chrome_posture_profile` - (List) Chrome posture profile configuration block. Contains:
  * `id` - (String) The unique identifier of the chrome posture profile.
  * `browser_type` - (String) The type of browser for this posture profile (e.g., "CHROME").
  * `crowd_strike_agent` - (Boolean) Indicates whether CrowdStrike agent is enabled for this posture profile.
  * `creation_time` - (String) The timestamp when the chrome posture profile was created.
  * `modified_by` - (String) The identifier of the user who last modified the chrome posture profile.
  * `modified_time` - (String) The timestamp when the chrome posture profile was last modified.