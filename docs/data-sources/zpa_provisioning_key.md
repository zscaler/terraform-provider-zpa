---
subcategory: "Provisioning Key"
layout: "zpa"
page_title: "ZPA: provisioning_key"
description: |-
  Gets a ZPA Provisioning Key details.

---

# zpa_provisioning_key

The **zpa_provisioning_key** data source provides details about a specific provisioning key created manually in the Zscaler Private Access administrator portal or via API.
This data source can be referenced in the following ZPA resources:

1. App Connector Groups
2. Service Edge Groups

## Example Usage

```hcl
# ZPA Posture Profile Data Source
data "zpa_provisioning_key" "example" {
 name = "Provisioning_Key"
}

output "zpa_provisioning_key" {
    value = data.zpa_provisioning_key.example
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the provisioning key.
* `association_type` (Required) Specifies the provisioning key type for App Connectors or ZPA Private Service Edges. The supported values are `CONNECTOR_GRP` and `SERVICE_EDGE_GRP`
