---
subcategory: "Application Server"
layout: "zscaler"
page_title: "ZPA: application_server"
description: |-
  Get information about ZPA Application Server in Zscaler Private Access cloud.
---

# Data Source: zpa_application_server

Use the **zpa_application_server** data source to get information about an application server created in the Zscaler Private Access cloud. This data source must be used in the following circumstances:

1. Server Group (When Dynamic Discovery is set to false)

## Zenith Community - ZPA Application Server

[![ZPA Terraform provider Video Series Ep5 - Application Server](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_application_servers.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEgCAI/video-terraform-provider-video-series-ep5-zpa-application-server)

## Example Usage

```hcl
# ZPA Application Server Data Source by Name
data "zpa_application_server" "example" {
 name = "server.example.com"
}
```

```hcl
# ZPA Application Server Data Source by ID
data "zpa_application_server" "example" {
 id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) This field defines the name of the server.
* `id` - (Optional) This field defines the id of the application server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) This field defines the description of the server.
* `address` - (string) This field defines the domain or IP address of the server.
* `enabled` - (bool) This field defines the status of the server.
* `app_server_group_ids` - (Set of String) This field defines the list of server groups IDs.
