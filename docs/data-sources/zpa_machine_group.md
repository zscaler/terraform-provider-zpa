---
page_title: "zpa_machine_group Data Source - terraform-provider-zpa"
subcategory: "Machine Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-machine-groups
  API documentation https://help.zscaler.com/zpa/obtaining-machine-group-details-using-api
  Get information about Machine Groups in Zscaler Private Access cloud.
---

# zpa_machine_group (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-machine-groups)
* [API documentation](https://help.zscaler.com/zpa/obtaining-machine-group-details-using-api)

Use the **zpa_machine_group** data source to get information about a machine group created in the Zscaler Private Access cloud. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

## Example Usage

```terraform
# ZPA Machine Group Data Source by name
data "zpa_machine_group" "example" {
  name = "MGR01"
}
```

```terraform
# ZPA Machine Group Data Source by id
data "zpa_machine_group" "example" {
  id = "1234567890"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the machine group to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (string) The ID of the machine group to be exported.
* `creation_time` (string)
* `description` (string)
* `enabled` (bool)
* `modified_by` (string)
* `modified_name` (string)
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
  * `machines` (string)
    * `creation_time` (string)
    * `description` (string)
    * `fingerprint` (string)
    * `id` (string)
    * `issued_cert_id` (string)
    * `machine_group_id` (string)
    * `machine_group_name` (string)
    * `machine_token_id` (string)
    * `modified_time` (string)
    * `modified_by` (string)
    * `name` (string)
    * `signing_cert` (string)
    * `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
    * `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
