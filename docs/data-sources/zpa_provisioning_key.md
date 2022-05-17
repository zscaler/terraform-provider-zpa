---
subcategory: "Provisioning Key"
layout: "zscaler"
page_title: "ZPA: provisioning_key"
description: |-
  Get information about Provisioning Key in Zscaler Private Access cloud.
---

# zpa_provisioning_key

Use the **zpa_provisioning_key** data source to get information about a provisioning key in the Zscaler Private Access portal or via API. This data source can be referenced in the following ZPA resources:

* App Connector Groups
* Service Edge Groups

-> **NOTE** The ``association_type`` parameter is required in order to distinguish between ``CONNECTOR_GRP`` and ``SERVICE_EDGE_GRP``

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
* `id` - (Optional) The ID of the posture profile to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `creation_time` - (Computed)
* `enabled` - (Computed)
* `expiration_in_epoch_sec` - (Computed)
* `ip_acl` - (Computed)
* `max_usage` - (Computed)
* `modified_by` - (Computed)
* `modified_time` - (Computed)
* `provisioning_key` - (Computed) Ignored in PUT/POST calls.
* `enrollment_cert_id` - (Computed)
* `enrollment_cert_name` - (Computed) Applicable only for GET calls, ignored in PUT/POST calls.
* `ui_config` - (Computed)
* `usage_count` - (Computed)
* `zcomponent_id` - (Computed)
* `zcomponent_name` - (Computed) Applicable only for GET calls, ignored in PUT/POST calls.

:warning: Notice that certificate and public_keys are omitted from the output.
