---
page_title: "zpa_trusted_network Data Source - terraform-provider-zpa"
subcategory: "Trusted Network"
description: |-
  Official documentation https://help.zscaler.com/client-connector/about-trusted-networks
  API documentation https://help.zscaler.com/zpa/obtaining-trusted-network-details-using-api
  Get information about Trusted Network in Zscaler Private Access cloud.
---

# zpa_trusted_network (Data Source)

* [Official documentation](https://help.zscaler.com/client-connector/about-trusted-networks)
* [API documentation](https://help.zscaler.com/zpa/obtaining-trusted-network-details-using-api)

The **zpa_trusted_network** data source to get information about a trusted network created in the Zscaler Private Access Mobile Portal. This data source can then be referenced within the following resources:

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

1. Access Policy
2. Forwarding Policy
3. Inspection Policy
4. Isolation Policy
5. Service Edge Group.

## Example Usage

```terraform
# ZPA Trusted Network Data Source
data "zpa_trusted_network" "example" {
 name = "trusted_network_name"
}

```

-> **NOTE** To query trusted network that are associated with a specific Zscaler cloud, it is required to append the cloud name to the name of the trusted network as the below example:

```terraform
# ZPA Posture Profile Data Source
data "zpa_trusted_network" "example1" {
 name = "Corporate-Network (zscalertwo.net)"
}

output "zpa_trusted_network" {
  value = data.zpa_trusted_network.example1.network_id
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
* `network_id` - (string)
* `zscaler_cloud` - (string)
* `zscaler_customer_id` - (string)
