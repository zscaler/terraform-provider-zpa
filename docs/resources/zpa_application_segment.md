---
page_title: "zpa_application_segment Resource - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-applications/API documentation https://help.zscaler.com/zpa/configuring-application-segments-using-api
  Creates and manages ZPA Application Segments
---

# zpa_application_segment (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-applications)
* [API documentation](https://help.zscaler.com/zpa/configuring-application-segments-using-api)

The **zpa_application_segment** resource creates an application segment in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule or access policy client forwarding rule.

## Zenith Community - ZPA Application Segment

[![ZPA Terraform provider Video Series Ep7 - Application Segment](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_application_segments.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEXCAY/video-zpa-terraform-provider-video-series-ep7-zpa-application-segment)

## Example 1 Usage

```terraform
# ZPA Application Segment resource
resource "zpa_application_segment" "this" {
    name              = "Example"
    description       = "Example"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
    tcp_port_ranges   = ["8080", "8080"]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.this.id
    server_groups {
        id = [ zpa_server_group.this.id]
    }
    depends_on = [ zpa_server_group.this, zpa_segment_group.this]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "this" {
  name            = "Example"
  description     = "Example"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "this" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.this.id ]
  }
  depends_on = [ zpa_app_connector_group.this ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "this" {
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

## Example 2 Usage

```terraform
# ZPA Application Segment resource
resource "zpa_application_segment" "this" {
    name              = "Example"
    description       = "Example"
    enabled           = true
    health_reporting  = "ON_ACCESS"
    bypass_type       = "NEVER"
    is_cname_enabled  = true
  tcp_port_range = [
    {
    from = "8080"
    to = "8080"
    }
  ]
  udp_port_range = [
    {
    from = "8080"
    to = "8080"
    }
  ]
    domain_names      = ["server.acme.com"]
    segment_group_id  = zpa_segment_group.this.id
    server_groups {
        id = [ zpa_server_group.this.id]
    }
    depends_on = [ zpa_server_group.this, zpa_segment_group.this]
}

# ZPA Segment Group resource
resource "zpa_segment_group" "this" {
  name            = "Example"
  description     = "Example"
  enabled         = true
}

# ZPA Server Group resource
resource "zpa_server_group" "this" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = false
  app_connector_groups {
    id = [ zpa_app_connector_group.this.id ]
  }
  depends_on = [ zpa_app_connector_group.this ]
}

# ZPA App Connector Group resource
resource "zpa_app_connector_group" "this" {
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
  dns_query_type                = "IPV4_IPV6"
}
```

## Example 3 Usage - Application Segment Extranet Configuration

```terraform
data "zpa_location_controller" "this" {
  name        = "ExtranetLocation01 | zscalerbeta.net"
  zia_er_name = "NewExtranet 8432"
}

data "zpa_location_group_controller" "this" {
  location_name = "ExtranetLocation01"
  zia_er_name   = "NewExtranet 8432"
}

data "zpa_extranet_resource_partner" "this" {
  name = "NewExtranet 8432"
}

resource "zpa_segment_group" "this" {
  name                   = "Example"
  description            = "Example"
  enabled                = true
}

resource "zpa_server_group" "this" {
  name              = "Example"
  description       = "Example"
  enabled           = true
  dynamic_discovery = true
  extranet_enabled  = true

  extranet_dto {
    zpn_er_id = data.zpa_extranet_resource_partner.this.id

    location_dto {
      id = data.zpa_location_controller.this.id
    }

    location_group_dto {
      id = data.zpa_location_group_controller.this.id
    }
  }
}

resource "zpa_application_segment" "this" {
  name              = "app01.acme.com"
  description       = "app01.acme.com"
  enabled           = true
  health_reporting  = "NONE"
  health_check_type = "NONE"
  bypass_type       = "NEVER"
  is_cname_enabled  = true
  tcp_port_ranges   = ["8080", "8080"]
  domain_names      = ["app01.acme.com"]
  segment_group_id  = zpa_segment_group.this.id
  server_groups {
    id = [zpa_server_group.this.id]
  }
  zpn_er_id {
    id = [data.zpa_extranet_resource_partner.this.id]
  }
}
```

## Schema

### Required

The following arguments are supported:

- `name` - (String) Name. The name of the App Connector Group to be exported.
- `domain_names` - (List) List of domains and IPs.
- `server_groups` - (Block Set) List of Server Group IDs
  - `id` - (Required)
- `segment_group_id` - (string) The unique identifier of the segment group.
- `tcp_port_ranges` - (List of String) TCP port ranges used to access the app.
- `udp_port_ranges` - (List of String) UDP port ranges used to access the app.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model below. We recommend this model as opposed of the legacy model via `tcp_port_ranges` and or `udp_port_ranges`.

- `tcp_port_range` - (Block Set) TCP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.

- `udp_port_range` (Block Set) UDP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.

-> **NOTE:** Application segments must have unique ports and cannot have overlapping domain names using the same tcp/udp ports across multiple application segments.

### Optional

- `description` (String) Description of the application.
- `bypass_type` (String) Indicates whether users can bypass ZPA to access applications. Supported values: `ALWAYS`, `NEVER`, `ON_NET`.
- `bypass_on_reauth` (Boolean) Supported values: `true`, `false`
- `double_encrypt` (Boolean) Whether Double Encryption is enabled or disabled for the app.
- `enabled` (Boolean) Whether this application is enabled or not.
- `health_reporting` (String) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
- `health_check_type` (String) Whether the health check is enabled (DEFAULT) or disabled (NONE) for the application. Supported values: `DEFAULT`, `NONE`.
- `icmp_access_type` - (String) The ICMP access type. Supported values: `PING_TRACEROUTING`, `PING`, `NONE`
- `ip_anchored` (Boolean) Supported values: `true`, `false`
- `is_cname_enabled` (Boolean) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors. Supported values: `true`, `false`
- `inspect_traffic_with_zia` (Boolean) Indicates if Inspect Traffic with ZIA is enabled for the application. When enabled, this leverages a single posture for securing internet/SaaS and private applications, and applies Data Loss Prevention policies to the application segment you are creating.Supported values: `true`, `false`
- `match_style` (String) Indicates if Multimatch is enabled for the application segment. If enabled (INCLUSIVE), the request allows traffic to match multiple applications. If disabled (`EXCLUSIVE`), the request allows traffic to match a single application. A domain can only be INCLUSIVE or EXCLUSIVE, and any application segment can only contain inclusive or exclusive domains.
Supported values: `EXCLUSIVE`, `INCLUSIVE`. [Learn More](https://help.zscaler.com/zpa/using-app-segment-multimatch)
- `tcp_keep_alive` (String) Whether the application is using TCP communication sockets or not. Supported values: ``1`` for Enabled and ``0`` for Disabled
- `passive_health_enabled` - Indicates if passive health checks are enabled on the application. (Boolean) Supported values: `true`, `false`

- `select_connector_close_to_app` (Boolean) Whether the App Connector is closest to the application (true) or closest to the user (false). Supported values: `true`, `false`

- `use_in_dr_mode` - (Boolean) Supported values: `true`, `false`
- `is_incomplete_dr_config` - (Boolean) Supported values: `true`, `false`
- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.
- `fqdn_dns_check` - (Boolean) When set to Enabled, Zscaler Client Connector receives CNAME DNS records from the App Connector for FQDN applications. Supported values: `true`, `false`
- `share_to_microtenants` (List) List of destination Microtenants to which the application segment is to be shared with.
- `zpn_er_id` (Block Set) - ZPN Extranet Resource
    - `id` - (String) The unique identifier of the zpn extranet resource

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/SecurityGeekIO/zscaler-terraformer)

Application Segment can be imported by using `<APPLICATION SEGMENT ID>` or `<APPLICATION SEGMENT NAME>` as the import ID.

```shell
terraform import zpa_application_segment.example <application_segment_id>
```

or

```shell
terraform import zpa_application_segment.example <application_segment_name>
```
