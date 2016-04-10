package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	//"github.com/jessevdk/go-flags"
)

func main() {

	var instance, region, mount, device string

	//var noop = flag.Bool("noop", true, "test operation, no action")
	flag.Bool("v", false, "verbose")
	var expires = flag.Int("expires", 1, "sets the expiration time in days")
	//var instance = flag.String("instance", "i-6ee11663", "instance-id")
	flag.StringVar(&region, "region", "", "region of instance (for remote snaps only)")
	flag.StringVar(&instance, "instance", "", "instance-id (for remote snaps only)")
	flag.StringVar(&device, "device", "", "device to snapshot (for remote snaps only, be careful!)")
	flag.Parse()

	mount = flag.Arg(0)

	if flag.Arg(1) != "" {
		println("error: multiple mounts provided, or flags after the mount argument")
		flag.Usage()
		os.Exit(1)
	}

	if mount == "" {
		println("error: must provide one mount to snapshot (example:  / )")
		flag.Usage()
		os.Exit(1)
	}

	// Set up the expiration time
	startingTime := time.Now().UTC()
	expireTag := fmt.Sprintf("%v", startingTime.AddDate(0, 0, *expires).Round(time.Second))
	if verbose() {
		println("Expiration time set to", expireTag)
	}

	// If we don't have a region, we're in for trouble
	region = verifyRegion(region)
	session := session.New(&aws.Config{Region: aws.String(region)})

	// If we didn't provide a device, look one up in /proc/mounts
	// This obviously only works on the local host.
	if device == "" {
		res, err := findDeviceFromMount(mount)
		if err != nil {
			fmt.Printf("error determining device for mount %s: %s\n", mount, err.Error())
			os.Exit(1)
		}
		device = res
	}

	instance, err := verifyInstance(session, instance)
	if err != nil {
		fmt.Printf("error finding instance (found '%s'): %s\n", instance, err.Error())
		os.Exit(1)
	}

	volumeId, err := findVolumeId(session, device, instance)
	if err != nil {
		fmt.Printf("error finding volume id for device %s: %s\n", device, err.Error())
		os.Exit(1)
	}

	// old autosnap uses hostname instead of instance-id
	// maybe we should find that...
	snapDesc := fmt.Sprintf("ebs-snap %s:%s:%s", instance, device, mount)
	snapId, err := snapVolume(session, volumeId, snapDesc)
	if err != nil {
		fmt.Printf("error creating snapshot for volume %s: %s\n", volumeId, err.Error())
		os.Exit(1)
	}
	if verbose() {
		println("Created snapshot with id: ", snapId)
	} else {
		println(snapId)
	}
	err = tagSnapshot(session, snapId, expireTag)
	if err != nil {
		println("error in tagging:", err)
		// delete here on error
	}

	os.Exit(0)

}
