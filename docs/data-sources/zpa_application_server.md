---
subcategory: "Application Segment"
layout: "zpa"
page_title: "ZPA: application_server"
description: |-
  Gets a ZPA Application Server details.

---
# zpa_application_server (Data Source)

The **zpa_application_server** data source provides details about a specific application server created in the Zscaler Private Access cloud. This data source must be used in the following circumstances:

1. Server Group (When Dynamic Discovery is set to false)

## Example Usage

```hcl
# ZPA Application Server Data Source
data "zpa_application_server" "example" {
 name = "server.example.com"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This field defines the name of the server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

### Read-Only

* `address` (String) This field defines the domain or IP address of the server.
* `enabled` (Boolean) This field defines the status of the server.
* `app_server_group_ids` (Set of String) This field defines the list of server groups IDs.
* `description` (String) This field defines the description of the server.
* `config_space` (String)
