---
subcategory: "Posture Profile"
layout: "zscaler"
page_title: "ZPA: posture_profile"
description: |-
  Get information about Posture Profile in Zscaler Private Access cloud.
---

# Data Source: zpa_posture_profile

Use the **zpa_posture_profile** data source to get information about a posture profile created in the Zscaler Private Access Mobile Portal. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

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

-> **NOTE** To query posture profiles that are associated with a specific Zscaler cloud, it is required to append the cloud name to the name of the posture profile as the below example:

```hcl
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example1" {
 name = "CrowdStrike_ZPA_ZTA_40 (zscalertwo.net)"
}
```

-> **NOTE** When associating a posture profile with one of supported resources, the following parameter must be exported: ``posture_udid`` instead of the ``id`` of the resource.

```hcl
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example1" {
 name = "CrowdStrike_ZPA_ZTA_40 (zscalertwo.net)"
}

output "zpa_posture_profile" {
  value = data.zpa_posture_profile.example1.posture_udid
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the posture profile to be exported.
* `id` - (Optional) The ID of the posture profile to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `creation_time` - (Computed)
* `domain` - (Computed)
* `master_customer_id` - (Computed)
* `modified_by` - (Computed)
* `modified_time` - (Computed)
* `posture_udid` - (Computed)
* `zscaler_cloud` - (Computed)
* `zscaler_customer_id` - (Computed)
