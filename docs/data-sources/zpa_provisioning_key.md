---
subcategory: "Provisioning Key"
layout: "zscaler"
page_title: "ZPA: provisioning_key"
description: |-
  Get information about Provisioning Key in Zscaler Private Access cloud.
---

# Data Source: zpa_provisioning_key

Use the **zpa_provisioning_key** data source to get information about a provisioning key in the Zscaler Private Access portal or via API. This data source can be referenced in the following ZPA resources:

* App Connector Groups
* Service Edge Groups

-> **NOTE** The ``association_type`` parameter is required in order to distinguish between ``CONNECTOR_GRP`` and ``SERVICE_EDGE_GRP``

## Zenith Community - ZPA Provisioning Keys

[![ZPA Terraform provider Video Series Ep3 - Provisioning Keys](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_provisioning_key.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEnCAI/video-zpa-terraform-provider-video-series-ep3-provisioning-keys)

## Example Usage

```hcl
# ZPA Provisioning Key for "CONNECTOR_GRP"
data "zpa_provisioning_key" "example" {
 name = "Provisioning_Key"
 association_type = "CONNECTOR_GRP"
}
```

```hcl
# ZPA Provisioning Key for "SERVICE_EDGE_GRP"
data "zpa_provisioning_key" "example" {
 name = "Provisioning_Key"
 association_type = "SERVICE_EDGE_GRP"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the provisioning key.
* `association_type` (Required) Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are `CONNECTOR_GRP` and `SERVICE_EDGE_GRP`
* `id` - (Optional) The ID of the provisioning key to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `creation_time` - (string)
* `enabled` - (bool)
* `expiration_in_epoch_sec` - (string)
* `ip_acl` - (string)
* `max_usage` - (string)
* `modified_by` - (string)
* `modified_time` - (string)
* `provisioning_key` - (string) Ignored in PUT/POST calls.
* `enrollment_cert_id` - (string)
* `enrollment_cert_name` - (string) Applicable only for GET calls, ignored in PUT/POST calls.
* `ui_config` - (string)
* `usage_count` - (string)
* `zcomponent_id` - (string)
* `zcomponent_name` - (string) Applicable only for GET calls, ignored in PUT/POST calls.
* `microtenant_id` (string) The ID of the microtenant the resource is to be associated with.
* `microtenant_name` (string) The name of the microtenant the resource is to be associated with.

:warning: Notice that certificate and public_keys are omitted from the output.
