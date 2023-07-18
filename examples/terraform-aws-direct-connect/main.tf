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
  amazon_side_asn = "65001"
}

resource "aws_dx_bgp_peer" "peer" {
  virtual_interface_id = aws_dx_transit_virtual_interface.example.id
  address_family       = "ipv4"
  bgp_asn              = 65002
  # This address is the remote address.
  customer_address     = "74.235.253.85"
  bgp_auth_key         = "1234567890"
  # This address is ignored, use the address of the host running the mock server
  amazon_address       = "1.1.1.1"
}

resource "aws_dx_transit_virtual_interface" "example" {
  connection_id = aws_dx_connection.this.id

  dx_gateway_id  = aws_dx_gateway.this.id
  name           = "tf-transit-vif-example"
  vlan           = 4094
  address_family = "ipv4"
  bgp_asn        = 65005
  mtu           = 1500
}

resource "aws_dx_gateway_association" "example" {
  dx_gateway_id         = aws_dx_gateway.this.id
  # TGW belongs to the EC2 service we only mocking the Direct Connect service
  # Use a fake TGW ID to avoid errors
  associated_gateway_id = "tgw-12345678"

  allowed_prefixes = [
    "10.255.255.0/30",
    "10.255.255.8/30",
  ]
}

# resource "aws_dx_connection" "this1" {
#   name      = "tfx-connection-exmple"
#   location  = "EqSe2"
#   bandwidth = "1Gbps"
#   tags = {
#     Test = "tfx-connection-exmple"
#   }
# }

# resource "aws_dx_bgp_peer" "peer1" {
#   virtual_interface_id = aws_dx_transit_virtual_interface.example1.id
#   address_family       = "ipv4"
#   bgp_asn              = 65002
#   # This address is the remote address.
#   customer_address     = "74.235.253.86"
#   bgp_auth_key         = "1234567890"
#   # This address is ignored, use the address of the host running the mock server
#   amazon_address       = "1.1.1.1"
# }

# resource "aws_dx_transit_virtual_interface" "example1" {
#   connection_id = aws_dx_connection.this1.id

#   dx_gateway_id  = aws_dx_gateway.this.id
#   name           = "tf-transit-vif-example1"
#   vlan           = 4094
#   address_family = "ipv4"
#   bgp_asn        = 65005
#   mtu           = 1500
# }
