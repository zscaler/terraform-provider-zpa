---
page_title: "zpa_ba_certificate Resource - terraform-provider-zpa"
subcategory: "Browser Access Certificate"
description: |-
  Official documentation https://help.zscaler.com/zpa/about-web-server-certificates
  API documentation https://help.zscaler.com/zpa/configuring-certificates-using-api
  Adds a certificate with a private key in Zscaler Private Access cloud.
---

# zpa_ba_certificate (Resource)

* [Official documentation](https://help.zscaler.com/zpa/about-web-server-certificates)
* [API documentation](https://help.zscaler.com/zpa/configuring-certificates-using-api)

Use the **zpa_ba_certificate** creates a browser access certificate with a private key in the Zscaler Private Access cloud. This resource is required when creating a browser access application segment resource.

## Example Usage

```terraform
# ZPA Browser Access Data Source
data "zpa_ba_certificate" "foo" {
  name = "example.acme.com"
}
```

```terraform
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

## Schema

### Required

The following arguments are supported:

- `name` - (Required) The name of the browser access certificate to be created.
- `cert_blob` - (Required) The content of the certificate in PEM format.

### Optional

In addition to all arguments above, the following attributes are exported:

- `description` - (string) - The description of the certificate.

## Import

This resource does not support importing.

## Let's Encrypt Certbot

This example demonstrates generatoring a domain certificate with letsencrypt
certbot https://letsencrypt.org/getting-started/

```
$ certbot certonly --manual --preferred-challenges dns --key-type rsa -d [DOMAIN]
```

Use letsencrypt's certbot to generate domain certificates in RSA output mode.
The generator's output corresponds to `zpa_ba_certificate` fields in the
following manner.

Zscaler Field          | Certbot file
--------------------|--------------
`certblob`          | `cert.pem`
`certblob`          | `privkey.pem`