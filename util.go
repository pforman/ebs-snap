package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func verbose() bool {
	if flag.Lookup("v").Value.String() == "true" {
		return true
	}
	return false
}

func verifyInstance(session *session.Session, instance string) (string, error) {
	svc := ec2.New(session)
	mdsvc := ec2metadata.New(session)
	// if there's no instance specified, go look it up in metadata
	if instance == "" {
		if verbose() {
			println("No instance-id specified, attempting to use local instance")
		}
		i, err := mdsvc.GetMetadata("instance-id")
		if err != nil {
			println("Cannot detect instance-id, exiting.")
			println("Provide an instance-id with '-instance' for remote operation")
			os.Exit(1)
		}
		return i, nil
	} else {
		// go verify the instance exists...
		params := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				aws.String(instance),
			},
		}
		_, err := svc.DescribeInstances(params)
		if err != nil {
			return "", err
		}
		return instance, nil
	}

	// should never get here
	return "", fmt.Errorf("unknown error in VerifyInstance")
}

func verifyRegion(region string) string {
	// if we don't have a region provided, then try to look one up
	// Check the env variable, then metadata
	if region == "" {
		if os.Getenv("AWS_REGION") != "" {
			// We'll be okay, the aws-sdk can use this
			return os.Getenv("AWS_REGION")
		}

		// We're only going to use this for metadata, since we're
		// trying to *find* the region..
		s := session.New()
		svc := ec2metadata.New(s)
		res, err := svc.Region()
		if err != nil {
			println("Error detecting region, exiting.")
			println("Provide a region with '-region' or set AWS_REGION")
			os.Exit(1)
		}
		return res
	}
	// if the user provided something, go with it.
	return region
}
