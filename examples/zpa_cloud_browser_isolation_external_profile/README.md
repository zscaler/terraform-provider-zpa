# Retrieve ZPA Cloud Browser Isolation External Profile

This example will show you how to use Terraform to implement a ZPA Cloud Browser Isolation External Profile.
This example codifies [this API](https://config.private.zscaler.com/swagger-ui.html#/cbi-profile-controller).

To run, configure your ZPA provider as described [Here](https://github.com/zscaler/terraform-provider-zpa/blob/master/docs/index.html.markdown)

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
