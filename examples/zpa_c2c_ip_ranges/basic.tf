## Example Usage - Using IP Range

resource "zpa_c2c_ip_ranges" "this" {
  name            = "Terraform_IP_Range01"
  description     = "Terraform_IP_Range01"
  enabled         = true
  location_hint   = "Created_via_Terraform"
  ip_range_begin  = "192.168.1.1"
  ip_range_end    = "192.168.1.254"
  location        = "San Jose, CA, USA"
  sccm_flag       = true
  country_code    = "US"
  latitude_in_db  = "37.33874"
  longitude_in_db = "-121.8852525"
}


## Example Usage - Using Subnet CIDR
resource "zpa_c2c_ip_ranges" "this" {
  name            = "Terraform_IP_Range01"
  description     = "Terraform_IP_Range01"
  enabled         = true
  location_hint   = "Created_via_Terraform"
  subnet_cidr     = "192.168.1.0/24"
  location        = "San Jose, CA, USA"
  sccm_flag       = true
  country_code    = "US"
  latitude_in_db  = "37.33874"
  longitude_in_db = "-121.8852525"
}
