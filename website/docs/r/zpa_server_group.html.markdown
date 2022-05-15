---
subcategory: "Server Group "
layout: "zpa"
page_title: "ZPA: server_group"
description: |-
  Creates and manages ZPA Server Group resource
---

# zpa_server_group (Resource)

The **zpa_server_group** resource creates a server group in the Zscaler Private Access cloud. This resource can then be referenced in an application segment or application server resource.

## Example Usage

```hcl
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

### Required

* `name` - (Required) This field defines the name of the server group.
* `app_connector_groups` - (Required)
  * `id` - (Required) The ID of this resource.

## Attributes Reference

* `config_space*` (Optional)
* `description` (Optional) This field is the description of the server group.
* `dynamic_discovery` (Optional) This field controls dynamic discovery of the servers.
* `enabled` (Optional) This field defines if the server group is enabled or disabled.
* `servers` (Block List) This field is a list of servers that are applicable only when dynamic discovery is disabled. Server name is required only in cases where the new servers need to be created in this API.

## Import

Server Groups can be imported; use `<SERVER GROUP ID>` or `<SERVER GROUP NAME>` as the import ID.

For example:

```shell
terraform import zpa_server_group.example <server_group_id>
```

or

```shell
terraform import zpa_server_group.example <server_group_name>
```
