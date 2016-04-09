package main

import (
	"fmt"
	"os"
  "time"
  "flag"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
  //"github.com/jessevdk/go-flags"
)

func snap (session *session.Session, volumeId string) (string, error) {
  svc := ec2.New(session)

  params := &ec2.CreateSnapshotInput{
    VolumeId:    aws.String(volumeId), // Required
    Description: aws.String("PF snap testing - delete me anytime"),
//    DryRun:      aws.Bool(true),
  }

  resp,err := svc.CreateSnapshot(params)
  if err != nil {
    // Print the error, cast err to awserr.Error to get the Code and
    // Message from an error.
    fmt.Println(err.Error())
    return "", err
  }

  // Pretty-print the response data.
  return *resp.SnapshotId, nil

}

func findVolumeId (session *session.Session, device string , instance string) (string, error) {

  svc := ec2.New(session)

  println("fVID: working with device", device,"instance", instance)
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

func main() {

  //var noop = flag.Bool("noop", true, "test operation, no action")
  var expires = flag.Int("expires", 1, "sets the expiration time in days")
  var instance = flag.String("instance", "i-6ee11663", "instance-id")
  var region = flag.String("region", "us-west-2", "region of instance")
  flag.Parse()

  if flag.Arg(0) == ""  {
    println ("no mount point?  bye!")
    os.Exit(1)
  }

  println ("going to try to snap ", flag.Arg(0))
  device,_ := findDeviceFromMount(flag.Arg(0))
  println ("main: found device", device)

  session := session.New(&aws.Config{Region: aws.String(*region)})

  startingTime := time.Now().UTC()
  fmt.Printf ("expires at %v\n", startingTime.AddDate(0,0,*expires).Round(time.Second))

  println (session, *instance)

  println ("here we go")
  volumeId, _ := findVolumeId(session, device, *instance)
  println ("woop woop found ", volumeId)
  snapId, _ := snap(session, volumeId)
  println ("OMG, snapped",snapId)

  os.Exit(0)

}
