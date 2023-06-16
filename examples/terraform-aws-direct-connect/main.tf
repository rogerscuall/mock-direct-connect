provider "aws" {
  # ... potentially other provider configuration ...
  access_key                  = "mock_access_key"
  region                      = "us-east-1"
  s3_force_path_style         = true
  secret_key                  = "mock_secret_key"
  skip_credentials_validation = true
  skip_metadata_api_check     = true
  skip_requesting_account_id  = true

  endpoints {
    directconnect = "http://localhost:8080"
    # s3            = "http://localhost:8080"
    # ec2           = "http://localhost:8080"
  }
}

resource "aws_dx_connection" "this" {
  name      = "tfx-connection-exmple"
  location  = "EqSe2"
  bandwidth = "1Gbps"
  tags = {
    Test = "tfx-connection-exmple"
  }
}

resource "aws_dx_gateway" "this" {
  name            = "tf-dxg-exampl123"
  amazon_side_asn = "64512"
}

resource "aws_dx_private_virtual_interface" "this" {
  connection_id = aws_dx_connection.this.id

  name             = "vif-foo"
  vlan             = 4094
  address_family   = "ipv4"
  bgp_asn          = 65002
  dx_gateway_id    = aws_dx_gateway.this.id
  amazon_address       = "1.1.1.1"
  # customer_address     = "20.127.114.9"
  # bgp_auth_key = "1234567890"
  tags = {
    "pvif" = "one"
  }

}

resource "aws_dx_public_virtual_interface" "this" {
  connection_id = aws_dx_connection.this.id

  name           = "vif-foo"
  vlan           = 4090
  address_family = "ipv4"
  bgp_asn        = 65350

  customer_address = "175.45.176.1/30"
  amazon_address   = "175.45.176.2/30"

  route_filter_prefixes = [
    "210.52.109.0/24",
    "175.45.176.0/22",
  ]
}

# This one can be used once the aws_dx_connection resource is available
# data "aws_dx_connection" "this" {
#   name = "tfx-connection-exmple"
# }
