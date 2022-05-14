---
subcategory: "Cloud Connector Group"
layout: "zpa"
page_title: "ZPA: cloud_connector_group"
description: |-
  Gets a ZPA Cloud Connector Group's details.
---
# zpa_cloud_connector_group

The **zpa_cloud_connector_group** data source provides details about a specific cloud connector group created in the Zscaler Private Access cloud. This data source is required when creating resources such as.

1. Access Policy Rules where the Object Type = `CLOUD_CONNECTOR_GROUP` is being used.

## Example Usage

```hcl
# ZPA Cloud Connector Group Data Source
data "zpa_cloud_connector_group" "foo" {
  name = "AWS-Cloud"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This field defines the name of the cloud connector group.

## Attribute Reference

* `description` (String) This field defines the description of the cloud connector group.
* `enabled` (Boolean) This field defines the status of the cloud connector group.
* `geolocation_id` (Number)
* `zia_cloud` (String)
* `zia_org_id` (Number)
