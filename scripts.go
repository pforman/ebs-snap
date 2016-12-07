package snap

import (
	"fmt"
	"os/exec"
	"strings"
)

func Script(script string) error {

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

	cmd, argv := args[0], args[1:]

	c := exec.Command(cmd, argv...)

	stdout, err := c.Output()
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
