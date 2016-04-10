package main

import (
  "bufio"
  "fmt"
  "os"
  "regexp"
  "runtime"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ec2"
  "github.com/aws/aws-sdk-go/aws/ec2metadata"


)

func findDeviceFromMount (mount string) (string, error) {

  // stub for Mac devel
  if runtime.GOOS != "linux" {
    println("*** Not on Linux?  You're going to have a bad day.")
    println("returning static device /dev/xvda for testing only.")
    return "/dev/xvda", nil
  }

  var device string = ""
  // Serious Linux-only stuff happening here...
  // If we can't read /proc/mounts, nothing good can happen.  Get out.
  file :=  "/proc/mounts"
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
        println ("fDFM: found device", result[1], " mount ", result[2])
        device = result[1]
      }
    }
  }
  if device == "" {
    return device, fmt.Errorf("No device found for mount %s", mount)
  }
  return device, nil
}

func verifyInstance (session *session.Session, instance string) (string, error) {
  svc := ec2.New(session)
  mdsvc := ec2metadata.New(session)
  // if there's no instance specified, go look it up in metadata
  if instance == "" {
    println ("No instance-id specified, attempting to use local instance")
    i, err := mdsvc.GetMetadata("instance-id")
    if err != nil {
      return "", err
    }
    println ("cool, returning instance",i,"hyuk-hyuk")
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
