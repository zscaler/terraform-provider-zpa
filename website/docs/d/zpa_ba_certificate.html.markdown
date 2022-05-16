---
layout: "zscaler"
page_title: "Zscaler Private Access (ZPA): ba_certificate"
sidebar_current: "docs-datasource-zpa-ba-certificate"
description: |-
  Get information about ZPA Browser Access Certificate in Zscaler Private Access cloud.
---

# zpa_ba_certificate

Use the **zpa_ba_certificate** data source to get information about a browser access certificate created in the Zscaler Private Access cloud. This data source is required when creating a browser access application segment resource.

## Example Usage

```hcl
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  name = "example.acme.com"
}
```

```hcl
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  id = "1234567890"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the browser access certificate to be exported.
* `id` - (Optional) The id of the browser access certificate to be exported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `cert_chain` - (Computed)
* `certificate` - (Computed) The certificate text is in PEM format.
* `cname` - (Computed)
* `creation_time` - (Computed)
* `description` - (Computed)
* `issued_by` - (Computed)
* `issued_to` - (Computed)
* `modified_time` - (Computed)
* `modifiedby` - (Computed)
* `san` - (Computed)
* `serial_no` - (Computed)
* `status` - (Computed)
* `valid_from_in_epochsec` - (Computed)
* `valid_to_in_epochsec` - (Computed)

:warning: Notice that certificate and public_keys are omitted from the output.
