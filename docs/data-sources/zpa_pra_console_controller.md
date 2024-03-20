---
subcategory: "Privileged Remote Access"
layout: "zscaler"
page_title: "ZPA): pra_console_controller"
description: |-
  Get information about ZPA privileged remote access console in Zscaler Private Access cloud.
---

# Data Source: zpa_pra_console_controller

The **zpa_pra_console_controller** data source gets information about a privileged remote access console created in the Zscaler Private Access cloud.
This resource can then be referenced in an privileged access policy credential and a privileged access portal resource.

## Example Usage

```hcl
# Retrieve PRA Console by Name
resource "zpa_pra_console_controller" "this" {
  name        = "PRA_Console"
}

# Retrieve PRA Console by ID
resource "zpa_pra_console_controller" "this" {
  id        = "1234567890"
}
```

## Attributes Reference

### Required

* `name` - (Required) The name of the privileged console.
* `id` - (Optional) The ID of the privileged console.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (Required) The description of the privileged console.
* `pra_application` - The Privileged Remote Access application segment resource
    - `id` - (String) The unique identifier of the Privileged Remote Access-enabled application.
* `pra_portals` - The Privileged Remote Access Portal resource
    - `id` - (List) The unique identifier of the privileged portal.

* `microtenant_id` (Optional) The unique identifier of the Microtenant for the ZPA tenant. If you are within the Default Microtenant, pass microtenantId as 0 when making requests to retrieve data from the Default Microtenant. Pass microtenantId as null to retrieve data from all customers associated with the tenant.
