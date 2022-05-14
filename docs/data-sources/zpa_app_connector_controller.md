---
subcategory: "App Connector Management"
layout: "zpa"
page_title: "ZPA: app_connector_controller"
description: |-
  Gets a ZPA App Connector details.

---
# zpa_app_connector_controller

The **zpa_app_connector_controller** data source provides details about a specific app connector created in the Zscaler Private Access cloud. This data source can then be referenced in the following resources:

* App Connector Group

## Example Usage

```hcl
# ZPA App Connector Data Source
data "zpa_app_connector" "example" {
  name = "AWS-VPC100-App-Connector"
}

output "zpa_app_connector" {
  value = data.zpa_app_connector.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the App Connector Group.
* `description` (Computed) Description of the App Connector Group.
* `enabled` - (Computed) Whether this App Connector Group is enabled or not. Default value: `true`. Supported values: `true`, `false`
* `latitude` - (Computed) Latitude of the App Connector Group. Integer or decimal. With values in the range of `-90` to `90`
* `longitude` - (Computed) Longitude of the App Connector Group. Integer or decimal. With values in the range of `-180` to `180`
