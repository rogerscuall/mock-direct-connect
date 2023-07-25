# Direct Connect Mock

At the current moment is not possible to test features of Direct Connect without access to one.
This is a mock of the Direct Connect API service. It is useful to test API (or Boto3, Terraform or Ansible) that uses the Direct Connect API or any other code that uses the Direct Connect API.
The following API calls are implemented:

* CreateBGPPeer
* CreateConnection
* CreateDXGateway
* CreateDirectConnectGatewayAssociation
* CreatePrivateVirtualInterface
* CreatePublicVirtualInterface
* CreateTransitVirtualInterface
* DeleteBGPPeer
* DeleteConnections
* DeleteDXGateway
* DeleteDirectConnectGatewayAssociation
* DeleteVirtualInterface
* DescribeConnections
* DescribeDXGateways
* DescribeDirectConnectGatewayAssociations
* DescribeVirtualInterfaces
* DescribeTags
* TagResource
* UpdateConnection
* UpdateDXGateway

## How to install

Download the binary from the [release page](https://github.com/rogerscuall/mock-direct-connect/releases) and put it in your path.

## How to use it

1. For testing purposes download the repo and move to the `mock-direct-connect` directory.
1. For testing with Terraform:

```bash
# Start the mock
./mock-direct-connect
# Run Terraform in another terminal
# Move to the terraform directory
cd examples/terraform-aws-direct-connect
# Run Terraform
terraform init
terraform apply
```

1. For testing with Python:

```bash
# Start the mock
./mock-direct-connect
# Run Python in another terminal
# Move to the python directory
cd examples/python
python3 main.py
```

## Relevant links

This is the link to the [AWS Direct Connect API Reference.](https://docs.aws.amazon.com/directconnect/latest/APIReference/API_Operations.html)

This is the link to the [AWS Direct Connect Go SDKv2.](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/directconnect#pkg-overview)

This like is very useful [AWS Direct API.](https://frichetten.com/blog/aws-api-protocols/)

* The header `X-Amz-Target` is used to determine action to perform. The header has value like this: `OvertureService.CreateConnection` this correspond CreateConnection action. In other words this is a substitute the the Action in the query string.
* The header `Content-Type` is used to determine the format of the request body. The header has value like this: `application/x-amz-json-1.1`.

## BGP Peer

The localhost (host where the mock is running) will run BGP with ASN 65001, the IP address would be the primary IP address of the main interface, the mock will auto-discovery this IP. You should see the log similar to this:

```CLI
INFO: 2023/06/25 11:22:11 logger.go:51: create BGP service is true
INFO: 2023/06/25 11:22:11 logger.go:51: creating BGP service
INFO: 2023/06/25 11:22:11 logger.go:51: primary IP address is 192.168.1.100
```
