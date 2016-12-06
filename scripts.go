package snap

import (
	"fmt"
	"os/exec"
	//"syscall"
)

func Script(script string) error {
	cmd := exec.Command(script)

	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	if Verbose() {
		fmt.Println("=== begin script output ===")
		fmt.Printf("%s", stdout)
		fmt.Println("=== end script output ===")

	}

	return nil

}
