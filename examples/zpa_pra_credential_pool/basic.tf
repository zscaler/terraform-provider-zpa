resource "zpa_pra_credential_pool" "this" {
  name            = "PRACredentialPool01"
  credential_type = "USERNAME_PASSWORD"
  credentials {
    id = [zpa_pra_credential_controller.this.id]
  }
}

resource "zpa_pra_credential_controller" "this" {
  name            = "John Doe"
  description     = "Created with Terraform"
  credential_type = "PASSWORD"
  user_domain     = "acme.com"
  password        = ""
}
