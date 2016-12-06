package main

import (
	"flag"
	"fmt"
	"os"


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
	var precommand, postcommand string
	var expires int
	var ret = 0

	//var noop = flag.Bool("noop", true, "test operation, no action")
	flag.Bool("v", false, "verbose mode, provides more info")
	flag.Bool("version", false, "print version string, then exit")
	flag.IntVar(&expires, "expires", expire_default, "sets the expiration time in days")
	//var instance = flag.String("instance", "i-6ee11663", "instance-id")
	flag.StringVar(&device, "device", "", "device to snapshot (for unmounted volumes only)")
	flag.StringVar(&precommand, "prescript", "", "command to run before snapshot")
	flag.StringVar(&postcommand, "postscript", "", "command to run after snapshot")
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

	if mount == "" {
		println("error: must provide one mount to snapshot (example:  / )")
		flag.Usage()
		os.Exit(1)
	}


	// If we don't have a region, we're in for trouble
	region = snap.VerifyRegion(region)
	s := session.New(&aws.Config{Region: aws.String(region)})

	// If we didn't provide a device, look one up in /proc/mounts
	// This obviously only works on the local host.
	if device == "" {
		res, err := snap.FindDeviceFromMount(mount)
		if err != nil {
			fmt.Printf("error determining device for mount %s: %s\n", mount, err.Error())
			os.Exit(1)
		}
		device = res
	}

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

	// Pre-script
	if precommand != "" {
		err = snap.Script(precommand)
		if err != nil {
			fmt.Printf("error in pre-run command '%s': %s\n", precommand, err.Error())
			os.Exit(1)
		}
	}

	err = snap.CreateSnapshot(s,instance,device,mount,volumeId,expires)
	if err != nil {
		fmt.Printf("error in creating snapshot: %s", err.Error())
		if postcommand == "" {
			os.Exit(1)
		} else {
			// exit with an error, but later...
			ret = 1
			fmt.Println("running postcommand hook")
		}
	}
	if postcommand != "" {
		err = snap.Script(postcommand)
		if err != nil {
			fmt.Printf("error in post-run command '%s': %s\n", precommand, err.Error())
			os.Exit(1)
		}
	}

	os.Exit(ret)

}
