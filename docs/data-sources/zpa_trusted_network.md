---
subcategory: "Trusted Network"
layout: "zscaler"
page_title: "ZPA: trusted_network"
description: |-
  Get information about Trusted Network in Zscaler Private Access cloud.
---

# zpa_trusted_network

The **zpa_trusted_network** data source to get information about a trusted network created in the Zscaler Private Access Mobile Portal. This data source can then be referenced within the following resources:

1. Access Policy
2. Forwarding Policy
3. Inspection Policy
4. Isolation Policy
5. Service Edge Group.

## Example Usage

```hcl
# ZPA Trusted Network Data Source
data "zpa_trusted_network" "example" {
 name = "trusted_network_name"
}

```

-> **NOTE** To query trusted network that are associated with a specific Zscaler cloud, it is required to append the cloud name to the name of the trusted network as the below example:

```hcl
# ZPA Posture Profile Data Source
data "zpa_trusted_network" "example1" {
 name = "Corporate-Network (zscalertwo.net)"
}

output "zpa_trusted_network" {
  value = data.zpa_trusted_network.example1.network_id
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
* `network_id` - (Computed)
* `zscaler_cloud` - (Computed)
* `zscaler_customer_id` - (Computed)
