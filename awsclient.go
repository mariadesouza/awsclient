package awsclient

import (
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elb"
)

// ELBInfo : Struct that contains basic info for load balancer
type ELBInfo struct {
	DNSName                string
	NumOfInstancesAttached int
	InstanceIDS            string
	SSLCert                string
	ListenStatus           string
}

//EC2InstanceInfo : Struct that contains basic info for ec2 instnaces
type EC2InstanceInfo struct {
	InstanceID  string
	IPAddress   string
	Status      string
	Name        string
	Environment string
	KeyName     string
}

//AWSClient : struct that has session varibales to use to qury AWS
type AWSClient struct {
	AWSSession *session.Session
	EC2Client  *ec2.EC2
	ELBClient  *elb.ELB
}

//New : Establish aws session using SharedConfigState
func New() *AWSClient {

	var client AWSClient
	// Load session from shared config
	client.AWSSession = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new EC2 client
	client.EC2Client = ec2.New(client.AWSSession)

	client.ELBClient = elb.New(client.AWSSession)

	return &client
}

// FindEC2InstancesByTagName : keyword string can have a wildcard
func (client *AWSClient) FindEC2InstancesByTagName(keyword string) ([]EC2InstanceInfo, error) {

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Name"),
				Values: []*string{
					aws.String(keyword),
				},
			},
		},
	}
	resp, err := client.EC2Client.DescribeInstances(params)
	if err != nil {
		return nil, err
	}
	instances := []EC2InstanceInfo{}
	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			instance := EC2InstanceInfo{}
			instance.Name = "None"
			instance.InstanceID = *inst.InstanceId
			instance.Status = *inst.State.Name
			instance.IPAddress = *inst.PrivateIpAddress
			instance.KeyName = *inst.KeyName
			for _, keys := range inst.Tags {
				if *keys.Key == "Name" {
					instance.Name = url.QueryEscape(*keys.Value)
				}
				if *keys.Key == "environment" {
					instance.Environment = url.QueryEscape(*keys.Value)
				}
			}
			instances = append(instances, instance)
		}
	}
	return instances, nil
}

//FindEC2InstanceByID : Find EC2 Instance By instance ID
func (client *AWSClient) FindEC2InstanceByID(instanceID string) (*EC2InstanceInfo, error) {

	params := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}

	resp, err := client.EC2Client.DescribeInstances(params)
	if err != nil {
		return nil, err
	}
	if len(resp.Reservations[0].Instances) > 1 { // not sure if this ever happens
		return nil, fmt.Errorf("found %d instances ", len(resp.Reservations[0].Instances))
	}
	instance := EC2InstanceInfo{}
	inst := resp.Reservations[0].Instances[0]
	instance.Name = "None"
	instance.InstanceID = *inst.InstanceId
	instance.Status = *inst.State.Name
	instance.IPAddress = *inst.PrivateIpAddress
	instance.KeyName = *inst.KeyName
	for _, keys := range inst.Tags {
		if *keys.Key == "Name" {
			instance.Name = url.QueryEscape(*keys.Value)
		}
		if *keys.Key == "environment" {
			instance.Environment = url.QueryEscape(*keys.Value)
		}
	}
	return &instance, nil

}

//FindLoadBalancersByName : FindLoadBalancers by Name - no wild card
func (client *AWSClient) FindLoadBalancersByName(name string) ([]ELBInfo, error) {

	params := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{
			aws.String(name),
		},
	}

	resp, err := client.ELBClient.DescribeLoadBalancers(params)
	if err != nil {
		return nil, err
	}

	var elbs []ELBInfo
	for _, elbInstance := range resp.LoadBalancerDescriptions {
		var elbInfo ELBInfo
		elbInfo.DNSName = *elbInstance.DNSName
		elbInfo.NumOfInstancesAttached = len(elbInstance.Instances)
		for _, instance := range elbInstance.Instances {
			if elbInfo.InstanceIDS != "" {
				elbInfo.InstanceIDS += ","
			}
			elbInfo.InstanceIDS += *instance.InstanceId

		}
		if len(elbInstance.ListenerDescriptions) > 0 {
			elbInfo.ListenStatus = fmt.Sprintf("Listening on %d port(s).", len(elbInstance.ListenerDescriptions))
		} else {
			elbInfo.ListenStatus = "Not set to listen on any port"
		}

		for _, listener := range elbInstance.ListenerDescriptions {
			if *listener.Listener.LoadBalancerPort == 443 {
				elbInfo.SSLCert = *listener.Listener.SSLCertificateId
			}
			elbInfo.ListenStatus += fmt.Sprintf(" Set to listen on port %d", *listener.Listener.LoadBalancerPort)
		}
		elbs = append(elbs, elbInfo)
	}

	return elbs, nil
}
