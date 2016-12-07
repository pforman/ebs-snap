package main

import (
	"flag"
	//"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pforman/ebs-snap"
	"fmt"
)

func main() {

	var region string
	var dryrun = false

	flag.Bool("clean", false, "clean mode, will remove expired snapshots")
	flag.Bool("dryrun", false, "dry run mode, shows what would be cleaned")
	flag.StringVar(&region, "region", "", "`region` of instance")
	flag.Bool("v", false, "verbose mode, provides more info")
	flag.Bool("version", false, "print version string, then exit")
	flag.Parse()

	// version is a quick exit
	if flag.Lookup("version").Value.String() == "true" {
		snap.PrintVersion()
		os.Exit(0)
	}

	if flag.Lookup("dryrun").Value.String() == "true" {
		dryrun = true
	}

	// Make sure they mean it, and weren't just looking for help
	if flag.Lookup("clean").Value.String() != "true" {
		flag.Usage()
		os.Exit(1)
	}

	// If we don't have a region, we're in for trouble
	region = snap.VerifyRegion(region)
	s := session.New(&aws.Config{Region: aws.String(region)})

	err := snap.CleanExpiredSnapshot(s,dryrun)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)

}
