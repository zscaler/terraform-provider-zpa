---
page_title: "zpa_ba_certificate Data source - terraform-provider-zpa"
subcategory: "Browser Access Certificate"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-web-server-certificates
  API documentation https://help.zscaler.com/zpa/configuring-certificates-using-api
  Get information about ZPA Browser Access Certificate in Zscaler Private Access cloud.
---

# zpa_ba_certificate (Data Source)

* [Official documentation](https://help.zscaler.com/zpa/about-web-server-certificates)
* [API documentation](https://help.zscaler.com/zpa/configuring-certificates-using-api)

Use the **zpa_ba_certificate** data source to get information about a browser access certificate created in the Zscaler Private Access cloud. This data source is required when creating a browser access application segment resource.

## Example Usage

```terraform
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  name = "example.acme.com"
}
```

```terraform
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  id = "1234567890"
}
```

## Schema

### Required

The following arguments are supported:

* `name` - (Required) The name of the browser access certificate to be exported.
* `id` - (Optional) The id of the browser access certificate to be exported.

### Read-Only

In addition to all arguments above, the following attributes are exported:

* `cert_chain` - (string) The certificate chain.
* `certificate` - (string) The certificate text is in PEM format.
* `cname` - (string) The canonical name (CNAME DNS records) of the certificate.
* `creation_time` - (string) The time the resource is created.
* `description` - (string) The description of the certificate.
* `issued_by` - (string) The unique identifier the certificate is issued by.
* `issued_to` - (string) The unique identifier the certificate is issued to.
* `modified_time` - (string) The time the certificate is modified.
* `modifiedby` - (string) The unique identifier of the tenant who modified the certificate.
* `san` - (string)  Subject Alternative Name field of the certificate
* `serial_no` - (string) The serial number of the certificate.
* `status` - (string) The status of the certificate.
* `valid_from_in_epochsec` - (string) The start date of the certificate.
* `valid_to_in_epochsec` - (string) The expiration date of the certificate.

~> **Warning**: Notice that certificate and public_keys are omitted from the output.
