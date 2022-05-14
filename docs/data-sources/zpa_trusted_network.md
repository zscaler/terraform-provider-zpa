---
subcategory: "Trusted Network"
layout: "zpa"
page_title: "ZPA: trusted network"
description: |-
  Gets a ZPA Trusted Network details.

---

# zpa_trusted_network

The **zpa_trusted_network** data source provides details about a specific trusted network created in the Zscaler Private Access Mobile Portal.
This data source is required when creating:

1. Access policy Rule
2. Access policy timeout rule
3. Access policy forwarding rule
4. Service Edge Group

## Example Usage

```hcl
# ZPA Trusted Network Data Source
data "zpa_trusted_network" "example" {
 name = "trusted_network_name"
}

output "zpa_trusted_network" {
  value = data.zpa_trusted_network.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the trusted network to be exported.
* `domain` - (Optional)
* `network_id` - (Optional)
* `zscaler_cloud` - (Optional)
* `master_customer_id` - (Optional)
