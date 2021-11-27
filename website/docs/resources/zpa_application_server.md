---
subcategory: "Application Server"
layout: "zpa"
page_title: "ZPA: application_server"
description: |-
  Creates a ZPA Application Server.
---

# zpa_application_server (Resource)

The **zpa_application_server** resource creates an application server in the Zscaler Private Access cloud. This resource can then be referenced in a server group.

## Example Usage

```hcl
# ZPA Application Server resource
resource "zpa_application_server" "server1" {
  name                          = "Example"
  description                   = "Example"
  address                       = "192.168.1.1"
  enabled                       = true
}
```

```hcl
# ZPA Application Server resource
resource "zpa_application_server" "server1" {
  name                          = "Example"
  description                   = "Example"
  address                       = "192.168.1.1"
  enabled                       = true
  app_server_group_ids          = [data.zpa_server_group.example.com]
}

data "zpa_server_group" "example" {
    name = "Example"
} 
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name. The name of the application server to be exported.
* `address` - (Required) Address. The address of the application server to be exported.

## Attributes Reference

* `app_server_group_ids` - (Optional) This field defines the list of server group IDs.
* `description` - (Optional) This field defines the description of the server.
* `enabled` - (Optional) This field defines the status of the server.
* `config_space` - (Optional)

## Import

Application Server can be imported by using `<APPLICATION SERVER ID>` as the import ID.

For example:

```shell
terraform import zpa_application_server.example <application_server_id>
```
