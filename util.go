package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func findDeviceFromMount(mount string) (string, error) {

	// stub for Mac devel
	if runtime.GOOS != "linux" {
		println("*** Not on Linux?  You're going to have a bad day.")
		println("returning static device /dev/xvda for testing only.")
		return "/dev/xvda", nil
	}

	var device string = ""
	// Serious Linux-only stuff happening here...
	// If we can't read /proc/mounts, nothing good can happen.  Get out.
	file := "/proc/mounts"
	_, err := os.Stat(file)
	if err != nil {
		fmt.Printf("Cannot stat file %s: %s", file, err.Error())
		os.Exit(1)
	}
	v, err := os.Open(file)
	if err != nil {
		fmt.Printf("Failed to open %s: %v", file, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(v)
	// leading slash on device to avoid matching things like "rootfs"
	r := regexp.MustCompile(`^(?P<device>/\S+) (?P<mount>\S+) `)
	for scanner.Scan() {
		result := r.FindStringSubmatch(scanner.Text())
		if len(result) > 1 {
			if result[2] == mount {
				println("fDFM: found device", result[1], " mount ", result[2])
				device = result[1]
			}
		}
	}
	if device == "" {
		return device, fmt.Errorf("No device found for mount %s", mount)
	}
	return device, nil
}

func verifyInstance(session *session.Session, instance string) (string, error) {
	svc := ec2.New(session)
	mdsvc := ec2metadata.New(session)
	// if there's no instance specified, go look it up in metadata
	if instance == "" {
		println("No instance-id specified, attempting to use local instance")
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
