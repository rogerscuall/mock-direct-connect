# Direct Connect Mock

At the current moment is not possible to test features of Direct Connect without access to one. This is a mock of the Direct Connect API service. It is useful to test Terraform code that uses the Direct Connect API or any other code that uses the Direct Connect API.

This like is very useful [AWS Direct API.](https://frichetten.com/blog/aws-api-protocols/#:~:text=Here%2C%20there%20is%20no%20X-Amz-Target%20header.%20It%20is,to%20intercept%20CLI%20traffic%20and%20inspect%20it%20yourself.)

* The header `X-Amz-Target` is used to determine action to perform. The header has value like this: `OvertureService.CreateConnection` this correspond CreateConnection action. In other words this is a substitute the the Action in the query string.
* The header `Content-Type` is used to determine the format of the request body. The header has value like this: `application/x-amz-json-1.1`.

This is the link to the [AWS Direct Connect API Reference.](https://docs.aws.amazon.com/directconnect/latest/APIReference/API_Operations.html)

This is the link to the [AWS Direct Connect Go SDKv2.](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/directconnect#pkg-overview)

## Create Connections

The create connection makes 3 calls:

1. OvertureService.CreateConnection
2. OvertureService.DescribeConnections
3. OvertureService.DescribeTags
