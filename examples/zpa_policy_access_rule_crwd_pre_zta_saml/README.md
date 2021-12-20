# Create Policy Set Rule for Posture Profile - SAML Attribute

This example will show you how to create a policy set rule to validate if the user's machine is compliant with the Posture Profile conditions according to SAML attribute information.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/policy-set-controller).

To run, configure your ZPA provider as described [Here](https://github.com/SecurityGeekIO/terraform-provider-zpa/blob/master/website/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```

## Destroy 💥

```bash
terraform destroy
```
