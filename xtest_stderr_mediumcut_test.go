package cwl

import (
	"testing"

	. "github.com/otiai10/mint"
)

func TestDecode_stderr_mediumcut(t *testing.T) {
	f := cwl("stderr-mediumcut.cwl")
	root := NewCWL()
	err := root.Decode(f)
	Expect(t, err).ToBe(nil)
	Expect(t, root.Version).ToBe("v1.0")
	Expect(t, root.Class).ToBe("CommandLineTool")
	Expect(t, root.Doc).ToBe("Test of capturing stderr output in a docker container.")
	Expect(t, root.Requirements[0].Class).ToBe("ShellCommandRequirement")
	Expect(t, len(root.Inputs)).ToBe(0)
	Expect(t, len(root.Outputs)).ToBe(1)
	Expect(t, root.Outputs[0].ID).ToBe("output_file")
	Expect(t, root.Outputs[0].Types[0].Type).ToBe("stderr")
	Expect(t, root.Arguments[0].Binding.ValueFrom).ToBe("echo foo 1>&2")
	Expect(t, root.Arguments[0].Binding.ShellQuote).ToBe(false)
	Expect(t, root.Stderr).ToBe("std.err")
}
