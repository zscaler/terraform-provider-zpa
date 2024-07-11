######### PASSWORDS IN THIS FILE ARE FAKE AND NOT USED IN PRODUCTION SYSTEMS #########

# Creates Credential of Type "USERNAME_PASSWORD"
resource "zpa_pra_credential_controller" "this" {
    name = "John Doe"
    description = "Created with Terraform"
    credential_type = "USERNAME_PASSWORD"
    user_domain = "acme.com"
    username = "jdoe"
    password = ""
}

# Creates Credential of Type "PASSWORD"
resource "zpa_pra_credential_controller" "this" {
    name = "John Doe"
    description = "Created with Terraform"
    credential_type = "PASSWORD"
    password = ""
}

# Creates Credential of Type "SSH_KEY"
resource "zpa_pra_credential_controller" "this" {
    name = "John Doe"
    description = "Created with Terraform"
    credential_type = "SSH_KEY"
    user_domain = "acme.com"
    username = "jdoe"
    private_key = <<-EOT
-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDEjc8pPoobS0l6
-----END PRIVATE KEY-----
    EOT
}