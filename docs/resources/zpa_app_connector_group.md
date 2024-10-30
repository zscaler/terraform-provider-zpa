---
page_title: "zpa_app_connector_group Resource - terraform-provider-zpa"
subcategory: "App Connector Group"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-connector-groups
  API documentation https://help.zscaler.com/zpa/configuring-app-connector-groups-using-api
  Creates and manages ZPA App Connector Groups
---

# zpa_app_connector_group (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-connector-groups)
* [API documentation](https://help.zscaler.com/zpa/configuring-app-connector-groups-using-api)

The **zpa_app_connector_group** resource creates a and manages app connector groups in the Zscaler Private Access (ZPA) cloud. This resource can then be associated with the following resources: server groups, log receivers and access policies.

## Zenith Community - ZPA App Connector Group

[![ZPA Terraform provider Video Series Ep2 - Connector Groups](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_app_connector_group.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEoCAI/video-zpa-terraform-provider-video-series-ep2-connector-groups)

## Example Usage

```terraform
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
  version_profile_name          = "New Release"
  dns_query_type                = "IPV4_IPV6"
  use_in_dr_mode                = true
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (String) Name of the App Connector Group.
- `enabled` - (Boolean) Whether this App Connector Group is enabled or not. Default value: `true`. Supported values: `true`, `false`
- `latitude` - (String) Latitude of the App Connector Group. Integer or decimal. With values in the range of `-90` to `90`
- `longitude` - (String) Longitude of the App Connector Group. Integer or decimal. With values in the range of `-180` to `180`
- `location` - (String) Location of the App Connector Group. i.e ``"San Jose, CA, USA"``
- `city_country` - (String) Whether Double Encryption is enabled or disabled for the app. i.e ``"San Jose, US"``
- `country_code` - (String) Provide a 2 letter [ISO3166 Alpha2 Country code](https://en.wikipedia.org/wiki/List_of_ISO_3166_country_codes). i.e ``"US"``, ``"CA"``

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

### Optional

- `description` (String) Description of the App Connector Group.
- `upgrade_day` - (String) App Connectors in this group will attempt to update to a newer version of the software during this specified day i.e ``SUNDAY``
- `upgrade_time_in_secs` - (String) App Connectors in this group will attempt to update to a newer version of the software during this specified time. Default value: `66600`. Integer in seconds (i.e., `-66600`). The integer should be greater than or equal to `0` and less than `86400`, in `15` minute intervals
- `override_version_profile` - (Boolean) Whether the default version profile of the App Connector Group is applied or overridden. Default: `false` Supported values: `true`, `false`
- `version_profile_id` - (String) The unique identifier of the version profile. Supported values are:
  - ``0`` = ``Default``
  - ``1`` = ``Previous Default``
  - ``2`` = ``New Release``
- `dns_query_type` - (String) Whether IPv4, IPv6, or both, are enabled for DNS resolution of all applications in the App Connector Group. Supported values are: ``IPV4``, ``IPV6``, or ``IPV4_IPV6``
- `tcp_quick_ack_app` - (Boolean) Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value. Supported values: `true`, `false`
- `tcp_quick_ack_assistant` - (Boolean) Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value. Supported values: `true`, `false`
- `tcp_quick_ack_read_assistant` - (Boolean) Whether TCP Quick Acknowledgement is enabled or disabled for the application. The tcpQuickAckApp, tcpQuickAckAssistant, and tcpQuickAckReadAssistant fields must all share the same value. Supported values: `true`, `false`
- `use_in_dr_mode` - (Boolean) Whether or not the App Connector Group is designated for disaster recovery. Supported values: `true`, `false`
- `pra_enabled` - (Boolean) Whether or not Privileged Remote Access is enabled on the App Connector Group. Supported values: `true`, `false`
- `waf_disabled` - (Boolean) Whether or not AppProtection is disabled for the App Connector Group. Supported values: `true`, `false`
- `microtenant_id` (String) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as `0` when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.
- `lss_app_connector_group` (boolean) Whether or not the App Connector Group is configured for the Log Streaming Service (LSS).

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

App Connector Group can be imported by using `<APP CONNECTOR GROUP ID>` or `<APP CONNECTOR GROUP NAME>`as the import ID.

```shell
terraform import zpa_app_connector_group.example <app_connector_group_id>
```

or

```shell
terraform import zpa_app_connector_group.example <app_connector_group_name>
```
