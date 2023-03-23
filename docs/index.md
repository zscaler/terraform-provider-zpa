---
layout: "zscaler"
page_title: "Provider: Zscaler Private Access (ZPA)"
description: |-
   The Zscaler Private Access provider is used to interact with ZPA API, to onboard new application segments, segment groups, server groups, application servers and create zero trust access policies. To use this  provider, you must create ZPA API credentials.
---

# Zscaler Private Access Provider (ZPA)

The Zscaler Private Access (ZPA) provider is used to interact with [ZPA](https://www.zscaler.com/products/zscaler-private-access) platform, to onboard new application segments, segment groups, server groups, and create zero trust access policies. To use this  provider, you must create ZPA API credentials. For details on API credentials, please visit the official product [help portal](https://help.zscaler.com/zpa/about-api-keys)

Use the navigation on the left to read about the available resources.

## Zenith Community - ZPA Terraform Provider Introduction

[![ZPA Terraform provider Video Series Ep1](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_terraform_provider_introduction.svg)](https://community.zscaler.com/t/video-zpa-terraform-provider-video-series-ep1/18691)

## Example Usage ZPA Production Cloud

For customers running this provider in their production tenant, the variable `ZPA_CLOUD` is optional. If provided, it must be followed by the value `PRODUCTION`.

```hcl
# Configure ZPA provider source and version
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "~> 2.7.0"
    }
  }
}

provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
}

resouce "zpa_application_segment" "app_segment" {
  # ...
}
```

## Example Usage ZPA Beta, GOV, Preview, and Dev Cloud

For customers who want to use this provider with ZPA Beta, Gov, Preview, and Dev Cloud, the following variable credentials `zpa_cloud` followed by the value `BETA`, `GOV`, `PREVIEW` or `DEV` values or via environment variable `ZPA_CLOUD=BETA`, `ZPA_CLOUD=GOV`, `ZPA_CLOUD=PREVIEW`, `ZPA_CLOUD=DEV`are required.

```hcl
# Configure ZPA provider source and version
terraform {
  required_providers {
    zpa = {
      source = "zscaler/zpa"
      version = "~> 2.7.0"
    }
  }
}

provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
  zpa_cloud             = "BETA" // Use `BETA`, `GOV`, `PREVIEW` or `DEV`
}

resouce "zpa_application_segment" "app_segment" {
  # ...
}
```

## Authentication

The ZPA provider offers various means of providing credentials for authentication. The following methods are supported:

* Directly in the provider block
* Environment variables
* From the JSON config file

### Static credentials

!> **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

Static credentials can be provided by specifying the `zpa_client_id`, `zpa_client_secret` and `zpa_customer_id` arguments in-line in the ZPA provider block:

**Usage:**

``` hcl
provider "zpa" {
  zpa_client_id         = "xxxxxxxxxxxxxxxx"
  zpa_client_secret     = "xxxxxxxxxxxxxxxx"
  zpa_customer_id       = "xxxxxxxxxxxxxxxx"
}
```

### Environment variables

You can provide credentials via the `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID`, `ZPA_CLOUD` environment variables, representing your ZPA API key credentials and customer ID, of your ZPA account, respectively.

~> **NOTE** `ZPA_CLOUD` environment variable is an optional parameter when running this provider in production, but required if running in the ZPA Beta Cloud, Gov Cloud, Preview Cloud or Dev Cloud.

```hcl
provider "zpa" {}
```

**macOS and Linux Usage:**

```sh
export ZPA_CLIENT_ID      = "xxxxxxxxxxxxxxxx"
export ZPA_CLIENT_SECRET  = "xxxxxxxxxxxxxxxx"
export ZPA_CUSTOMER_ID    = "xxxxxxxxxxxxxxxx"
terraform plan
```

**Windows Powershell:**

```powershell
env:ZPA_CLIENT_ID      = 'xxxxxxxxxxxxxxxx'
env:ZPA_CLIENT_SECRET  = 'xxxxxxxxxxxxxxxx'
env:ZPA_CUSTOMER_ID    = 'xxxxxxxxxxxxxxxx'
terraform plan
```

### Configuration file

You can use a configuration file to specify your credentials. The
file location must be `$HOME/.zscaler/credentials.json` on Linux and OS X, or
`"%USERPROFILE%\.zscaler/credentials.json"` for Windows users.
If we fail to detect credentials inline, or in the environment variable, Terraform will check
this location.

Usage:

```hcl
provider "zpa" {}
```

credentials.json file:

```json
{
  "zpa_client_id":"zpa_client_id",
  "zpa_client_secret": "zpa_client_secret",
  "zpa_customer_id": "zpa_customer_id"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `zpa_client_id` - (Required) ZPA client ID, is equivalent to a username.
* `zpa_client_secret` - (Required) ZPA client secret, is equivalent to a secret password.
* `zpa_customer_id` - (Required) ZPA customer ID, is equivalent to your ZPA tenant identification.
* `zpa_cloud` - (Required) ZPA Cloud name `BETA`, `GOV`, `PREVIEW` or `DEV`. Only required when running in the ZPA beta cloud.

### Optional

* `zpa_cloud` - (Optiona) ZPA Cloud name `PRODUCTION`. Optional when running in the ZPA production cloud.

## Support

This template/solution are released under an as-is, best effort, support
policy. These scripts should be seen as community supported and Zscaler
Business Development Team will contribute our expertise as and when possible.
We do not provide technical support or help in using or troubleshooting the components
of the project through our normal support options such as Zscaler support teams,
or ASC (Authorized Support Centers) partners and backline
support options. The underlying product used (Zscaler Private Access API) but the
scripts or templates are still supported, but the support is only for the
product functionality and not for help in deploying or using the template or
script itself. Unless explicitly tagged, all projects or work posted in our
[GitHub repository](https://github.com/zscaler) or sites other
than our official [Downloads page](https://help.zscaler.com/login-tickets)
are provided under the best effort policy.
