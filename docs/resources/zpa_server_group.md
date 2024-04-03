---
page_title: "zpa_server_group Resource - terraform-provider-zpa"
subcategory: "Server Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-server-groups
  API documentation https://help.zscaler.com/zpa/configuring-server-groups-using-api
  Creates and manages ZPA Server Group resource
---

# zpa_server_group (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-server-groups)
* [API documentation](https://help.zscaler.com/zpa/configuring-server-groups-using-api)

The **zpa_server_group** resource creates a server group in the Zscaler Private Access cloud. This resource can then be referenced in an application segment or application server resource.

## Zenith Community - ZPA Server Groups

[![ZPA Terraform provider Video Series Ep4 - Server Groups](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_server_groups.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEmCAI/video-zpa-terraform-provider-video-series-ep4-server-groups)

## Example Usage

```terraform
# Create a Server Group resource with Dynamic Discovery Enabled
resource "zpa_server_group" "example" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = true
  app_connector_groups {
    id = [ zpa_app_connector_group.example.id ]
  }
  depends_on = [ zpa_app_connector_group.example ]
}

# Create a App Connector Group
resource "zpa_app_connector_group" "example" {
  name                          = "Example"
  description                   = "Example"
  enabled                       = true
  city_country                  = "San Jose, CA"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}
```

```terraform
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
    id = [ zpa_app_connector_group.example.id ]
  }
  depends_on = [ zpa_app_connector_group.example, zpa_application_server.server ]
}

# Create an application server
resource "zpa_application_server" "example" {
  name                          = "Example"
  description                   = "Example"
  address                       = "server.example.com"
  enabled                       = true
}

# Create a App Connector Group
resource "zpa_app_connector_group" "example" {
  name                          = "Example"
  description                   = "Example"
  enabled                       = true
  city_country                  = "San Jose, CA"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  override_version_profile      = true
  version_profile_id            = 0
  dns_query_type                = "IPV4"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (String) This field defines the name of the server group.
* `app_connector_groups` - (Required)
  * `id` - (Required) The ID of this resource.

### Optional

In addition to all arguments above, the following attributes are exported:

* `description` (String) This field is the description of the server group.
* `dynamic_discovery` (String) This field controls dynamic discovery of the servers. Supported values are `true` and `false`
* `enabled` (String) This field defines if the server group is enabled or disabled.
* `servers` (Block List) This field is a list of application servers that are applicable only when dynamic discovery is disabled `false`.
* `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Server Groups can be imported; use `<SERVER GROUP ID>` or `<SERVER GROUP NAME>` as the import ID.

For example:

```shell
terraform import zpa_server_group.example <server_group_id>
```

or

```shell
terraform import zpa_server_group.example <server_group_name>
```
