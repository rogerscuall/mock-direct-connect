import boto3

# Replace 'your_local_endpoint_url' with the actual URL of your mock service running locally
endpoint_url = 'http://localhost:8080'  

# Create the client with the custom endpoint URL
client = boto3.client('directconnect', endpoint_url=endpoint_url)

# Replace 'your_connection_name' with a meaningful name for your connection
connection_name = 'MyMockConnection'

# Create the Direct Connect connection
response = client.create_connection(
    connectionName=connection_name,
    location='YourLocation',  # Replace with the desired location for the connection
    bandwidth='1Gbps'  # Replace with the desired bandwidth
    # connectionType='dedicated',  # Replace with the desired connection type
    # providerName='YourProvider',  # Replace with the desired provider name
)

# Check if the connection was created successfully
if 'connectionId' in response:
    connection_id = response['connectionId']
    print(f"Connection with ID {connection_id} created successfully.")
    response = client.describe_connections(connectionId=connection_id)
    print(response)
else:
    print("Connection creation failed.")


client.delete_connection(connectionId=connection_id)

# Verify connection was deleted
try:
    response = client.describe_connections(connectionId=connection_id)
    for connection in response['connections']:
        if connection['connectionId'] == connection_id:
            if connection['connectionState'] == 'deleted':
                print(f"Connection with ID {connection_id} deleted successfully.")
except Exception as e:
    print(f"Connection with ID {connection_id} deleted successfully.")