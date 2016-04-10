package main

import (
  "bufio"
  "fmt"
	"os"
  "regexp"

)

func findDeviceFromMount (mount string) (string, error) {

  // stub for Mac devel
  //return "/dev/xvda", nil
  var device string = ""
  // Serious Linux-only stuff happening here...
  file :=  "/proc/mounts"
  v, err := os.Open(file)
  if err != nil {
    fmt.Printf("Failed to open %s: %v", file, err)
    return "", err
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

func verifyInstance(instance string) (string, error) {
  // if there's no instance specified, go look it up in metadata
  //if instance == "" &&  {
  return "", nil
}
