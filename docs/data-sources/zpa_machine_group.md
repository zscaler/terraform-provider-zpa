---
subcategory: "Machine Group"
layout: "zscaler"
page_title: "ZPA: machine_group"
description: |-
  Get information about Machine Groups in Zscaler Private Access cloud.
---

# Data Source: zpa_machine_group

Use the **zpa_machine_group** data source to get information about a machine group created in the Zscaler Private Access cloud. This data source can then be referenced in an Access Policy, Timeout policy, Forwarding Policy, Inspection Policy or Isolation Policy.

## Example Usage

```hcl
# ZPA Machine Group Data Source by name
data "zpa_machine_group" "example" {
  name = "MGR01"
}
```

```hcl
# ZPA Machine Group Data Source by id
data "zpa_machine_group" "example" {
  id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the machine group to be exported.
* `id` - (Optional) The ID of the machine group to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `creation_time` (string)
* `description` (string)
* `enabled` (bool)
* `modified_by` (string)
* `modified_name` (string)
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
