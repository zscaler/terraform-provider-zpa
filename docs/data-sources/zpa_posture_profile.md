---
page_title: "zpa_posture_profile Data Source - terraform-provider-zpa"
subcategory: "Posture Profile"
description: |-
  Official documentation https://help.zscaler.com/client-connector/about-device-posture-profiles
  API documentation https://help.zscaler.com/zpa/obtaining-posture-profile-details-using-api
  Get information about Posture Profile in Zscaler Private Access cloud.
---

# zpa_posture_profile (Data Source)

* [Official documentation](https://help.zscaler.com/client-connector/about-device-posture-profiles)
* [API documentation](https://help.zscaler.com/zpa/obtaining-posture-profile-details-using-api)

Use the **zpa_posture_profile** data source to get information about a posture profile created in the Zscaler Private Access Mobile Portal. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

## Example Usage

```terraform
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example1" {
 name = "CrowdStrike_ZPA_ZTA_40"
}
```

```terraform
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example2" {
 name = "Detect SentinelOne"
}
```

```terraform
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example3" {
 name = "domain_joined"
}
```

-> **NOTE** To query posture profiles that are associated with a specific Zscaler cloud, it is required to append the cloud name to the name of the posture profile as the below example:

```terraform
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example1" {
 name = "CrowdStrike_ZPA_ZTA_40 (zscalertwo.net)"
}
```

-> **NOTE** When associating a posture profile with one of supported resources, the following parameter must be exported: ``posture_udid`` instead of the ``id`` of the resource.

```terraform
# ZPA Posture Profile Data Source
data "zpa_posture_profile" "example1" {
 name = "CrowdStrike_ZPA_ZTA_40 (zscalertwo.net)"
}

output "zpa_posture_profile" {
  value = data.zpa_posture_profile.example1.posture_udid
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the posture profile to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the posture profile to be exported.
* `creation_time` - (string)
* `domain` - (string)
* `master_customer_id` - (string)
* `modified_by` - (string)
* `modified_time` - (string)
* `posture_udid` - (string)
* `zscaler_cloud` - (string)
* `zscaler_customer_id` - (string)
