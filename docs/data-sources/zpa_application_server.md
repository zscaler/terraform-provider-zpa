---
page_title: "zpa_application_server Data Source - terraform-provider-zpa"
subcategory: "Application Server"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-servers
  API documentation https://help.zscaler.com/zpa/configuring-servers-using-api
  Get information about ZPA Application Server in Zscaler Private Access cloud.
---

# zpa_application_server (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-servers)
* [API documentation](https://help.zscaler.com/zpa/configuring-servers-using-api)

Use the **zpa_application_server** data source to get information about an application server created in the Zscaler Private Access cloud. This data source must be used in the following circumstances:

**NOTE:** To ensure consistent search results across data sources, please avoid using multiple spaces or special characters in your search queries.

1. Server Group (When Dynamic Discovery is set to false)

## Zenith Community - ZPA Application Server

[![ZPA Terraform provider Video Series Ep5 - Application Server](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_application_servers.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEgCAI/video-terraform-provider-video-series-ep5-zpa-application-server)

## Example Usage

```terraform
# ZPA Application Server Data Source by Name
data "zpa_application_server" "example" {
 name = "server.example.com"
}
```

```terraform
# ZPA Application Server Data Source by ID
data "zpa_application_server" "example" {
 id = "1234567890"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) This field defines the name of the server.
* `id` - (Optional) This field defines the id of the application server.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `description` - (string) This field defines the description of the server.
* `address` - (string) This field defines the domain or IP address of the server.
* `enabled` - (bool) This field defines the status of the server.
* `app_server_group_ids` - (Set of String) This field defines the list of server groups IDs.
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
