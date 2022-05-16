---
layout: "zscaler"
page_title: "Zscaler Private Access (ZPA): machine_group"
sidebar_current: "docs-datasource-zpa-machine-group"
description: |-
  Get information about Machine Groups in Zscaler Private Access cloud.
---

# zpa_machine_group

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

* `creation_time` (Computed)
* `description` (Computed)
* `enabled` (Computed)
* `modified_by` (Computed)
* `modified_name` (Computed)
* `name` (Computed)
  * `machines` (Computed)
    * `creation_time` (Computed)
    * `description` (Computed)
    * `fingerprint` (Computed)
    * `id` (Computed)
    * `issued_cert_id` (Computed)
    * `machine_group_id` (Computed)
    * `machine_group_name` (Computed)
    * `machine_token_id` (Computed)
    * `modified_time` (Computed)
    * `modifiedby` (Computed)
    * `name` (Computed)
    * `signing_cert` (Computed)
