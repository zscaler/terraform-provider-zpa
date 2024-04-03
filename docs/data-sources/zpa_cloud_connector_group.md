---
page_title: "zpa_cloud_connector_group Data Source - terraform-provider-zpa"
subcategory: "Cloud Connector Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-web-server-certificates
  API documentation https://help.zscaler.com/zpa/obtaining-cloud-connector-group-details-using-api
  Get information about ZPA Cloud Connector Group in Zscaler Private Access cloud.
---

# zpa_cloud_connector_group (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-web-server-certificates)
* [API documentation](https://help.zscaler.com/zpa/obtaining-cloud-connector-group-details-using-api)

Use the **zpa_cloud_connector_group** data source to get information about a cloud connector group created from the Zscaler Private Access cloud. This data source can then be referenced within an Access Policy rule

~> **NOTE:** A Cloud Connector Group resource is created in the Zscaler Cloud Connector cloud and replicated to the ZPA cloud. This resource can then be referenced in a Access Policy Rule where the Object Type = `CLOUD_CONNECTOR_GROUP` is being used.

## Example Usage

```terraform
# ZPA Cloud Connector Group Data Source
data "zpa_cloud_connector_group" "foo" {
  name = "AWS-Cloud"
}
```

```terraform
# ZPA Cloud Connector Group Data Source
data "zpa_cloud_connector_group" "foo" {
  id = "1234567890"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (String) This field defines the name of the cloud connector group.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (string) This field defines the id of the cloud connector group.
* `description` (string) - This field defines the description of the cloud connector group.
* `enabled` (bool) - This field defines the status of the cloud connector group.
* `creation_time` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `modified_by` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `modified_time` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `geolocation_id` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `zia_cloud` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
* `zia_org_id` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.

* `cloud_connectors` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `creation_time` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `description` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `enabled` (bool) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `figerprint` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `ip_acl` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `issued_cert_id` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `modified_by` (string) - Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `modified_time`(string)- Only applicable for a GET request. Ignored in PUT/POST/DELETE requests.
  * `name` (string) - This field defines the name of the cloud connector group.
  * `enrollment_cert` (string) - This field defines the name of the cloud connector group.
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
