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

	/******  Update as necessary  *******/
	const version = "0.1.0"
	const expire_default = 7 // days
	/************************************/

	var instance, region, mount, device string
	var precommand, postcommand string
	var expires int

	//var noop = flag.Bool("noop", true, "test operation, no action")
	flag.Bool("v", false, "verbose mode, provides more info")
	flag.Bool("version", false, "print version string, then exit")
	flag.IntVar(&expires, "expires", expire_default, "sets the expiration time in days")
	//var instance = flag.String("instance", "i-6ee11663", "instance-id")
	flag.StringVar(&region, "region", "", "region of instance (for remote snaps only)")
	flag.StringVar(&instance, "instance", "", "instance-id (for remote snaps only)")
	flag.StringVar(&device, "device", "", "device to snapshot (for remote snaps only, be careful!)")
	flag.StringVar(&precommand, "prescript", "", "command to run before snapshot")
	flag.StringVar(&postcommand, "postscript", "", "command to run after snapshot")
	flag.Parse()

	// version is a quick exit
	if flag.Lookup("version").Value.String() == "true" {
		printVersion(version)
		os.Exit(0)
	}

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
	expireTag := fmt.Sprintf("%v", startingTime.AddDate(0, 0, expires).Round(time.Second))
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

	// Pre-script
	if precommand != "" {
		err = runScript(precommand)
		if err != nil {
			fmt.Printf("error in pre-run command '%s': %s\n", precommand, err.Error())
			os.Exit(1)
		} else {
			if verbose() {
				println("pre-run script completed successfully")
			}
		}
	}

	// old autosnap uses hostname instead of instance-id
	// maybe we should find that...
	snapDesc := fmt.Sprintf("ebs-snap %s:%s:%s", instance, device, mount)
	snapId, err := snapVolume(session, volumeId, snapDesc)
	if err != nil {
		fmt.Printf("error creating snapshot for volume %s: %s\n", volumeId, err.Error())
		os.Exit(1)
	}

	// Tagging
	err = tagSnapshot(session, snapId, volumeId, expireTag)
	if err != nil {
		println("error in tagging:", err)
		// delete here on error, if we can...
	}

	// Post-script
	if postcommand != "" {
		err = runScript(postcommand)
		if err != nil {
			fmt.Printf("error in post-run command '%s': %s\n", postcommand, err.Error())
		} else {
			if verbose() {
				println("post-run script completed successfully")
			}
		}
	}

	// final output
	if verbose() {
		println("Created snapshot with id: ", snapId)
	} else {
		println(snapId)
	}

	os.Exit(0)

}
