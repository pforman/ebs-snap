package main

import (
	"testing"

	// "github.com/pforman/ebs-snap"
)

func TestRunScript(t *testing.T) {
	cases := []struct {
		cmd  string
		want bool
	}{
		{"/bin/true", true},
		{"/bin/false", false},
		{"true", true},
		{"  true", true},
		{" true  still  not  false  ", true},
	}
	for _, v := range cases {
		err := runScript(v.cmd)
		if (err != nil) && v.want {
			t.Errorf("runScript(%q) produced an error on a correct case: %v", v.cmd, err)
		}
		if (err == nil) && !v.want {
			t.Errorf("runScript(%q) failed to produce an error on an incorrect case: %v", v.cmd, err)
		}
	}
}
