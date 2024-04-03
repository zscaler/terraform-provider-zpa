---
page_title: "zpa_server_group Data Source - terraform-provider-zpa"
subcategory: "Server Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-server-groups
  API documentation https://help.zscaler.com/zpa/configuring-server-groups-using-api
  Get information about Server Groups in Zscaler Private Access cloud.
---

# zpa_server_group (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-server-groups)
* [API documentation](https://help.zscaler.com/zpa/configuring-server-groups-using-api)

Use the **zpa_server_group** data source to get information about a server group created in the Zscaler Private Access cloud. This data source can then be referenced in an application segment, application server and Access Policy rule.

## Zenith Community - ZPA Server Groups

[![ZPA Terraform provider Video Series Ep4 - Server Groups](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_server_groups.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEmCAI/video-zpa-terraform-provider-video-series-ep4-server-groups)

## Example Usage

```terraform
# ZPA Server Group Data Source
data "zpa_server_group" "example" {
 name = "server_group_name"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the server group to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `id` - (Optional) The ID of the server group to be exported.
* `config_space` - (string)
* `description` - (string) This field is the description of the server group.
* `dynamic_discovery` - (bool) This field controls dynamic discovery of the servers.
* `enabled` - (bool) This field defines if the server group is enabled or disabled.
* `ip_anchored` - (bool)
* `app_connector_groups` (string)This field is a json array of app-connector-id only.
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.
