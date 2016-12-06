package snap

import (
	"bufio"
	//"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func FindVolumeId(session *session.Session, device string, instance string) (string, error) {

	svc := ec2.New(session)

	// println("fVID: working with device", device, "instance", instance)
	params := &ec2.DescribeVolumesInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("attachment.instance-id"),
				Values: []*string{
					aws.String(instance),
				},
			},
			{
				Name: aws.String("attachment.device"),
				Values: []*string{
					aws.String(device),
				},
			},
		},
	}
	resp, err := svc.DescribeVolumes(params)

	if err != nil {
		return "", err
	}

	if len(resp.Volumes) > 0 {
		// is there ever some case when we might get back 2?
		return *resp.Volumes[0].VolumeId, nil
	}
	return "", fmt.Errorf("unable to find volume ID for device %s on instance %s", device, instance)
}

func FindDeviceFromMount(mount string) (string, error) {

	// stub for Mac devel
	if runtime.GOOS != "linux" {
		println("*** Not on Linux?  You're going to have a bad day.")
		println("You should set -device to specify the device to snap.")
		os.Exit(1)
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
				if Verbose() {
					println("Found device", result[1], " mount ", result[2])
				}
				device = result[1]
			}
		}
	}
	if device == "" {
		return device, fmt.Errorf("No device found for mount %s", mount)
	}
	return device, nil
}
