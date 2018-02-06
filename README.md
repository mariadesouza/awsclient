# awsclient

## Synopsis

    Package for ease of use of AWS SDK for Go to query AWS resources

## Pre-requisite

### AWS credentials

    The session is created from the configuration values from the shared credentials.
    This support has be enabled by setting the environment variable, AWS_SDK_LOAD_CONFIG=1

### PACKAGES
    go get -u github.com/aws/aws-sdk-go
    go get -u github.com/awslabs/aws-sdk-go/aws
    go get -u github.com/aws/aws-sdk-go/service/ec2
    go get -u github.com/aws/aws-sdk-go/service/elb

## Exported Functions

###  New()

Loads session from shared config and returns a pointer to a client

Returns:
    \*awsclient.AWSClient

### (\*awsclient.AWSClient) FindEC2InstancesByTagName

Takes a keyword  string that can have a wildcard and looks for EC2 instances that have a tag Name with the value

    Input:
      keyword string

    Returns:
      []EC2InstanceInfo,
      error

## (\*awsclient.AWSClient) FindEC2InstanceByID

Finds an ec2 instance based on InstanceID

    Input:
      instanceID string

    Returns:
      []EC2InstanceInfo,
      error


### (\*awsclient.AWSClient) FindLoadBalancersByName

Looks for a load Balancer based on the name passed to it.
No wildcard characters allowed. Full case sensitive names needed.

    Input:
      name string

    Returns:
      []ELBInfo,
      error

## Contributors

Maria DeSouza <maria.desouza@incentivenetworks.com>
