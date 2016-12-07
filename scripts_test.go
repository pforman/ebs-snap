package snap

import (
	"flag"
	"testing"
	// "github.com/pforman/ebs-snap"
)

func TestRunScript(t *testing.T) {
	// required or Verbose() blows up
	flag.Bool("v", false, "verbose mode, provides more info")

	cases := []struct {
		cmd  string
		want bool
	}{
		{"/usr/bin/true", true},
		{"/usr/bin/false", false},
		{"true", true},
		{"  true", true},
		{" true  still  not  false  ", true},
	}
	for _, v := range cases {
		err := Script(v.cmd)
		if (err != nil) && v.want {
			t.Errorf("runScript(%q) produced an error on a correct case: %v", v.cmd, err)
		}
		if (err == nil) && !v.want {
			t.Errorf("runScript(%q) failed to produce an error on an incorrect case: %v", v.cmd, err)
		}
	}
}
