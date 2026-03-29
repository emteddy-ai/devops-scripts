package devops

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type EC2Metadata struct {
	ec2iface.EC2API
}

func NewEC2Metadata() EC2Metadata {
	return EC2Metadata{
		Ec2APIClient: ec2.New(session.New(), &aws.Config{Region: aws.String("us-west-2")}),
	}
}

func (c *EC2Metadata) DescribeInstances() ([]*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{}
	return c.DescribeInstances(input)
}

func (c *EC2Metadata) DescribeInstances(input *ec2.DescribeInstancesInput) ([]*ec2.Instance, error) {
	req, out := c.Ec2APIClient.DescribeInstances(input)
	err := req.Send()
	if err != nil {
		return nil, err
	}
	return out.Reservations, nil
}

func (c *EC2Metadata) GetEc2Instances() ([]*ec2.Instance, error) {
	reservations, err := c.DescribeInstances()
	if err != nil {
		return nil, err
	}
	instances := make([]*ec2.Instance, 0)
	for _, reservation := range reservations {
		instances = append(instances, reservation.Instances...)
	}
	return instances, nil
}

func (c *EC2Metadata) GetEc2InstanceIds() ([]string, error) {
	instances, err := c.GetEc2Instances()
	if err != nil {
		return nil, err
	}
	ids := make([]string, len(instances))
	for i, instance := range instances {
		ids[i] = *instance.InstanceId
	}
	return ids, nil
}

func (c *EC2Metadata) GetEc2InstancePublicDnsNames() ([]string, error) {
	ids, err := c.GetEc2InstanceIds()
	if err != nil {
		return nil, err
	}
	publicDnsNames := make([]string, len(ids))
	for i, id := range ids {
		publicDnsNames[i] = fmt.Sprintf("ec2-%s.us-west-2.compute.amazonaws.com", id)
	}
	return publicDnsNames, nil
}

func (c *EC2Metadata) GetEc2InstancePublicDnsNamesByTags(key string, value string) ([]string, error) {
	instances, err := c.GetEc2Instances()
	if err != nil {
		return nil, err
	}
	publicDnsNames := make([]string, 0)
	for _, instance := range instances {
		tags := make(map[string]*string)
		for _, tag := range instance.Tags {
			tags[*tag.Key] = tag.Value
		}
		if tags[key] != nil && *tags[key] == aws.String(value) {
			publicDnsNames = append(publicDnsNames, *instance.PublicDnsName)
		}
	}
	return publicDnsNames, nil
}

func (c *EC2Metadata) GetEc2InstanceTags(instanceId string) (map[string]*string, error) {
	input := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceId),
				},
			},
		},
	}
	req, out := c.Ec2APIClient.DescribeTags(input)
	err := req.Send()
	if err != nil {
		return nil, err
	}
	tags := make(map[string]*string)
	for _, tag := range out.TagList {
		tags[*tag.Key] = tag.Value
	}
	return tags, nil
}

func GetEc2InstanceTags(instanceId string) (map[string]*string, error) {
	c := NewEC2Metadata()
	return c.GetEc2InstanceTags(instanceId)
}

func GetEc2InstancePublicDnsName(instanceId string) (string, error) {
	c := NewEC2Metadata()
	instances, err := c.GetEc2Instances()
	if err != nil {
		return "", err
	}
	for _, instance := range instances {
		if *instance.InstanceId == instanceId {
			return *instance.PublicDnsName, nil
		}
	}
	return "", fmt.Errorf("instance %s not found", instanceId)
}

func GetEc2InstancePublicDnsNames() ([]string, error) {
	c := NewEC2Metadata()
	return c.GetEc2InstancePublicDnsNames()
}

func GetEc2InstancePublicDnsNamesByTags(key string, value string) ([]string, error) {
	c := NewEC2Metadata()
	return c.GetEc2InstancePublicDnsNamesByTags(key, value)
}

func GetEc2InstanceIds() ([]string, error) {
	c := NewEC2Metadata()
	return c.GetEc2InstanceIds()
}

func GetEc2Instances() ([]*ec2.Instance, error) {
	c := NewEC2Metadata()
	return c.GetEc2Instances()
}