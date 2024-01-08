package trace

import (
	"testing"

	tu "github.com/bayashi/witness/testutil"
)

func TestInfo(t *testing.T) {
	trace := Info()
	if len(trace) != 1 {
		t.Error("trace length should be 1.")
	}

	if ok, msg := tu.Match(`/witness/trace/trace_test\.go:\d+$`, trace[0]); !ok {
		t.Error(msg)
	}
}

func TestSkipMyself(t *testing.T) {
	line := "/home/usr/go/pkg/mod/github.com/bayashi/witness@v0.0.8/trace/trace.go"
	if !skipMyself(line) {
		t.Error("Expect to skip, but not skipped")
	}
}
