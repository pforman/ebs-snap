package snap

import (
	"fmt"
	"os/exec"
)

func Script(script string) error {
	cmd := exec.Command(script)

	var err error

	rawargs := strings.Split(script, " ")

	// squash multiple spaces to prevent "" arguments
	// do this first in case of leading spaces
	args := make([]string, 0, 10)
	for i, v := range rawargs {
		if rawargs[i] != "" {
			args = append(args, v)
		}
	}

	cmd, args := args[0], args[1:]

	c := exec.Command(cmd, args...)

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
