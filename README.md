# Direct Connect Mock

At the current moment is not possible to test features of Direct Connect without access to one. This is a mock of the Direct Connect API service. It is useful to test Terraform code that uses the Direct Connect API or any other code that uses the Direct Connect API.

## How to install

Download the binary from the [release page](https://github.com/rogerscuall/mock-direct-connect/releases) and put it in your path.

## How to use it

1. For testing purposes download the repo and move to the `mock-direct-connect` directory.
1. Download the binary from the [release page](https://github.com/rogerscuall/mock-direct-connect/releases) to the current directory.

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

## Relevant links

This is the link to the [AWS Direct Connect API Reference.](https://docs.aws.amazon.com/directconnect/latest/APIReference/API_Operations.html)

This is the link to the [AWS Direct Connect Go SDKv2.](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/directconnect#pkg-overview)

This like is very useful [AWS Direct API.](https://frichetten.com/blog/aws-api-protocols/)

* The header `X-Amz-Target` is used to determine action to perform. The header has value like this: `OvertureService.CreateConnection` this correspond CreateConnection action. In other words this is a substitute the the Action in the query string.
* The header `Content-Type` is used to determine the format of the request body. The header has value like this: `application/x-amz-json-1.1`.
