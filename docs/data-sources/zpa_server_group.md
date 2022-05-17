---
subcategory: "Server Group"
layout: "zscaler"
page_title: "ZPA: server_group"
description: |-
  Get information about Server Groups in Zscaler Private Access cloud.
---

# zpa_server_group

Use the **zpa_server_group** data source to get information about a server group created in the Zscaler Private Access cloud. This data source can then be referenced in an application segment, application server and Access Policy rule.

## Example Usage

```hcl
# ZPA Server Group Data Source
data "zpa_server_group" "example" {
 name = "server_group_name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the server group to be exported.
* `id` - (Optional) The ID of the server group to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `config_space` - (Computed)
* `description` - (Computed) This field is the description of the server group.
* `dynamic_discovery` - (Computed) This field controls dynamic discovery of the servers.
* `enabled` - (Computed) This field defines if the server group is enabled or disabled.
* `ip_anchored` - (Computed)
* `app_connector_groups` (Computed)This field is a json array of app-connector-id only.
