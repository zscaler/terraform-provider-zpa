# Application Server Controller Example

This example will show you how to use Terraform to implement the ZPA application server controller resource.
This example codifies [this API](https://help.zscaler.com/zpa/api-reference#/app-server-controller/addAppServer).

To run, configure your ZPA provider as described [Here](https://github.com/willguibr/terraform-provider-zpa/blob/master/website/docs/index.html.markdown)

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
