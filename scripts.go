package main

import (
  "fmt"
  "os/exec"
  "syscall"
)

func preScript(script string) error {
  cmd := exec.Command(script)

  var waitStatus syscall.WaitStatus
  if err := cmd.Run(); err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      waitStatus = exitError.Sys().(syscall.WaitStatus)
      fmt.Printf("error found: %v\n", err)
      fmt.Printf("exit status was: %d\n", waitStatus.ExitStatus())
      return err
    }
  } else {
    // Command was successful
    waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
    fmt.Printf("no error, exit status was: %d\n", waitStatus.ExitStatus())
    return nil
  }
  return nil

}
