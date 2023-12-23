package trace

import (
	"regexp"
	"testing"
)

func TestInfo(t *testing.T) {
	trace := Info()
	if len(trace) != 1 {
		t.Error("trace length should be 1.")
	}
	var traceRegexp = regexp.MustCompile(`/witness/trace/trace_test\.go:\d+$`)
	if !traceRegexp.MatchString(trace[0]) {
		t.Errorf("trace was not match Regexp:`%s`, Got:`%s`", traceRegexp.String(), trace[0])
	}
}

func TestSkipMyself(t *testing.T) {
	line := "/home/usr/go/pkg/mod/github.com/bayashi/witness@v0.0.8/trace/trace.go"
	if !skipMyself(line) {
		t.Error("Expect to skip, but not skipped")
	}
}
