package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pforman/ebs-snap"
	//"github.com/jessevdk/go-flags"
)

func main() {

	/******  Update as necessary  *******/
	const version = "0.1.0"
	const expire_default = 7 // days
	/************************************/

	var instance, region, mount, device string
	var expires int

	//var noop = flag.Bool("noop", true, "test operation, no action")
	flag.Bool("v", false, "verbose mode, provides more info")
	flag.Bool("version", false, "print version string, then exit")
	flag.IntVar(&expires, "expires", expire_default, "sets the expiration time in days")
	flag.StringVar(&region, "region", "", "`region` of instance")
	flag.StringVar(&instance, "instance", "", "`instance-id` of the instance to snapshot")
	flag.StringVar(&device, "device", "", "`device` to snapshot")
	flag.Parse()

	// version is a quick exit
	if flag.Lookup("version").Value.String() == "true" {
		snap.PrintVersion(version)
		os.Exit(0)
	}

	mount = flag.Arg(0)

	if flag.Arg(1) != "" {
		println("error: multiple mounts provided, or flags after the mount argument")
		flag.Usage()
		os.Exit(1)
	}

	if device == "" {
		println("error: -device is required for remote snapshots")
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
	expireTag := fmt.Sprintf("%v", startingTime.AddDate(0, 0, expires).Round(time.Second))
	if snap.Verbose() {
		println("Expiration time set to", expireTag)
	}

	// If we don't have a region, we're in for trouble
	region = snap.VerifyRegion(region)
	s := session.New(&aws.Config{Region: aws.String(region)})

	instance, err := snap.VerifyInstance(s, instance)
	if err != nil {
		fmt.Printf("error finding instance (found '%s'): %s\n", instance, err.Error())
		os.Exit(1)
	}

	volumeId, err := snap.FindVolumeId(s, device, instance)
	if err != nil {
		fmt.Printf("error finding volume id for device %s: %s\n", device, err.Error())
		os.Exit(1)
	}

	snap.CreateSnapshot(s,instance,device,mount,volumeId,expires)
	/*
	// old autosnap uses hostname instead of instance-id
	// maybe we should find that...
	snapDesc := fmt.Sprintf("ebs-snap %s:%s:%s", instance, device, mount)
	snapId, err := snap.EC2Snapshot(s, volumeId, snapDesc)
	if err != nil {
		fmt.Printf("error creating snapshot for volume %s: %s\n", volumeId, err.Error())
		os.Exit(1)
	}
	if snap.Verbose() {
		println("Created snapshot with id: ", snapId)
	} else {
		println(snapId)
	}
	err = snap.TagSnapshot(s, snapId, volumeId, expireTag)
	if err != nil {
		println("error in tagging:", err)
		// delete here on error, if we can...
	}
	*/
	os.Exit(0)

}
