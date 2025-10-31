---
page_title: "zpa_application_segment_multimatch_bulk Data source - terraform-provider-zpa"
subcategory: "Application Segment"
description: |-
  Official documentation https://help.zscaler.com/zpa/using-app-segment-multimatch
  API documentation https://help.zscaler.com/zpa/configuring-application-segment-multimatch-using-api
  Get application segments by domain that are incompatible with application segment Multimatch
---

# zpa_application_segment_multimatch_bulk (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/using-app-segment-multimatch)
* [API documentation](https://help.zscaler.com/zpa/configuring-application-segment-multimatch-using-api)

Use the **zpa_application_segment_multimatch_bulk** data source to get application segments by domain that are incompatible with application segment Multimatch Zscaler Private Access cloud. 

## Example Usage

```terraform
data "zpa_application_segment_multimatch_bulk" "this" {
  domain_names = ["server1.bd-hashicorp.com","server2.bd-hashicorp.com"]
}
```

## Schema

### Required

The following arguments are supported:

* `domain_names` - (Required) The list of domains.


