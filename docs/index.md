---
layout: "zscaler"
page_title: "Provider: Zscaler Private Access (ZPA)"
description: |-
   The Zscaler Private Access provider is used to interact with ZPA API, to onboard new application segments, segment groups, server groups, application servers and create zero trust access policies. To use this  provider, you must create ZPA API credentials.
---

# Zscaler Private Access Provider (ZPA)

The Zscaler Private Access (ZPA) provider is used to interact with [ZPA](https://www.zscaler.com/products/zscaler-private-access) platform, to onboard new application segments, segment groups, server groups, and create zero trust access policies. To use this  provider, you must create ZPA API credentials. For details on API credentials, please visit the official product [help portal](https://help.zscaler.com/zpa/about-api-keys)

Use the navigation on the left to read about the available resources.

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](guides/troubleshooting.md) for guidance on typical problems.

## Feature Availability and API Parity

-> **Important:** The ZPA Terraform provider maintain parity with publicly available API endpoints. In some instances, certain features or attributes available via the Zscaler UI may not be immediately available through the API, and therefore cannot be included in the Terraform provider. This does not indicate that the provider is lagging behind; rather, it reflects that we implement only the features that are currently exposed by the public API.

If there is a feature or attribute you would like to see included in the provider, you are welcome to:
- Submit a feature request via [GitHub Issues](https://github.com/zscaler/terraform-provider-zpa/issues)
- Contact Zscaler Global Support by opening a support ticket

Our team continuously works with product teams to expand API coverage and will incorporate new features into the provider as they become publicly available through the API.

## Zscaler OneAPI New Framework

The ZPA Terraform Provider now offers support for [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) Oauth2 authentication through [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE** As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity). 

**NOTE** Notice that OneAPI and Zidentity is not currently supported for the following clouds: `GOV` and `GOVUS`. Refer to the [Legacy API Framework](#legacy-api-framework) for more information on how authenticate to these environments

## Zenith Community - ZPA Terraform Provider Introduction

[![ZPA Terraform provider Video Series Ep1](https://raw.githubusercontent.com/zscaler/terraform-provider-zpa/master/images/zpa_terraform_provider_introduction.svg)](https://community.zscaler.com/zenith/s/question/0D54u00009evlEpCAI/video-zpa-terraform-provider-video-series-ep1)

## Examples Usage - Client Secret Authentication

```hcl
# Configure the Zscaler Private Access Provider
terraform {
    required_providers {
        zpa = {
            version = "~> 4.0.0"
            source = "zscaler/zpa"
        }
    }
}

# Configure the ZPA Provider (OneAPI Authentication)
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "zpa" {
  client_id = "[ZSCALER_CLIENT_ID]"
  client_secret = "[ZSCALER_CLIENT_SECRET]"
  vanity_domain = "[ZSCALER_VANITY_DOMAIN]"
  zscaler_cloud = "[ZSCALER_CLOUD]"
  customer_id   = "[ZPA_CUSTOMER_ID]"
}
```

## Examples Usage - Private Key Authentication

```hcl
# Configure the Zscaler Private Access Provider
terraform {
    required_providers {
        zpa = {
            version = "~> 4.0.0"
            source = "zscaler/zpa"
        }
    }
}

# Configure the ZPA Provider (OneAPI Authentication) - Private Key
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "zpa" {
  client_id     = "[ZSCALER_CLIENT_ID]"
  private_key   = "[ZSCALER_PRIVATE_KEY]"
  vanity_domain = "[ZSCALER_VANITY_DOMAIN]"
  zscaler_cloud = "[ZSCALER_CLOUD]"
  customer_id   = "[ZPA_CUSTOMER_ID]"
}
```

**NOTE**: The `zscaler_cloud` is optional and only required when authenticating to other environments i.e `beta`

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

For the resources and data sources examples, please check the [examples](https://github.com/zscaler/terraform-provider-zpa/tree/master/examples) directory.

## Authentication - OneAPI New Framework

As of version v4.0.0, this provider supports authentication via the new Zscaler API framework [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi)

Zscaler OneAPI uses the OAuth 2.0 authorization framework to provide secure access to Zscaler Private Access (ZPA) APIs. OAuth 2.0 allows third-party applications to obtain controlled access to protected resources using access tokens. OneAPI uses the Client Credentials OAuth flow, in which client applications can exchange their credentials with the authorization server for an access token and obtain access to the API resources, without any user authentication involved in the process.

* [ZPA API](https://help.zscaler.com/oneapi/understanding-oneapi#:~:text=Workload%20Groups-,ZPA%20API,-Zscaler%20Private%20Access)

### OneAPI (API Client Scope)

OneAPI Resources are automatically created within the ZIdentity Admin UI based on the RBAC Roles
applicable to APIs within the various products. For example, in ZPA, navigate to `Administration -> Role
Management` and select `Add API Role`.

Once this role has been saved, return to the ZIdentity Admin UI and from the Integration menu
select API Resources. Click the `View` icon to the right of Zscaler APIs and under the ZIA
dropdown you will see the newly created Role. In the event a newly created role is not seen in the
ZIdentity Admin UI a `Sync Now` button is provided in the API Resources menu which will initiate an
on-demand sync of newly created roles.

### Default Environment variables

You can provide credentials via the `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET`, `ZSCALER_VANITY_DOMAIN`, `ZSCALER_CLOUD` environment variables, representing your Zidentity OneAPI credentials `clientId`, `clientSecret`, `vanityDomain` and `zscaler_cloud` respectively.

| Argument        | Description                                                                                         | Environment Variable     |
|-----------------|-----------------------------------------------------------------------------------------------------|--------------------------|
| `client_id`     | _(String)_ Zscaler API Client ID, used with `clientSecret` or `PrivateKey` OAuth auth mode.         | `ZSCALER_CLIENT_ID`      |
| `client_secret` | _(String)_ Secret key associated with the API Client ID for authentication.                         | `ZSCALER_CLIENT_SECRET`  |
| `privateKey`    | _(String)_ A string Private key value.                                                              | `ZSCALER_PRIVATE_KEY`    |
| `customer_id`   | _(String)_ A string that contains the ZPA customer ID which identifies the tenant                   | `ZPA_CUSTOMER_ID`    |
| `microtenant_id`| _(String)_ A string that contains the ZPA microtenant ID which identifies the tenant                | `ZPA_MICROTENANT_ID`    |
| `vanity_domain` | _(String)_ Refers to the domain name used by your organization.                                     | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud`         | _(String)_ The name of the Zidentity cloud, e.g., beta.                                             | `ZSCALER_CLOUD`          |

### Alternative OneAPI Cloud Environments

OneAPI supports authentication and can interact with alternative Zscaler enviornments i.e `beta`, `alpha` etc. To authenticate to these environments you must provide the following values:

| Argument         | Description                                                                                         |   | Environment Variable     |
|------------------|-----------------------------------------------------------------------------------------------------|---|--------------------------|
| `vanity_domain`   | _(String)_ Refers to the domain name used by your organization |   | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud`          | _(String)_ The name of the Zidentity cloud i.e beta      |   | `ZSCALER_CLOUD`          |

For example: Authenticating to Zscaler Beta environment:

```sh
export ZSCALER_VANITY_DOMAIN="acme"
export ZSCALER_CLOUD="beta"
```

## Legacy API Framework

### ZPA native authentication

* As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

### Examples Usage

```hcl
# Configure the Zscaler Internet Access Provider
terraform {
    required_providers {
        zpa = {
            version = "~> 4.0.0"
            source = "zscaler/zpa"
        }
    }
}

# Configure the ZPA Provider (Legacy Authentication)
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "zpa" {
  zpa_client_id            = "[ZPA_CLIENT_ID]"
  zpa_client_secret        = "[ZPA_CLIENT_SECRET]"
  zpa_customer_id          = "[ZPA_CUSTOMER_ID]"
  zpa_cloud                = "[ZPA_CLOUD]"
  use_legacy_client        = "[ZSCALER_USE_LEGACY_CLIENT]"
}
```

### Environment variables

You can provide credentials via the `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID`, `ZPA_CLOUD` environment variables, representing your ZPA API key credentials and customer ID, of your ZPA account, respectively.

~> **NOTE 1** `ZPA_CLOUD` environment variable is an optional parameter when running this provider in production; however, this parameter is **ONLY** required when provisioning resources in in the following ZPA Clouds `BETA`, `GOV`, `GOVUS`, `PREVIEW` or `ZPATWO`

~> **NOTE 2** `ZPA_MICROTENANT_ID` environment variable is an optional parameter when provisioning resources within a ZPA microtenant

```terraform
provider "zpa" {}
```

**macOS and Linux Usage:**

```sh
export ZPA_CLIENT_ID                = "xxxxxxxxxxxxxxxx"
export ZPA_CLIENT_SECRET            = "xxxxxxxxxxxxxxxx"
export ZPA_CUSTOMER_ID              = "xxxxxxxxxxxxxxxx"
export ZSCALER_USE_LEGACY_CLIENT    = true
terraform plan
```

**Windows Powershell:**

```powershell
$env:ZPA_CLIENT_ID='xxxxxxxxxxxxxxxx'
$env:ZPA_CLIENT_SECRET='xxxxxxxxxxxxxxxx'
$env:ZPA_CUSTOMER_ID='xxxxxxxxxxxxxxxx'
$env:ZSCALER_USE_LEGACY_CLIENT=true
terraform plan
```

## Argument Reference

The following arguments are supported:

### Required

* ``zpa_client_id`` - (Required) ZPA client ID, is equivalent to a username.
* ``zpa_client_secret`` - (Required) ZPA client secret, is equivalent to a secret password.
* ``zpa_customer_id`` - (Required) ZPA customer ID, is equivalent to your ZPA tenant identification.
* ``zpa_cloud`` - (Required) ZPA Cloud name `BETA`, `GOV`, `GOVUS`, `PREVIEW` or `ZPATWO` clouds.
* ``use_legacy_client`` - (Required) Enable legacy API client. Supported values `true` or `false`.

### Optional

* `zpa_cloud` - (Optional) ZPA Cloud name `PRODUCTION`.

~> **NOTE** `ZPA_CLOUD` environment variable is an optional parameter when running this provider in production; however, this parameter is **ONLY** required when provisioning resources in in the following ZPA Clouds `BETA`, `GOV`, `GOVUS`, `PREVIEW` or `ZPATWO`

### Zscaler Private Access Microtenant

A Microtenant is a delegated administrator responsibility that is assigned to an admin by an admin with Microtenant administrator privileges. Microtenants are defined by an authentication domain and assigned to admins based on country, department, and company for role-based administration control. A Microtenant is created within a tenant and is used when departments or subsidiaries within an organization want to manage their configurations independently.[Read More](https://help.zscaler.com/zpa/about-microtenants)

To manage a microtenant using the ZPA Terraform provider, the administrator for the parent or default tenant, must first provision the microtenant using the resource `zpa_microtenant_controller`. The resource will output the administrator credentials for the new microtenant, which can then be provided to the microtenant administrator.

The microtenant administrator can then create his own microtenant API credentials required to authenticate via API to the ZPA platform. From that point, the administrator can then individually manage his own resources in an isolated manner.
When authenticating to microtenant via API using the ZPA Terraform provider, the administrator must provide the following environment variable credentials: `ZPA_CLIENT_ID`, `ZPA_CLIENT_SECRET`, `ZPA_CUSTOMER_ID`, `ZPA_CLOUD`, `ZPA_MICROTENANT_ID`

~> **NOTE 1** The environment variable `ZPA_MICROTENANT_ID` is mandatory when provisioning/managing resources exclusively within a Microtenant.

~> **NOTE 2** `ZPA_CLOUD` environment variable is an optional parameter when running this provider in production; however, this parameter is **ONLY** required when provisioning resources in in the following ZPA Clouds `BETA`, `GOV`, `GOVUS`, `PREVIEW` or `ZPATWO`

## Argument Reference - OneAPI

Before starting with this Terraform provider you must create an API Client in the Zscaler Identity Service portal [Zidentity](https://help.zscaler.com/zidentity/what-zidentity) or have create an API key via the legacy method.

* `client_id` - (Required) This is the client ID for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_ID` environment variable.

* `client_secret` - (Required) This is the client secret for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_SECRET` environment variable. `client_secret` conflicts with `private_key`.

* `private_key` - (Required) This is the private key for obtaining the API token (can be represented by a filepath, or the key itself). It can also be sourced from the `ZSCALER_PRIVATE_KEY` environment variable. `private_key` conflicts with `client_secret`. The format of the PK is PKCS#1 unencrypted (header starts with `-----BEGIN RSA PRIVATE KEY-----` or PKCS#8 unencrypted (header starts with `-----BEGIN PRIVATE KEY-----`).

* `vanity_domain` - (Required) This refers to the domain name used by your organization.. It can also be sourced from the `ZSCALER_VANITY_DOMAIN`.

* `zscaler_cloud` - (Required) This refers to Zscaler cloud name where API calls will be directed to i.e `beta`. It can also be sourced from the `ZSCALER_CLOUD`.

* `customer_id` - (Required) A string that contains the the ZPA customer ID which identifies the tenant. Can also be sourced from the `ZPA_CUSTOMER_ID` environment variable. It is required when interacting with ZPA Cloud via OneAPI framework.

* `microtenant_id` - (Optional) A string that contains the the ZPA customer ID which identifies the microtenant ID. Can also be sourced from the `ZPA_MICROTENANT_ID` environment variable. It is required when interacting with ZPA Microtenant ID feature via OneAPI framework.

* `http_proxy` - (Optional) This is a custom URL endpoint that can be used for unit testing or local caching proxies. Can also be sourced from the `ZSCALER_HTTP_PROXY` environment variable.

* `parallelism` - (Optional) Number of concurrent requests to make within a resource where bulk operations are not possible. [Learn More](https://help.zscaler.com/oneapi/understanding-rate-limiting)

* `max_retries` - (Optional) Maximum number of retries to attempt before returning an error, the default is `5`.

* `request_timeout` - (Optional) Timeout for single request (in seconds) which is made to Zscaler, the default is `0` (means no limit is set). The maximum value can be `300`.

* `zpa_client_id` - (Required) A string that contains the legacy ZPA client ID.  Can also be sourced from the `ZPA_CLIENT_ID` environment variable.. Required when setting the attribute `use_legacy_client`

* `zpa_client_secret` - (Required) A string that contains the the legacy ZPA client Secret. Can also be sourced from the `ZPA_CLIENT_SECRET` environment variable. Required when setting the attribute `use_legacy_client`

* `zpa_customer_id` - (Required) A string that contains the the legacy ZPA customer ID which identifies the tenant. Can also be sourced from the `ZPA_CUSTOMER_ID` environment variable. Required when setting the attribute `use_legacy_client`

* `microtenant_id` - (Optional) A string that contains the the ZPA customer ID which identifies the microtenant ID. Can also be sourced from the `ZPA_MICROTENANT_ID` environment variable. Required when interacting a microtenant via the legacy API framework by setting the attribute `use_legacy_client` 

* `zpa_cloud` - (Optional) This refers to the the legacy ZPA cloud name where api calls will be forward to. Can also be sourced from the `ZPA_CLOUD` environment variable. Required when interacting with a ZPA cloud other than `PRODUCTION` via the legacy API framework by setting the attribute `use_legacy_client`. 

Currently the following cloud names are supported:
  * `PRODUCTION`
  * `BETA`
  * `GOV`
  * `GOVUS`
  * `ZPATWO`
  * `ZSPREVIEW`

* `use_legacy_client` - (Optional) This parameter is required when using the legacy API framework. Can also be sourced from the `ZSCALER_USE_LEGACY_CLIENT` environment variable.
