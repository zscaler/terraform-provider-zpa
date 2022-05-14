---
subcategory: "Service Edge Controller"
layout: "zpa"
page_title: "ZPA: service_edge_controller"
description: |-
  Retrieves details of a ZPA Service Edge Controller.
  
---
# zpa_service_edge_controller

The **zpa_service_edge_controller** data source provides details about a specific Service Edge Controller created in the Zscaler Private Access cloud. This data source can then be referenced in the following resources:

* Service Edge Group
* Provisioning Key

## Example Usage

```hcl
# ZPA Service Edge Controller Data Source
data "zpa_service_edge_controller" "example" {
  name = "On-Prem-PSE"
}

output "zpa_service_edge_controller" {
  value = data.zpa_service_edge_controller.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the Service Edge Controller.
* `description` (Computed) Description of the Service Edge Controller.
* `enabled` - (Computed) Whether this Service Edge Controller is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `latitude` - (Computed) Latitude of the Service Edge Controller. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (Computed) Longitude of the Service Edge Controller. Integer or decimal. With values in the range of `-180` to `180`
* `location` - (Computed) Location of the Service Edge Controller.
