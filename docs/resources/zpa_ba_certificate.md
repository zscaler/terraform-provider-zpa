---
subcategory: "Browser Access Certificate"
layout: "zscaler"
page_title: "ZPA: ba_certificate"
description: |-
  Adds a certificate with a private key in Zscaler Private Access cloud.
---

# Data Source: zpa_ba_certificate

Use the **zpa_ba_certificate** creates a browser access certificate with a private key in the Zscaler Private Access cloud. This resource is required when creating a browser access application segment resource.

## Example Usage

```hcl
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  name = "example.acme.com"
}
```

```hcl
# ZPA Browser Access resource
resource "zpa_ba_certificate" "this" {
    name = "server.example.com"
    description = "server.example.com"
    cert_blob = <<-EOT
-----BEGIN CERTIFICATE-----
MIIDyzCCArOgAwIBAgIUekBD+iu64583B3u5ew7Bqj2O5cQwDQYJKoZIhvcNAQEL
BQAwgY0xCzAJBgNVBAYTAkNBMRkwFwYDVQQIDBBCcml0aXNoIENvbHVtYmlhMRIw
EAYDVQQHDAlWYW5jb3V2ZXIxFTATBgNVBAoMDEJELUhhc2hpQ29ycDEVMBMGA1UE
-----END CERTIFICATE-----
    EOT
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the browser access certificate to be created.
* `cert_blob` - (Required) The content of the certificate in PEM format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (string) - The description of the certificate.
