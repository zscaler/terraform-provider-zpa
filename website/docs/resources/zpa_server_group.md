---
subcategory: "Server Group "
layout: "zpa"
page_title: "ZPA: server_group"
description: |-
  Creates a ZPA Server Group resource
---

# zpa_server_group (Resource)

The **zpa_server_group** resource creates a server group in the Zscaler Private Access cloud. This resource can then be referenced in an application segment or application server resource.

## Example Usage

```hcl
# ZPA Server Group resource with Dynamic Discovery Enabled
resource "zpa_server_group" "example" {
  name = "Example"
  description = "Example"
  enabled = true
  dynamic_discovery = true
  app_connector_groups {
    id = [data.zpa_app_connector_group.aws_connector_group.id]
  }
}

data "zpa_app_connector_group" "aws_connector_group" {
  name = "AWS-Connector-Group"
}
```

```hcl
# ZPA Server Group resource with Dynamic Discovery Disabled
resource "zpa_server_group" "example" {
  name = "Example"
  description = "Example"
  enabled = true
  dynamic_discovery = false
servers {
    id = [zpa_application_server.example.id]
  }
  app_connector_groups {
    id = [data.zpa_app_connector_group.aws_connector_group.id]
  }
}

data "zpa_app_connector_group" "aws_connector_group" {
  name = "AWS-Connector-Group"
}

resource "zpa_application_server" "example" {
  name                          = "Example"
  description                   = "Example"
  address                       = "server.example.com"
  enabled                       = true
}
```

### Required

* `name` - (Required) This field defines the name of the server group.

`app_connector_groups` - (Required)

* `id` - (Required) The ID of this resource.

## Attributes Reference

* `config_space*` (String)
* `description` (String) This field is the description of the server group.
* `dynamic_discovery` (Boolean) This field controls dynamic discovery of the servers.
* `enabled** (Boolean) This field defines if the server group is enabled or disabled.
* `ip_anchored` (Boolean)
* `servers` (Block List) This field is a list of servers that are applicable only when dynamic discovery is disabled. Server name is required only in cases where the new servers need to be created in this API.

## Import

Server Groups can be imported; use `<SERVER GROUP ID>` as the import ID.

For example:

```shell
terraform import zpa_server_group.example 216196257331290863
```
