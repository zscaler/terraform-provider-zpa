---
page_title: "Best Practices"
description: "Recommended practices for managing Zscaler Private Access (ZPA) with the Zscaler Terraform provider"
---

# Best Practices

This guide collects the recommended practices for managing Zscaler Private Access (ZPA) with Terraform: the general Terraform practices as they apply to this provider, followed by the practices specific to ZPA — the composition chain, application segments, policy rules and microtenants.

-> **Disclaimer:** Please refer to our [General Support Statement](support.md) before proceeding with the use of this provider. For specific errors, see the [Troubleshooting Guide](troubleshooting.md).

## General practices

Terraform practice that applies to the ZPA provider as it does to the other Zscaler providers. ZPA-specific behaviour — the composition chain, application segments, policy rules and microtenants — is in [ZPA-specific practices](#zpa-specific-practices).

### Authentication model

- **DO** use OneAPI (Zidentity OAuth2) when the tenant has migrated — `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET` (or `ZSCALER_PRIVATE_KEY`), `ZSCALER_VANITY_DOMAIN` and `ZPA_CUSTOMER_ID`
- **DO** use the legacy framework (`use_legacy_client = true`) for tenants not yet migrated to Zidentity — the docs call it "the recommended authentication method" for those tenants
- **DO** set `zscaler_cloud` only for `beta` or a FedRAMP cloud (`gov` / `govus`, v4.4.6+) — the docs say it is not required for production commercial clouds
- **DON'T** put a legacy cloud name (`PRODUCTION`, `GOV`, `GOVUS` …) in `zscaler_cloud` — those belong in `zpa_cloud` on the legacy path
- **DON'T** leave both credential sets in one job — the path is chosen **only** by `use_legacy_client` / `ZSCALER_USE_LEGACY_CLIENT`, and the other set is silently ignored
- **DON'T** set both `client_secret` and `private_key` — the argument reference documents them as conflicting
- **DON'T** hard-code credentials in the provider block — the docs warn it "risks secret leakage should this file be committed to public version control"

#### Determining the auth path

The provider decides the path from one input: `use_legacy_client`, falling back to `ZSCALER_USE_LEGACY_CLIENT`. When it is false, `ZPA_CLIENT_ID` / `ZPA_CLIENT_SECRET` are never consulted, however correct they are — which is why a job carrying both namespaces fails with a confusing 401 rather than an explicit error.

API access is granted by role. Create an API Role in ZPA (*Administration > Role Management*), then, if the role does not appear in Zidentity, use *Sync Now* under *Integration > API Resources*.

**Do — empty provider blocks; credentials from the environment**

```hcl
provider "zpa" {}   # OneAPI: ZSCALER_CLIENT_ID, ZSCALER_CLIENT_SECRET, ZSCALER_VANITY_DOMAIN
#         + ZPA_CUSTOMER_ID (required on both paths)
# Legacy: ZPA_CLIENT_ID, ZPA_CLIENT_SECRET, ZPA_CUSTOMER_ID, ZSCALER_USE_LEGACY_CLIENT=true
# FedRAMP over OneAPI (v4.4.6+): ZSCALER_CLOUD=gov | govus
```

**Don't — both namespaces in one job**

```sh
export ZSCALER_CLIENT_ID=...       # OneAPI set
export ZPA_CLIENT_ID=...           # legacy set — silently ignored
export ZSCALER_CLOUD=PRODUCTION    # legacy cloud name — not a OneAPI value
```

### Credential handling

- **DO** let the provider read credentials from environment variables, injected from a secret store at job start (Vault, AWS Secrets Manager, Azure Key Vault — reached through your CI's OIDC)
- **DO** keep state in an encrypted, access-controlled remote backend — secrets you set on resources are stored in state in plain text
- **DO** protect saved plan files (`tfplan`) like state — they contain the values of the variables used to produce them
- **DO** prefer `private_key` over `client_secret` where the customer can manage keys
- **DON'T** count on write-only (`*_wo`) arguments or ephemeral resources — Terraform 1.11 supports them, but the provider does not expose any today
- **DON'T** commit `terraform.tfvars` with secrets — commit `terraform.tfvars.example`
- **DON'T** pass secrets with `-var "client_secret=…"` — it lands in shell history and CI logs
- **DON'T** output secrets at all, even marked `sensitive`
- **DON'T** leave `TF_LOG=DEBUG` on in routine CI — debug logs contain request and response bodies

#### Two kinds of secret

**Provider credentials** — client ID, client secret, private key, and the legacy credentials. Provider configuration is never written to Terraform state. When the provider reads these from the environment, as the provider documentation recommends, they never appear in state, in plan output, or in code.

**Secrets that are resource attributes** — a `provisioning_key`, and the administrator credentials that `zpa_microtenant_controller` outputs for a new microtenant. Terraform must store these in state, and stores them in plain text. Marking variables `sensitive = true` only redacts CLI output; it is not encryption. The control is the backend: encryption at rest, tight IAM, and state locking.

**Don't**

```hcl
variable "api_key" {
  default = "9f8a7b6c5d4e"     # in git forever
}

output "api_key" {
  value     = var.api_key
  sensitive = true              # redacted on screen, still readable
}
```

**Do**

```hcl
provider "zpa" {}   # credentials from ZSCALER_* + ZPA_CUSTOMER_ID — never in state

variable "admin_password" {
  type      = string
  sensitive = true   # from TF_VAR_admin_password; still lands in state
}                    # -> encrypt and restrict the backend

# Output only non-sensitive identifiers
output "segment_id" {
  value = zpa_application_segment.crm.id
}
```

-> **Two operational habits** Add `trivy config --severity HIGH,CRITICAL --exit-code 1 .` to CI to catch secrets in tfvars and insecure backend configuration before merge.

Rotate static credentials every 90 days at most.

### Resource addressing: for_each, not count

- Terraform tracks every resource by an **address**. `count` addresses by position, `for_each` addresses by key
- `count` → `rules[0]`, `rules[1]`, `rules[2]` … **the index is the identity**
- Remove the 2nd of 20 items and every later element shifts down one position, so all 18 addresses now point at different objects — Terraform reads that as **18 destroy/recreates**, not one delete
- `for_each` → `rules["block_gambling"]` — the key is the identity, so removing one key touches **only** that object
- **DO** use `for_each` with a map keyed by **purpose**; use `count` only for on/off (`count = var.enabled ? 1 : 0`)

Terraform does not identify a resource by its name or its content — it identifies it by an address in state. With `count`, that address contains the list index, so **the position in the list becomes the object's identity**. Delete the middle entry from the input list and Terraform sees that every later index now holds a different object, so it destroys and recreates them. One deletion becomes N changes. With `for_each` the address carries a stable key, so removing `"streaming"` removes exactly `rules["streaming"]` and nothing else is re-evaluated.

The cost is higher on a ZPA tenant than in a generic Terraform context. Every needless destroy/recreate is a write against the documented rate limits (10 writes per 10 seconds), and a recreated segment or policy rule does not exist until it is rebuilt. Key by purpose (`"crm"`), never by number (`"app_1"`).

**Don't**

```hcl
resource "zpa_segment_group" "groups" {
  count = length(var.group_names)   # address = groups[0], groups[1] …
  name  = var.group_names[count.index]
}
# Remove "Finance" from the list:
#   groups[1] was Finance -> now Sales -> destroy + recreate
```

**Do**

```hcl
resource "zpa_segment_group" "groups" {
  for_each = var.groups                 # map(object({ description = string }))
  name     = each.key                   # address = groups["finance"]
  enabled  = true
}
```

### Drift and ignore_changes

- **Rare, but it happens:** an attribute you set is missing — or returned empty — in the API response
- Terraform stores what the API returns, so the value never lands in state and every plan shows a change
- **This is an API-side issue, not a provider bug** — the provider can only record what it is given
- Common shapes: booleans omitted when false, API-injected defaults, write-only secrets never returned, list ordering — and, documented in the ZPA troubleshooting guide, `domain_names` that the API lowercases
- **DO** check the provider CHANGELOG and upgrade first — many drift cases are fixed in a later release
- **DON'T** reach for `ignore_changes` as the *first* move, and never `= all` — it hides genuine drift too

#### Two available responses, in order of preference

1. **Raise a bug with the respective engineering team** so it is fixed at the API backend. This is the real fix.
2. **Temporarily** add `lifecycle { ignore_changes = [<field>] }` scoped to that field until the root cause is resolved.

Both the time limit and the scoping matter. A narrowly scoped suppression with a comment and a ticket number is a documented workaround that someone can find and remove later. `ignore_changes = all` is how a configuration silently diverges from the tenant for months, and it will mask the *next*, unrelated problem too.

Write-only fields are a special case and not a defect: passwords, pre-shared keys and API keys are never returned by design. After an import the secret is not recoverable from the API and must be re-supplied through a variable.

**Don't**

```hcl
# Blanket suppression hides genuine drift as well as the bug
lifecycle {
  ignore_changes = all
}
```

**Do**

```hcl
# Narrowly scoped to the affected field, with a comment and a ticket
# reference. Remove once the API fix ships.
lifecycle {
  ignore_changes = [description]   # SUP-XXXX: absent from GET
}
```

### Parallelism and provider tuning

- **DO** run Terraform with its default parallelism — the troubleshooting guide says **do not change `-parallelism`**; it cannot speed up a run
- `-parallelism` is **global to the whole run**. There is no per-resource, per-endpoint or per-provider setting, so lowering it to help one resource type slows down **everything else** in the same configuration
- Rate limiting needs no help: the docs state that "rate limiting requires no configuration" — on a 429 the provider waits out `Retry-After` and retries automatically
- The provider **arguments** `parallelism` (from v4.4.11), `backoff`, `min_wait_seconds` and `max_wait_seconds` (from v4.4.12) are deprecated and ignored — remove them

#### Provider tuning arguments: leave them at the defaults

The argument reference states it explicitly: rate limiting requires no configuration; `max_retries` and `request_timeout` are the **only** supported tuning attributes; neither changes request throughput, so lowering them cannot make plan or apply faster. Environment variables belonging to the Zscaler Go SDK, such as `ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES` or `ZSCALER_CLIENT_REQUEST_TIMEOUT`, are **not** part of the provider configuration and are ignored.

| Provider | `max_retries` default | `request_timeout` default (range) | Deprecated and ignored |
|---|---|---|---|
| **ZPA** | 100 (maximum 100) | 240 s (0–300) | `parallelism` (v4.4.11+); `backoff`, `min_wait_seconds`, `max_wait_seconds` (v4.4.12+) |

A `request_timeout` of `0` selects the SDK's built-in 60-second timeout; it does **not** disable the timeout.

**Don't — tuning that does nothing, or harm**

```sh
export ZSCALER_CLIENT_RATE_LIMIT_MAX_RETRIES=2   # ignored by the provider
export ZSCALER_CLIENT_REQUEST_TIMEOUT=15         # ignored by the provider
provider "zpa" {
  parallelism = 1                # deprecated, no effect
  backoff     = true             # deprecated and ignored from v4.4.12
}
```

#### Conditionally disabled providers: skip_credentials_validation

A shared configuration deployed to several environments may declare the ZPA provider where Zscaler is not used, with every `zpa_*` resource disabled through `count = 0`. Terraform still configures the provider, which validates credentials, so a provider block without credentials fails the run. Set `skip_credentials_validation = true` (or `ZSCALER_SKIP_CREDENTIALS_VALIDATION=true`, v4.4.11+) in those environments only; any resource or data source that does attempt an API call fails with an explanatory error.

**Do**

```hcl
provider "zpa" {
  skip_credentials_validation = !var.zscaler_enabled
}

resource "zpa_segment_group" "this" {
  count = var.zscaler_enabled ? 1 : 0
  name  = "example"
}
```

If a plan is slow on a very large configuration, the troubleshooting guide's procedure applies: don't change `-parallelism` or set SDK variables; plan routine changes with `-refresh=false` (in HCP Terraform / TFE, `TF_CLI_ARGS_plan="-refresh=false"`) and run a full refresh periodically; split configurations that keep growing — see [Execution, scale and state](#execution-scale-and-state).

### State organisation

- How you lay out state and which backend you use is **your call** — any remote backend works
- **DO** split on whatever boundary matches your teams — by product, by rule type, by environment, by approval group
- **DO** use `terraform_remote_state` (or equivalent) to share IDs across states rather than merging them
- **DON'T** import the same object into two states — every object needs exactly one owner

#### ZPA constraints on layout

| Provider | Constraint on state layout and execution | Section |
|---|---|---|
| **ZPA** | One state per microtenant, with the same scope on every resource and data source | [Microtenant scoping](#microtenant-scoping) |

This section is deliberately light on prescription. State layout and backend choice are architecture decisions for the operator. What *is* ours to say is the constraint above, and it is a consequence of API behaviour rather than an opinion.

**One workable layout among many**

```text
infrastructure/zpa/
├── prod/platform/         # connector groups, server groups
├── prod/apps-crm/         # segments + policy for one application domain
└── prod/mt-finance/       # one state per microtenant (ZPA_MICROTENANT_ID)
```

### Modules

- **DO** keep `main.tf`, `variables.tf`, `outputs.tf`, `versions.tf`, `README.md` in every module
- **DO** declare `required_providers` in **every** module, and keep the `version` constraint at the environment (root) level only — the provider troubleshooting guides recommend exactly this to avoid *registry does not have a provider named hashicorp/…*
- **DO** configure providers **only in the root module**
- **DO** give each module a single clear **purpose** — purpose is the boundary, not product
- **DO** declare multiple providers in `required_providers` when the *feature* genuinely spans products
- **DON'T** put a `provider` block in a reusable module
- **DON'T** build a `zscaler_everything` god module bundling unrelated products
- **DON'T** hard-code IDs like `"72058304855457833"` — use data sources

#### Purpose as the module boundary

The blunt formulation — "never mix ZIA and ZPA in one module" — is incorrect. The real boundary is **purpose**: "can I describe what this module does in one sentence without using the word *and*?" A module implementing one feature that spans two products — IP Source Anchoring (a ZIA gateway defined by a ZPA server group and segment), or ZTC forwarding rules that take `src_workload_groups` IDs from the ZIA provider — needs both providers. Even then it declares both only in `required_providers` and configures neither; the root module owns authentication, which is what makes aliasing and multi-tenant composition work.

**Legitimately cross-product — IP Source Anchoring**

```sh
# modules/zia-zpa-source-ip-anchoring/versions.tf — one purpose, two providers
terraform {
  required_providers {
    zia = { source = "zscaler/zia" }
    zpa = { source = "zscaler/zpa" }
  }
}
# no provider block — the root composes and authenticates both
data "zpa_server_group" "this" { name = "Server_Group_IP_Source_Anchoring" }
resource "zia_forwarding_control_zpa_gateway" "this" {
  name = "ZPA_GW01"
  type = "ZPA"
  zpa_server_group {
    external_id = data.zpa_server_group.this.id
    name        = data.zpa_server_group.this.name
  }
}
```

### Version pinning

- **DO** use a pessimistic constraint: `zpa = "~> 4.0"` — v4 is the OneAPI floor
- **DO** set `required_version` from the features you use: `>= 1.3` minimum (stable `optional()` defaults); `>= 1.7` if you use `removed {}` or mock providers. Every provider CHANGELOG states support for Terraform **v1.x**
- **DO** commit `.terraform.lock.hcl` — the troubleshooting guides show the *Failed to install provider* checksum error you get without it — and upgrade in a dedicated PR tested in non-prod first
- **DO** read the provider CHANGELOG before filing a bug
- **DON'T** omit version constraints
- **DON'T** pin an exact GA version (`= 4.0.3`) without a specific reason — it blocks bug fixes
- **DON'T** run `init -upgrade` inside a feature PR

| Provider | Latest release (verified 15 Oct 2026) | Recommended constraint |
|---|---|---|
| `zscaler/zpa` | v4.4.15 (14 Oct 2026) | `~> 4.0` |

Several fixes described in this guide are version-gated. Recent examples: v4.4.12 fixed `zpa_application_segment_inspection` (removed `apps_config` blocks, `DYNAMIC` drift); v4.4.13 refreshes large configurations from list endpoints and adds `iam_idp_id` / `iam_idp_name` to the `zpa_scim_groups` data source; v4.4.14 fixed `payload.size.exceeded` on large segment groups and deprecated `zpa_segment_group.applications`; v4.4.15 adds `tags` to `zpa_application_segment`.

**Do**

```sh
terraform {
  required_version = ">= 1.7.0, < 2.0.0"   # 1.7+: removed {}, mock providers
  required_providers {
    zpa = { source = "zscaler/zpa", version = "~> 4.0" }
  }
}
```

### Naming, variables and outputs

- All lowercase; underscores in Terraform identifiers, hyphens in names that appear in the Zscaler portal
- Names communicate purpose: `internal_hr_portal`, not `app1` / `a` / `test`
- Carry environment context in a local rather than repeating it: `prefix = "${var.environment}-${var.project}"`
- Separate required from optional variables, and give every variable a `description` and a `type`
- Use `validation` blocks to fail at plan time rather than at the API — for example lowercase `domain_names`
- Output identifiers, never secrets ([Credential handling](#credential-handling))

### CI/CD

1. **Validate** — `terraform fmt -check`, `validate`, `tflint`
2. **Scan** — `trivy config`, `checkov`
3. **Plan** — `terraform plan -out=tfplan`, save the artifact
4. **Apply** — `terraform apply tfplan` (the reviewed artifact)

- **DON'T** re-run plan inside the apply job — you'd apply something nobody reviewed
- **DO** keep host-cloud (backend) auth and Zscaler API auth as separate steps
- **DON'T** add an activation stage — a successful ZPA apply is already live

The saved-artifact rule is the one most often bypassed, usually to work around an expired plan. Re-planning in the apply job silently discards human review — it is the CI equivalent of `--auto-approve`.

Authenticating to S3 or Azure Blob for state is a completely different concern from authenticating to the Zscaler API. Keep the two credential paths visibly separate.

On HCP Terraform or Terraform Enterprise, routine plans can skip the refresh with the workspace variable `TF_CLI_ARGS_plan="-refresh=false"`. Keep a scheduled full-refresh plan so changes made outside Terraform are still detected.

-> **Two scheduling habits** Run drift detection on a **daily schedule**, not per-commit.

Never run integration tests against production from PR CI.

### Testing

- **DO** run static analysis (`fmt`, `validate`, `tflint`) on every PR
- **DO** run `terraform plan` against a **non-prod tenant** on every PR
- **DO** use `terraform test` (1.6+) for variable validation
- **DO** tag test resources and run a sweeper to clean up
- **DON'T** expect mock providers to catch API behaviour — no operand validation, no `DUPLICATE_ITEM`
- **DON'T** test in production

Mock providers (Terraform 1.7+) validate *module wiring*, not Zscaler semantics. They are valuable and cheap, but a green mock suite says nothing about whether the API accepts the configuration. Pair them with a sandbox tenant.

ZPA has no activation step: deletes take effect immediately and teardown needs nothing extra.

### Import and brownfield adoption

- **DO** use **zscaler-terraformer** for brownfield adoption — `import` writes the HCL *and* imports into state; `generate` writes HCL only
- **DO** adopt **one resource type per PR**, and confirm an empty `terraform plan` before the next
- **DO** import individual objects by ID (many resources also accept the name — check the resource page)
- **DO** keep the default parallelism, even for a large initial import
- **DO** re-supply write-only secrets after import — the API won't return them
- **DO** use `removed {}` (Terraform 1.7+) to stop managing something without destroying it
- **DON'T** hand-write HCL for a tenant that already exists — that is what the tool is for
- **DON'T** use `terraform state rm` on a custom object — it orphans the API object and the next apply hits `DUPLICATE_ITEM`

zscaler-terraformer is a supported Zscaler CLI that authenticates with the same credentials as the providers, reads a live tenant over the API, and emits the Terraform configuration and the import into state, so the first plan after an import should come back clean. Scope it with `--resources`: the product name for a full sweep, or a single resource type when onboarding incrementally — the pattern to recommend. `--prefix` namespaces the generated names, `--validate` runs `terraform init` and `terraform validate` over the output, and `--collect-logs` captures SDK debug logs for a support case.

Generated HCL is correct but not idiomatic: names come from the API, resources are flat rather than grouped with `for_each`, and there are no variables. Refactor it before it becomes the long-term source of truth.

**Do — brownfield: don't hand-write it**

```sh
zscaler-terraformer import --resources "zpa"                    # whole product
zscaler-terraformer import --resources "zpa_segment_group"      # one type
zscaler-terraformer --prefix "prod" --validate import --resources "zpa"
```

-> **Coverage and caveats** List what the installed version covers with `--supported-resources="zpa"`, and confirm the resource types you need are in it.

It cannot recover write-only values — passwords, pre-shared keys, API keys — because the API never returns them. Re-supply those from your secret store after the import.

It reads `ZPA_CUSTOMER_ID`. Supply credentials as environment variables rather than CLI flags so they do not land in shell history.

`terraform state rm` is the most frequently misused command in this context. It appears to offer a clean exit from a bad state, but it leaves the object live in the tenant and untracked, and the next apply tries to create it again. `removed {}` expresses the same intent safely and is visible in code review.

**Don't**

```sh
terraform state rm zpa_application_segment.crm
# Object still exists in the tenant. Next apply -> DUPLICATE_ITEM or a duplicate.
```

**Do**

```hcl
removed {
  from = zpa_application_segment.crm
  lifecycle { destroy = false }
}
```

## ZPA-specific practices

ZPA has no activation step, no tenant-wide write lock and no contiguous rule-order contract. Its constraints are its own: a strict resource composition chain, application segments with domain and port rules, policy rules expressed as typed operands and ordered through a reorder resource, and microtenant scoping. Everything here comes from the ZPA provider documentation — docs/index.md, the troubleshooting guide and the resource pages — the CHANGELOG and release notes at v4.4.15.

### Scope and resource model

The zscaler/zpa provider registers **55 resources and 71 data sources** (zpa/provider.go at v4.4.15). The excess of data sources is deliberate: roughly thirty objects are provisioned by another system — an identity provider, a posture vendor, the ZPA console, or cloud orchestration — and are readable but not creatable.

| Layer | Purpose | Representative resources and data sources |
|---|---|---|
| Connectivity | Where ZPA reaches into the network | `zpa_app_connector_group`; `zpa_app_connector_controller` (data source) |
| Application | What is being protected | `zpa_application_server`, `zpa_application_segment` |
| Grouping | How applications and servers are grouped for policy | `zpa_segment_group`, `zpa_server_group` |
| Identity | Where users come from | `zpa_idp_controller`, `zpa_scim_groups`, `zpa_saml_attribute` (data sources) |
| Posture | Device-trust signals | `zpa_posture_profile`, `zpa_trusted_network` (data sources) |
| Policy | Allow, deny, forward, inspect, isolate | `zpa_policy_access_rule_v2`, `zpa_policy_forwarding_rule_v2`, `zpa_policy_access_rule_reorder` |
| Tenancy | Microtenant scoping (optional) | `zpa_microtenant_controller` |

#### The composition chain

ZPA enforces a creation order: an application segment cannot be created before the segment group it belongs to, and a policy rule cannot reference a segment that does not yet exist.

```hcl
application server -> server group -> segment group -> application segment -> policy rule
```

Terraform derives this order from attribute references. Configurations that pass literal IDs lose it and fail intermittently, because Terraform has no dependency edge to schedule against. Reference by attribute, not by ID — it orders destroys as well as creates.

#### Segment group membership is set on the application segment (v4.4.14)

- **DO** associate each application segment with its segment group through `segment_group_id` on `zpa_application_segment` — to move a segment, change its `segment_group_id` (updated in place)
- **DON'T** manage membership with the `applications` block of `zpa_segment_group` — it is **deprecated** and will be removed in a future major release; from v4.4.14 adding or removing applications there is rejected at plan time

The `zpa_segment_group` page explains why: "An application segment always belongs to exactly one segment group, and that association can only be changed through the application segment." v4.4.14 also stopped sending the application list on segment-group updates, which fixed `payload.size.exceeded` on large groups. Configurations that list the group's current applications keep working with a deprecation warning; remove the block to clear it.

**Don't**

```hcl
resource "zpa_segment_group" "apps" {
  name = "Apps"
  applications { id = [zpa_application_segment.crm.id] }   # deprecated
}
```

**Do**

```hcl
resource "zpa_segment_group" "apps" {
  name    = "Apps"
  enabled = true
}
resource "zpa_application_segment" "crm" {
  name             = "crm"
  segment_group_id = zpa_segment_group.apps.id   # membership lives here
  # …
}
```

#### Data-source-only objects

- **Identity** — `zpa_idp_controller`, `zpa_saml_attribute`, `zpa_scim_attribute_header`, `zpa_scim_groups`
- **Posture** — `zpa_posture_profile`, `zpa_machine_group`, `zpa_trusted_network`
- **Cloud orchestration** — `zpa_app_connector_controller`, `zpa_private_cloud_controller`, `zpa_cloud_connector_group`
- **Tenant catalogs** — `zpa_enrollment_cert`, `zpa_customer_version_profile`, `zpa_policy_type`, `zpa_isolation_profile`
- **DO** set `iam_idp_id` or `iam_idp_name` on the `zpa_scim_groups` data source (v4.4.13+) when several Zidentity IAM IdPs map to the same ZPA IdP — the group name is no longer unique, and without them the lookup returns the first match. The two are mutually exclusive and require `name`

**Don't**

```hcl
resource "zpa_idp_controller" "okta" { name = "Okta" }   # does not exist
```

**Do**

```hcl
data "zpa_idp_controller" "okta" { name = "Okta" }       # look it up

data "zpa_scim_groups" "engineering" {
  name         = "Engineering"
  idp_name     = data.zpa_idp_controller.okta.name
  iam_idp_name = "Okta-Prod"     # v4.4.13+: disambiguate across IAM IdPs
}
```

### Authentication

- **DO** confirm the auth path before writing any provider block — OneAPI and legacy are both supported
- **DO** supply `ZPA_CUSTOMER_ID` on **both** paths
- **DO** use OneAPI on the FedRAMP clouds — supported from **v4.4.6**; earlier versions are not supported
- **DON'T** set `zscaler_cloud = "PRODUCTION"` on OneAPI — `PRODUCTION` is a legacy `zpa_cloud` value
- **DON'T** set `zscaler_cloud` at all for a commercial production tenant on OneAPI — omit it
- **DON'T** mix `ZSCALER_*` and `ZPA_*` variables in one job — the provider selects one namespace and ignores the other
- **DON'T** keep `parallelism`, `backoff`, `min_wait_seconds` or `max_wait_seconds` in the provider block — ignored from v4.4.11 / v4.4.12 ([Parallelism and provider tuning](#parallelism-and-provider-tuning))

The cloud argument differs between the two paths and the values are not interchangeable. On OneAPI, `zscaler_cloud` is optional and used only for `beta` or the FedRAMP clouds. On legacy, `zpa_cloud` (`ZPA_CLOUD`) is required only when the target is not production; the docs list `BETA`, `GOV`, `GOVUS`, `PREVIEW` and `ZPATWO`.

**OneAPI — Zidentity tenants**

```hcl
provider "zpa" {}
# ZSCALER_CLIENT_ID        (required)
# ZSCALER_CLIENT_SECRET    (required; or ZSCALER_PRIVATE_KEY)
# ZSCALER_VANITY_DOMAIN    (required)
# ZPA_CUSTOMER_ID          (required)
# ZPA_MICROTENANT_ID       (only for a microtenant)
# ZSCALER_CLOUD            (optional — "beta", or "gov" / "govus" from v4.4.6)
```

**Legacy — pre-Zidentity tenants**

```hcl
provider "zpa" {
  use_legacy_client = true
}
# ZPA_CLIENT_ID, ZPA_CLIENT_SECRET, ZPA_CUSTOMER_ID
# ZPA_CLOUD  <- only when not production: BETA | GOV | GOVUS | PREVIEW | ZPATWO
# ZSCALER_USE_LEGACY_CLIENT=true
```

- **DO** inject credentials as environment variables from a secret store; mark credential-bearing variables `sensitive = true`
- **DON'T** commit `client_secret`, `private_key` or `customer_id` in `.tfvars`

The OneAPI client secret is not persisted to state. State does contain object identifiers and resource secrets such as a `provisioning_key` or the administrator credentials output by `zpa_microtenant_controller`, so the backend still needs restricted access. [Credential handling](#credential-handling) covers the general handling.

### Application segments, lifecycle and drift

- **DO** write `domain_names` in **lowercase**, or wrap the list in `lower()` — the troubleshooting guide: the API converts them to lowercase, "which causes a drift"
- **DO** use the `tcp_port_range` / `udp_port_range` blocks (`from` / `to`) — the segment pages: "We recommend this model as opposed of the legacy model via `tcp_port_ranges` and or `udp_port_ranges`"
- **DON'T** give two segments overlapping domain names on the same TCP/UDP ports — "Application segments must have unique ports and cannot have overlapping domain names using the same tcp/udp ports"
- **DO** manage AppProtection (Auto App Protection) and Active Directory Inspection segments with `zpa_application_segment_inspection` **only**
- **DON'T** keep updating a converted segment through `zpa_application_segment` — the next update disables the inspection settings
- **DO** let Terraform order destroys through attribute references
- **DON'T** `terraform state rm` a ZPA resource to move it — use `removed {}` plus an import into the new resource

#### AppProtection and AD Inspection segments

The `zpa_application_segment` page opens with the rule: enabling Auto App Protection or Active Directory Inspection on a segment — for example in the Admin Portal — converts it into an *inspection* application segment, which only `zpa_application_segment_inspection` supports. Import the converted segment into `zpa_application_segment_inspection` and remove it from `zpa_application_segment`.

- Use one `common_apps_dto` block with one `apps_config` block per application; each domain must also appear in `domain_names`, and each port in the segment's port ranges. Removing an `apps_config` block deletes that application
- `application_protocol` is required in every `apps_config`, even with `auto_app_protect_enabled = true` — the API then reports `DYNAMIC`, and the provider keeps the configured value in state (v4.4.12+)
- `certificate_id` is required for HTTPS when `auto_app_protect_enabled` is unset or false (validated at plan time from v4.4.12); `trust_untrusted_cert = true` is the portal's *Use Untrusted Certificates*; `auto_app_protect_enabled` and `adp_enabled` are mutually exclusive

**Do**

```hcl
resource "zpa_application_segment_inspection" "jenkins" {
  name                     = "prod-jenkins"
  domain_names             = ["jenkins.example.com"]   # lowercase
  auto_app_protect_enabled = true
  segment_group_id         = zpa_segment_group.apps.id
  server_groups { id = [zpa_server_group.web.id] }

  tcp_port_range {
    from = "443"
    to   = "443"
  }

  common_apps_dto {
    apps_config {
      domain               = "jenkins.example.com"
      application_protocol = "HTTPS"   # required, even with auto-protect
      application_port     = "443"
      app_types            = ["INSPECT"]
    }
  }
}
```

#### References and drift

ZPA returns `RESOURCE_IN_USE` when destroying a segment group, server group or app connector group that a policy rule still references. Reference the object by attribute, so Terraform computes the destroy order along the composition chain, rather than removing references by hand. Recurring drift has four usual causes: uppercase characters in `domain_names`, booleans omitted from the API response when false, sets whose order the API does not preserve, and write-only fields that are never returned — [Drift and ignore_changes](#drift-and-ignore_changes) covers the response.

### Policy rules, ordering and operands

- **DO** use the **v2** policy resources (`zpa_policy_access_rule_v2`, `zpa_policy_timeout_rule_v2`, `zpa_policy_forwarding_rule_v2`, `zpa_policy_inspection_rule_v2`, `zpa_policy_isolation_rule_v2`) — the docs recommend them when a rule needs more than 1,000 resource criteria
- **DON'T** mix v1 and v2 resources for the same policy type
- **DO** set rule order with **`zpa_policy_access_rule_reorder`** — `rule_order` on the rule resources is **deprecated** "for all ZPA policy types"
- **DON'T** try to reorder rules created by **Zscaler Deception** — not supported; a regular rule's order must be greater than the Deception rules'
- **DO** drop `policy_set_id` — optional since v3.2.0, derived from the policy type, and the docs say it "can be safely removed from existing configurations"
- **DO** reference application segments and groups by attribute, not by literal ID
- **DON'T** assume `lhs` and `rhs` mean the same thing for every `object_type` — they do not
- **DO** use ISO-3166 Alpha-2 in `COUNTRY_CODE` operands — the troubleshooting guide quotes the error for anything else
- **Credential rules:** the first condition block, `object_type = "CONSOLE"`, is mandatory

#### Ordering

ZPA policy rules are evaluated top-down by rule order within a policy set. The documented way to manage that order is the dedicated reorder resource; `rule_order` on the individual resources still works but carries a deprecation warning. Some older examples still set `rule_order` with large offsets. Prefer the reorder resource and contiguous orders — offset only where Zscaler Deception rules require it (see the `zpa_policy_access_rule_reorder` page).

**Do**

```hcl
resource "zpa_policy_access_rule_reorder" "this" {
  policy_type = "ACCESS_POLICY"
  rules {
    id    = zpa_policy_access_rule_v2.allow_crm.id
    order = "1"
  }
  rules {
    id    = zpa_policy_access_rule_v2.deny_default.id
    order = "2"
  }
}
```

#### Operands

A rule carries an ordered list of `conditions` blocks combined by the rule-level operator; within each block, operands are combined by that block's operator. `object_type` is a closed enumeration, and what goes in each operand depends on it (abridged, from the v2 resource pages and their examples):

| object_type | `lhs` (`entry_values`) | `rhs` / `values` |
|---|---|---|
| `APP` | — | `values`: application segment ID |
| `APP_GROUP` | — | `values`: segment group ID |
| `CLIENT_TYPE` | — | `values`: `zpn_client_type_*` client types |
| `SAML` | SAML attribute ID | attribute value to match |
| `SCIM` | SCIM attribute ID | attribute value to match |
| `SCIM_GROUP` | IdP ID (`zpa_idp_controller`) | SCIM group ID (`zpa_scim_groups`) |
| `PLATFORM` | `mac`, `ios`, `windows`, `android`, `linux` | `"true"` |
| `POSTURE` | posture UDID | `"true"` / `"false"` |
| `TRUSTED_NETWORK` | network ID | `"true"` |
| `COUNTRY_CODE` | ISO-3166 Alpha-2 code | `"true"` |

A mismatch returns `400 INVALID_INPUT` or `Invalid operand type`, usually without saying which operand is at fault. Inspecting the planned JSON before apply is the fastest way to isolate it.

```sh
terraform plan -out=tfplan
terraform show -json tfplan | jq '.resource_changes[]
  | select(.type=="zpa_policy_access_rule_v2") | .change.after.conditions'
```

### Microtenant scoping

- **DO** set the scope with `ZPA_MICROTENANT_ID` — the resource pages call the environment variable "the recommended method", and docs/index.md says it is "mandatory when provisioning/managing resources exclusively within a Microtenant"
- **DO** keep one state per microtenant, with the same scope on every resource **and** data source
- **DON'T** put `microtenant_id` inside policy `operands` for `SAML`, `SCIM`, `SCIM_GROUP`, `IDP`, `POSTURE` or `TRUSTED_NETWORK` — not supported; ZPA uses the parent tenant's IDs automatically
- **DON'T** mix scopes in one configuration
- **DO** treat `zpa_microtenant_controller` output as a secret — it outputs the administrator credentials for the new microtenant, which land in state
- Microtenants need the licence and feature flag; the controller is "in limited availability and requires additional license"

~> **Why mixed scoping is dangerous** A resource created in a microtenant but read without the microtenant ID returns 404. Terraform interprets a 404 as "the object no longer exists" and schedules a create. The result is a duplicate object rather than an error.

ZPA is the only Zscaler provider with this concept.

**Do**

```hcl
# One state per microtenant; this job exports ZPA_MICROTENANT_ID
provider "zpa" {}

data "zpa_segment_group" "apps" {
  name = "Apps"          # read in the same microtenant scope
}
```

### Execution, scale and state

- **DO** run with Terraform's default parallelism, and don't set Zscaler Go SDK variables — the troubleshooting guide: "They are not part of the provider's configuration and cannot speed up a run"
- **DO** expect run time to follow the documented rate limits — 20 GET and 10 POST/PUT/DELETE requests per 10 seconds; an apply that changes 2,000 resources takes at least 2,000 seconds
- **DO** upgrade to v4.4.13+ for large configurations — refresh reads from paginated list endpoints
- **DO** plan routine changes with `-refresh=false` and run a full refresh periodically
- **DO** split configurations that keep growing — the guide suggests by segment group or application domain
- **DON'T** carry over a tenant-wide write lock or an activation stage from ZIA or ZTC — ZPA has neither

From v4.4.13, `zpa_application_segment`, `zpa_application_server`, `zpa_segment_group`, `zpa_server_group`, `zpa_app_connector_group` and the v1 and v2 access, timeout and forwarding policy rules refresh from their paginated list endpoint (500 objects per page), fetched once per run and per microtenant. The troubleshooting guide's example: 3,600 application segments now refresh in 8 GETs instead of 3,600. Applies are unchanged — each create, update or delete is still one write.

Policy-rule endpoints reject concurrent modification with `api.concurrent.access.error` ("Unable to modify the resource due to concurrent change requests"). The provider has retried that response since v3.3.23, and v4.4.11 removed the old guidance to limit concurrency to one for the policy resources.

| Version | What changed for configurations | Action |
|---|---|---|
| **v4.4.6** | OneAPI on FedRAMP clouds (`gov` / `govus`) | Use lowercase `gov` / `govus` in `zscaler_cloud` |
| **v4.4.11** | `skip_credentials_validation`; `parallelism` deprecated | Remove `parallelism` |
| **v4.4.12** | Inspection-segment fixes; `backoff` / `*_wait_seconds` deprecated | Upgrade; remove the arguments |
| **v4.4.13** | List-endpoint refresh; `iam_idp_id` / `iam_idp_name` on `zpa_scim_groups` | Upgrade for large configurations |
| **v4.4.14** | `zpa_segment_group.applications` deprecated; `payload.size.exceeded` fixed | Move membership to `segment_group_id` |
| **v4.4.15** | `tags` on `zpa_application_segment` | Optional — tag segments with `tag_key_id` / `tag_value_id` |

- **DO** pin `zpa = "~> 4.0"` — v4 is the minimum for OneAPI; latest v4.4.15 — and commit `.terraform.lock.hcl`
- **DO** configure providers only in the root module
- **DON'T** re-run plan inside the apply job

The ZPA pipeline is the four-stage form from [CI/CD](#cicd) — validate, scan, plan, apply — with no activation stage, and test teardown needs nothing extra. A module implementing a feature that spans ZIA and ZPA, such as IP Source Anchoring, declares both providers in `required_providers` and configures neither ([Modules](#modules)). For brownfield adoption, `zscaler-terraformer import --resources "zpa"` covers ZPA ([Import and brownfield adoption](#import-and-brownfield-adoption)).

### Troubleshooting reference

Rows marked † are quoted from the ZPA troubleshooting guide.

| Error | Likely cause | Fix |
|---|---|---|
| `401 unauthorized` | Wrong credentials, wrong vanity domain, or mixed variable namespaces | Verify the auth path; check `ZPA_CUSTOMER_ID` is set on both paths |
| † `Could not open credentials file … zpa_client_id …` | No credentials configured | Set the environment variables |
| Cloud GOV not supported for OneAPI | Below v4.4.6, or a legacy cloud name in `zscaler_cloud` | Upgrade; use lowercase `gov` / `govus` |
| † `domain_names` change on every plan | Uppercase characters — the API lowercases them | Lowercase them, or `lower()` |
| Segment create rejected | Overlapping domain on the same port as another segment | Make ports / domains unique per segment |
| Inspection settings disabled after an apply | AppProtection / AD Inspection segment in `zpa_application_segment` | Import into `zpa_application_segment_inspection` |
| Removed `apps_config` still present; `DYNAMIC` drift | Below v4.4.12 | Upgrade |
| `payload.size.exceeded` updating a segment group | Below v4.4.14 on a group with many applications | Upgrade; move membership to `segment_group_id` |
| Plan error adding applications to a segment group | `applications` is deprecated (v4.4.14) | Set `segment_group_id` on the application segment |
| `400 INVALID_INPUT` on a policy rule | `lhs` / `rhs` do not match `object_type` | Check operand semantics; inspect planned JSON |
| † `COUNTRY_CODE LHS must be "valid ISO-3166 Alpha-2…"` | Invalid country code in an operand | Use Alpha-2 (`"US"`) |
| Warning: `rule_order` is deprecated | Order set on the rule resource | Use `zpa_policy_access_rule_reorder` |
| `api.concurrent.access.error` | Concurrent policy-rule writes | Retried automatically (v3.3.23+) |
| Duplicate object after an apply | Microtenant scope mismatch → 404 → create | Set `ZPA_MICROTENANT_ID` consistently |
| `zpa_scim_groups` returns the wrong group | Same group name under several IAM IdPs | Set `iam_idp_name` or `iam_idp_id` (v4.4.13+) |
| `RESOURCE_IN_USE` on destroy | Segment or server group still referenced by a policy rule | Reference by attribute so Terraform orders the destroy |
| † Slow `terraform plan` on large configurations | One GET per resource before v4.4.13; rate limits | Upgrade; `-refresh=false` for routine plans; split |
