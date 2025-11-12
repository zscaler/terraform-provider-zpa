---
page_title: "zpa_application_segment_weightedlb_config Resource - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Manage the weighted load balancer configuration for an application segment.
---

# zpa_application_segment_weightedlb_config (Resource)

Use the **zpa_application_segment_weightedlb_config** resource to create, update, or delete the weighted load balancer configuration for an application segment in the Zscaler Private Access cloud.

## Example Usage

```terraform
resource "zpa_application_segment_weightedlb_config" "this" {
  application_name = "app02"
  weighted_load_balancing = true
  application_to_server_group_mappings {
    name = "Example100"
    weight = "100"
  }
}
```

## Schema

### Optional

The following arguments are supported:

* `application_id` - (Optional) The unique identifier of the application segment. Either `application_id` or `application_name` must be provided.
* `application_name` - (Optional) The name of the application segment. Either `application_id` or `application_name` must be provided.
* `microtenant_id` - (Optional) The microtenant identifier to scope the request.
* `weighted_load_balancing` - (Optional) Indicates if the application load balancing configuration for application segments is enabled (`true`) or disabled (`false`).
* `application_to_server_group_mappings` - (Optional) A list of server group mappings associated with the application segment.

The nested `application_to_server_group_mappings` block supports the following:

* `id` - (Optional) The unique identifier of the server group mapping.
* `name` - (Optional) The name of the server group. If only `name` is supplied the provider resolves the identifier automatically.
* `passive` - (Optional) Indicates if the server group should remain passive (`true`) or load balance requests (`false`).
* `weight` - (Optional) The weight assigned to the server group. Higher weights indicate that a greater number of requests are routed to the associated App Connectors; weights of `0` do not route requests.

### Read-Only

The following attributes are exported:

* `application_id` - The unique identifier of the application segment.
* `application_name` - The name of the application segment.
* `application_to_server_group_mappings` - The effective list of server group mappings returned by the API, including resolved identifiers, names, passive status, and weights.


