---
subcategory: "Cloud Connector Group"
layout: "zscaler"
page_title: "ZPA: cloud_connector_group"
description: |-
  Get information about ZPA Cloud Connector Group in Zscaler Private Access cloud.
---

# zpa_cloud_connector_group

Use the **zpa_cloud_connector_group** data source to get information about a cloud connector group created from the Zscaler Private Access cloud. This data source can then be referenced within an Access Policy rule

~> **NOTE:** A Cloud Connector Group resource is created in the Zscaler Cloud Connector cloud and replicated to the ZPA cloud. This resource can then be referenced in a Access Policy Rule where the Object Type = `CLOUD_CONNECTOR_GROUP` is being used.

## Example Usage

```hcl
# ZPA Cloud Connector Group Data Source
data "zpa_cloud_connector_group" "foo" {
  name = "AWS-Cloud"
}
```

```hcl
# ZPA Cloud Connector Group Data Source
data "zpa_cloud_connector_group" "foo" {
  id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This field defines the name of the cloud connector group.
* `id` - (Optional) This field defines the id of the cloud connector group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` (Computed) - This field defines the description of the cloud connector group.
* `enabled` (Computed) - This field defines the status of the cloud connector group.
* `creation_time` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `modified_by` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `modified_time` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `geolocation_id` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `zia_cloud` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `zia_org_id` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.

* `cloud_connectors` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `creation_time` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `description` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `enabled` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `figerprint` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `ip_acl` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `issued_cert_id` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `modified_by` (Computed) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `modified_time`(Computed)- Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `name` (Computed) - This field defines the name of the cloud connector group.
  * `enrollment_cert` (Computed) - This field defines the name of the cloud connector group.
