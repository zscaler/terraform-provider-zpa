---
subcategory: "Posture Profiles"
layout: "zpa"
page_title: "ZPA: posture profile"
description: |-
  Gets a ZPA Posture Profile details.

---
# zpa_posture_profile

The **zpa_posture_profile** data source provides details about a specific posture profile created in the Zscaler Private Access Mobile Portal.
This data source is required when creating:

1. Access policy Rule
2. Access policy timeout rule
3. Access policy forwarding rule

## Example Usage

```hcl
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example1" {
 name = "CrowdStrike_ZPA_ZTA_40"
}
```

```hcl
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example2" {
 name = "Detect SentinelOne"
}
```

```hcl
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example3" {
 name = "domain_joined"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the posture profile to be exported.
* `domain` - (Optional)
* `posture_udid` - (Optional)
* `zscaler_cloud` - (Optional)
* `zscaler_customer_id` - (Optional)
