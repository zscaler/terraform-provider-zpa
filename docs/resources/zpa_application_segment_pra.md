---
page_title: "zpa_application_segment_pra Resource - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-privileged-remote-access-applications
  API documentation https://help.zscaler.com/zpa/configuring-application-segments-using-api
  Creates and manages ZPA Application Segment for Privileged Remote Access.
---

# zpa_application_segment_pra (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-privileged-remote-access-applications)
* [API documentation](https://help.zscaler.com/zpa/configuring-application-segments-using-api)

The **zpa_application_segment_pra** resource creates an application segment for Privileged Remote Access in the Zscaler Private Access cloud. This resource can then be referenced in an access policy rule, access policy timeout rule, access policy client forwarding rule and inspection policy. This resource supports Privileged Remote Access for both `RDP` and `SSH`.

## Example Usage

```terraform
resource "zpa_application_segment_pra" "this" {
  name             = "PRA_Example"
  description      = "PRA_Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["22", "22", "3389", "3389"]
  domain_names     = ["ssh_pra.example.com", "rdp_pra.example.com"]
  segment_group_id = zpa_segment_group.this.id
  server_groups {
    id = [zpa_server_group.this.id]
  }
  common_apps_dto {
    apps_config {
      name                 = "ssh_pra"
      domain               = "ssh_pra.example.com"
      application_protocol = "SSH"
      application_port     = "22"
      enabled = true
      app_types = [ "SECURE_REMOTE_ACCESS" ]
    }
    apps_config {
    name                 = "rdp_pra"
    domain               = "rdp_pra.example.com"
    application_protocol = "RDP"
    connection_security  = "ANY"
    application_port     = "3389"
    enabled = true
    app_types = [ "SECURE_REMOTE_ACCESS" ]
    }
  }
}
```

## Schema

### Required

The following arguments are supported:

- `name` (String) Name. The name of the App Connector Group to be exported.
- `domain_names` - (Required) List of domains and IPs.
- `server_groups` - (Required) List of Server Group IDs
- `segment_group_id` - (String) The unique identifier of the Segment Group.
- `common_apps_dto` (Block Set, Min: 1) List of applications (e.g., Inspection, Browser Access or Privileged Remote Access)
  - `apps_config:` (Block Set, Min: 1) List of applications to be configured
    - `name` - (String) Name of the Privileged Remote Access
    - `domain` - (String) Domain name of the Privileged Remote Access
    - `application_protocol` - (String) Protocol for the Privileged Remote Access. Supported values: `RDP` and `SSH`
    - `application_port` - (String) Port for the Privileged Remote Access
    - `app_types` (List of String) Indicates the type of application as Privileged Remote Access. Supported value: `SECURE_REMOTE_ACCESS`
    - `connection_security` - (String) - Parameter required when `application_protocol` is of type `RDP`
    - `enabled` - (Boolean) Whether this application is enabled or not
- `tcp_port_ranges` - (List of String) TCP port ranges used to access the app.
- `udp_port_ranges` - (List of String) UDP port ranges used to access the app.

!> **WARNING:** Removing PRA applications from the `common_apps_dto.apps_config` block will cause the provider to force a replacement of the application segment.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model:

- `tcp_port_range` - (Block Set) TCP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.

- `udp_port_range` - (Block Set) UDP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.


-> **NOTE:** Application segments must have unique ports and cannot have overlapping domain names using the same tcp/udp ports across multiple application segments.

### Optional

- `description` - (String) Description of the application.
- `bypass_type` (String) Indicates whether users can bypass ZPA to access applications. Default value is: `NEVER` and supported values are: `ALWAYS`, `NEVER` and `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
- `double_encrypt` (Boolean) Whether Double Encryption is enabled or disabled for the app. Supported values are `true` and `false`
- `enabled` - (Boolean) Whether this application is enabled or not. Supported values are `true` and `false`
- `health_reporting` (String) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
- `health_check_type` (String) Default: `DEFAULT`. Supported values: `DEFAULT`, `NONE`
- `icmp_access_type` - (String) The ICMP access type. Supported values: `PING_TRACEROUTING`, `PING`, `NONE`
- `ip_anchored` - (Boolean) Whether Source IP Anchoring for use with ZIA is enabled or disabled for the application.
- `is_cname_enabled` (Boolean) Indicates if the Zscaler Client Connector (formerly Zscaler App or Z App) receives CNAME DNS records from the connectors. Supported values: `true`, `false`
- `tcp_keep_alive` (String) Whether the application is using TCP communication sockets or not. Supported values: ``1`` for Enabled and ``0`` for Disabled
- `passive_health_enabled` - Indicates if passive health checks are enabled on the application. (Boolean) Supported values: `true`, `false`

- `select_connector_close_to_app` (Boolean) Whether the App Connector is closest to the application (true) or closest to the user (false). Supported values: `true`, `false`

- `use_in_dr_mode` - (Boolean) Supported values: `true`, `false`
- `is_incomplete_dr_config` - (Boolean) Supported values: `true`, `false`
- `microtenant_id` (String) The ID of the microtenant the resource is to be associated with.

⚠️ **WARNING:**: The attribute ``microtenant_id`` is optional and requires the microtenant license and feature flag enabled for the respective tenant. The provider also supports the microtenant ID configuration via the environment variable `ZPA_MICROTENANT_ID` which is the recommended method.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZPA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Application Segment can be imported by using `<APPLICATION SEGMENT ID>` or `<APPLICATION SEGMENT NAME>` as the import ID.

```shell
terraform import zpa_application_segment_pra.example <application_segment_id>
```

or

```shell
terraform import zpa_application_segment_pra.example <application_segment_name>
```
