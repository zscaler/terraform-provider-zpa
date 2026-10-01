---
page_title: "zpa_application_segment_inspection Resource - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-appprotection-applications
  API documentation https://help.zscaler.com/zpa/configuring-application-segments-using-api
  Creates and manages ZPA Application Segment for Inspection
---

# zpa_application_segment_inspection (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-appprotection-applications)
* [API documentation](https://help.zscaler.com/zpa/configuring-application-segments-using-api)

The **zpa_application_segment_inspection** resource creates an inspection application segment in the Zscaler Private Access cloud. This resource can then be referenced in an access policy inspection rule. This resource supports Inspection for both `HTTP` and `HTTPS`.

~> **NOTE: This is the only resource that supports AppProtection and Active Directory Inspection.** The segment-level **Auto App Protection** (`auto_app_protect_enabled`) and **Active Directory Inspection** (`adp_enabled`) settings, and the per-application **Use Untrusted Certificates** option (`trust_untrusted_cert`), are supported exclusively here. Enabling Auto App Protection or Active Directory Inspection on a segment in the ZPA Admin Portal converts it into an inspection application segment: if that segment is currently managed by [`zpa_application_segment`](zpa_application_segment.md), import it into this resource and remove it from `zpa_application_segment`. Continuing to manage a converted segment with `zpa_application_segment` is not supported — that resource has no inspection attributes, so any update it applies disables these settings.

## Example Usage

```terraform
data "zpa_ba_certificate" "jenkins" {
  name = "jenkins.example.com"
}

resource "zpa_application_segment_inspection" "this" {
  name             = "ZPA_Inspection_Example"
  description      = "ZPA_Inspection_Example"
  enabled          = true
  health_reporting = "ON_ACCESS"
  bypass_type      = "NEVER"
  is_cname_enabled = true
  tcp_port_ranges  = ["443", "443"]
  domain_names     = ["jenkins.example.com"]
  segment_group_id = zpa_segment_group.this.id
  server_groups {
    id = [zpa_server_group.this.id]
  }
  common_apps_dto {
    apps_config {
      domain               = "jenkins.example.com"
      application_protocol = "HTTPS"
      application_port     = "443"
      certificate_id       = data.zpa_ba_certificate.jenkins.id
      app_types            = [ "INSPECT" ]
    }
  }
}
```

## Schema

### Required

The following arguments are supported:

- `name` (String) Name. The name of the App Connector Group to be exported.
- `domain_names` (List of String) List of domains and IPs.
- `server_groups` - (Block Set) List of Server Group IDs
  - `id` - (String)

- `segment_group_id` - (string) The unique identifier of the segment group.
- `common_apps_dto` (Block Set, Min: 1) List of applications (e.g., Inspection, Browser Access or Privileged Remote Access)
  - `apps_config:` (Block Set, Min: 1) List of applications to be configured
    - `domain` - (String) Domain name of the inspection application.
      - **NOTE** The domain name configured in this attribute **MUST** also be present in `domain_names` list.

    - `application_protocol` (String) Protocol for the Inspection Application Segment. Supported values: `HTTP` and `HTTPS`. Required by the API on every write. When `auto_app_protect_enabled` is `true` the API manages the protocol itself and reports `DYNAMIC`; the provider keeps the configured value in state so the configuration does not drift.

    - `application_port` - (String) Port for the inspection application.
      - **NOTE** The ports configured in this attribute **MUST** also be present in the port list.

    - `app_types` (List of String) Indicates the type of application as inspection. Supported value: `INSPECT`
    - `certificate_id` (string) - ID of the signing certificate. This field is **required** if the ``application_protocol`` is set to `HTTPS` and `auto_app_protect_enabled` is unset or `false`; the provider rejects the configuration at plan time otherwise. It is not required when `auto_app_protect_enabled` is `true`, because the API manages the certificate itself. The ``certificate_id`` is **NOT** supported if the application_protocol is set to `HTTP`.
    - `trust_untrusted_cert` (Boolean) Enables **Use Untrusted Certificates** for the inspection application, allowing inspection of applications that present self-signed or otherwise untrusted certificates. This is the Terraform equivalent of the Admin Portal's "Use Untrusted Certificates" checkbox and is only supported on this resource.
    - `enabled` (Boolean) Whether this application is enabled or not
- `tcp_port_ranges` - (List of String) TCP port ranges used to access the app.
- `udp_port_ranges` - (List of String) UDP port ranges used to access the app.

-> **NOTE:**  TCP and UDP ports can also be defined using the following model below. We recommend this model as opposed of the legacy model via `tcp_port_ranges` and or `udp_port_ranges`.

- `tcp_port_range` - (Block Set) TCP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.

- `udp_port_range` - (Block Set) UDP port ranges used to access the app.
  - `from:` (String) The starting port for a port range.
  - `to:` (String) The ending port for a port range.

-> **NOTE:** Application segments must have unique ports and cannot have overlapping domain names using the same tcp/udp ports across multiple application segments.

### Optional

- `auto_app_protect_enabled` (Boolean) Enables **Auto App Protection** for the application segment, so that its traffic is inspected by AppProtection. This is the segment-level toggle shown in the Admin Portal and is only supported on this resource. Supported values: `true`, `false`

  ~> **NOTE:** When `auto_app_protect_enabled` is `true`, the ZPA API manages the protocol of every `common_apps_dto.apps_config` entry itself and reports it as `DYNAMIC` (with no port), regardless of the value sent. `application_protocol` must still be set in every `apps_config` block — the API rejects the request without it — and the provider keeps the configured `application_protocol` and `application_port` in state so the configuration does not drift. With `auto_app_protect_enabled` unset or `false`, the configured `HTTP`/`HTTPS` protocol is applied as-is.
- `adp_enabled` (Boolean) Enables **Active Directory Inspection** (Active Directory Protection) for the application segment, so that its Kerberos/LDAP/SMB traffic is inspected. Only supported on this resource, and mutually exclusive with `auto_app_protect_enabled`. Supported values: `true`, `false`
- `description` - (String) Description of the application.
- `bypass_on_reauth` (Boolean) Supported values: `true`, `false`
- `bypass_type` (String) Indicates whether users can bypass ZPA to access applications. Default value is: `NEVER` and supported values are: `ALWAYS`, `NEVER` and `ON_NET`. The value `NEVER` indicates the use of the client forwarding policy.
- `double_encrypt` (Boolean) Whether Double Encryption is enabled or disabled for the app. Supported values are `true` and `false`
- `enabled` - (Boolean) Whether this application is enabled or not. Supported values are `true` and `false`
- `health_reporting` (String) Whether health reporting for the app is Continuous or On Access. Supported values: `NONE`, `ON_ACCESS`, `CONTINUOUS`.
- `health_check_type` (String) Default: `DEFAULT`. Supported values: `DEFAULT`, `NONE`
- `icmp_access_type` - (String) The ICMP access type. Supported values: `PING_TRACEROUTING`, `PING`, `NONE`
- `ip_anchored` - (Boolean) Whether Source IP Anchoring for use with ZIA is enabled or disabled for the application.
- `fqdn_dns_check` - (Boolean) When set to Enabled, Zscaler Client Connector receives CNAME DNS records from the App Connector for FQDN applications. Supported values: `true`, `false`
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

Inspection Application Segment can be imported by using `<APPLICATION SEGMENT ID>` or `<APPLICATION SEGMENT NAME>` as the import ID.

```shell
terraform import zpa_application_segment_inspection.example <application_segment_id>
```

or

```shell
terraform import zpa_application_segment_inspection.example <application_segment_name>
```
