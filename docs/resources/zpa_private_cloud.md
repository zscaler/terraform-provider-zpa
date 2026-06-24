---
page_title: "zpa_private_cloud Resource - terraform-provider-zpa"
subcategory: "Private Clouds"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-private-cloud-controllers
  API documentation https://help.zscaler.com/zpa/about-private-cloud-controllers
  Creates and manages ZPA Private Cloud in Zscaler Private Access cloud.
---

# zpa_private_cloud (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-private-cloud-controllers)
* [API documentation](https://help.zscaler.com/zpa/about-private-cloud-controllers)

The **zpa_private_cloud** resource creates and manages a private cloud in the Zscaler Private Access cloud. This resource can then be associated with the related Private Cloud component groups (Site Controller, SIEM, Private Exporter, and Private Broker groups).

## Example Usage - Private Cloud Group Controller Association

```terraform
resource "zpa_private_cloud" "this" {
  name               = "PrivateCloud01"
  description        = "Example private cloud"
  enabled            = true
  re_enroll_period   = "90"
  fire_drill_enabled = false
  sitec_preferred    = false
  remote_lss         = false

  site_controller_group_ids {
    id = [zpa_private_cloud_group.this.id]
  }
}

resource "zpa_private_cloud_group" "this" {
  name                     = "PrivateCloudGroup01"
  description              = "Example private cloud group"
  enabled                  = true
  country_code             = "US"
  city_country             = "San Jose, US"
  latitude                 = "37.33874"
  longitude                = "-121.8852525"
  location                 = "San Jose, CA, USA"
  upgrade_day              = "SUNDAY"
  upgrade_time_in_secs     = "66600"
  version_profile_id       = "0"
  override_version_profile = true
  is_public                = "TRUE"
}
```

## Example Usage - App Connector Group Controller Association

```terraform
resource "zpa_private_cloud" "this" {
  name               = "PrivateCloud01"
  description        = "Example private cloud"
  enabled            = true
  re_enroll_period   = "90"
  fire_drill_enabled = false
  sitec_preferred    = false
  remote_lss         = false

  assistant_groups_ids {
    id = [zpa_app_connector_group.this.id]
  }
}

resource "zpa_app_connector_group" "this" {
  name                          = "AppConnectorGroup01"
  description                   = "AppConnectorGroup01"
  enabled                       = true
  city_country                  = "San Jose, US"
  country_code                  = "US"
  latitude                      = "37.338"
  longitude                     = "-121.8863"
  location                      = "San Jose, CA, US"
  upgrade_day                   = "SUNDAY"
  upgrade_time_in_secs          = "66600"
  version_profile_id            = "0"
  override_version_profile      = true
  dns_query_type                = "IPV4_IPV6"
  use_in_dr_mode                = false
}
```

## Example Usage - Service Edge Group Controller Association

```terraform
resource "zpa_private_cloud" "this" {
  name               = "PrivateCloud01"
  description        = "Example private cloud"
  enabled            = true
  re_enroll_period   = "90"
  fire_drill_enabled = false
  sitec_preferred    = false
  remote_lss         = false

  private_broker_group_ids {
    id = [zpa_service_edge_group.this.id]
  }
}

resource "zpa_service_edge_group" "this" {
  name                        = "ServiceEdgeGroup01"
  description                 = "ServiceEdgeGroup01"
  enabled                     = true
  is_public                   = false
  upgrade_day                 = "SUNDAY"
  city_country                = "San Jose, US"
  country_code                = "US"
  latitude                    = "37.338"
  longitude                   = "-121.8863"
  location                    = "San Jose, CA, US"
  upgrade_time_in_secs        = "66600"
  version_profile_id            = "0"
  override_version_profile      = true
}
```

## Example Usage - With Fire Drill Site

```terraform
resource "zpa_private_cloud" "this" {
  name               = "PrivateCloud01"
  description        = "Example private cloud"
  enabled            = true
  re_enroll_period   = "86400"
  fire_drill_enabled = true
  sitec_preferred    = true

  zpn_fire_drill_site {
    fire_drill_interval           = "7"
    fire_drill_interval_time_unit = "DAYS"
  }
  site_controller_group_ids {
    id = [zpa_private_cloud_group.this.id]
  }
}

resource "zpa_private_cloud_group" "this" {
  name                     = "PrivateCloudGroup01"
  description              = "Example private cloud group"
  enabled                  = true
  country_code             = "US"
  city_country             = "San Jose, US"
  latitude                 = "37.33874"
  longitude                = "-121.8852525"
  location                 = "San Jose, CA, USA"
  upgrade_day              = "SUNDAY"
  upgrade_time_in_secs     = "66600"
  version_profile_id       = "0"
  override_version_profile = true
  is_public                = "TRUE"
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (String) Name of the Private Cloud

### Optional

- `id` - (String) The ID of the Private Cloud
- `description` - (String) Description of the Private Cloud
- `enabled` - (Boolean) Whether this Private Cloud is enabled or not. Supported values: `true`, `false`
- `re_enroll_period` - (String) Specify the number of days before the certificate expires to automatically renew the enrollment certificate for Private Cloud Controllers, App Connectors, and ZPA Private Service Edges. Supported Values between `30` and `180` days
- `fire_drill_enabled` - (Boolean) Whether fire drill is enabled for the Private Cloud. Supported values: `true`, `false`
- `sitec_preferred` - (Boolean) By default, App Connectors and ZPA Private Service Edges establish a control channel to the Zscaler Zero Trust Exchange (ZTE) during normal operations. Select Private Cloud Controller to allow App Connectors and ZPA Private Service Edges to use Private Cloud Controllers for control channels like configuration downloads, configuration updates, and logging even during normal operations and when not in Business Continuity. Select Public Cloud to use the ZTE. Supported values: `true`, `false`
- `remote_lss` - (Boolean) Enable to allow Logging through LSS App Connectors. Supported values: `true`, `false`
- `microtenant_id` - (String) Microtenant ID for the Private Cloud

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

- `assistant_groups_ids` - (List) The App Connector Groups associated with the Private Cloud.
    - `id` - (Set of String) The set of App Connector Group IDs.

- `site_controller_group_ids` - (List) The Private Cloud Group associated with the Private Cloud.
    - `id` - (Set of String) The set of Private Cloud Group IDs.

- `siem_ids` - (List) The LSS App Connector Group associated with the Private Cloud.
    - `id` - (Set of String) The set of App Connector Groups IDs

- `private_broker_group_ids` - (List)  The Service Edge Groups associated with the Private Cloud.
    - `id` - (Set of String) The set of Service Edge Groups IDs.

- `zpn_fire_drill_site` - (Block List) The fire drill site configuration for the Private Cloud. This block requires the attribute `fire_drill_enabled` to be set to true.
    - `id` - (String) The unique identifier of the fire drill site.
    - `microtenant_id` - (String) The microtenant ID of the fire drill site.
    - `fire_drill_interval` - (String) The fire drill interval.
    - `fire_drill_interval_time_unit` - (String) The fire drill interval time unit. Supported values: `SECONDS`, `MINUTES`, `HOURS`, `DAYS`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Private Cloud can be imported by using `<PRIVATE CLOUD ID>` or `<PRIVATE CLOUD NAME>` as the import ID.

```shell
terraform import zpa_private_cloud.example <private_cloud_id>
```

or

```shell
terraform import zpa_private_cloud.example <private_cloud_name>
```
