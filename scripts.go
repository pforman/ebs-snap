package main

import (
  "os/exec"
  "strings"
  //"syscall"
)

func runScript(script string) error {

  var err error

  rawargs := strings.Split(script, " ")

  // squash multiple spaces to prevent "" arguments
  // do this first in case of leading spaces
  args := make([]string,0,10)
  for i,v := range rawargs {
    if rawargs[i] != "" {
      args = append(args,v)
    }
  }

  cmd, args := args[0],args[1:]

  c := exec.Command(cmd,args...)

  // var waitStatus syscall.WaitStatus = 0
  if err = c.Run(); err != nil {
    /* If we want the exit status specifically
    if exitError, ok := err.(*exec.ExitError); ok {
      waitStatus = exitError.Sys().(syscall.WaitStatus)
    }
  } else {
    // Command was successful, this should be 0
    waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
    */
  }

  // Return the whole error.  If it's missing, not executable, fails
  // or whatever, we're going to fail the run anyway
  return err
}
