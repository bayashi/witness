package witness

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bayashi/witness/diff"
	"github.com/bayashi/witness/obj"
	"github.com/bayashi/witness/report"
	"github.com/bayashi/witness/trace"
)

// Witness is a context of the fail report
type Witness struct {
	got      *obj.Object
	expect   *obj.Object
	name     string
	messages []map[string]string // additional info as {"label": "message"}
	showDiff bool                // If true, show a diff string for "got" and "expect"
	showRaw  bool                // If true, show raw values as string(raw string or dumped string) for "got" and "expect"
}

// You can write "witness.New(witness.ShowDiff, witness.NotShowRaw)" instead of raw boolean
const ShowDiff, ShowRaw, NotShowDiff, NotShowRaw = true, true, false, false

// You don't need to call `New`. You can call `Got` or `Expect` directly without calling `New` like below.
//
//	witness.Got("abc").Fail(t, "somehow")
//
// You should call `New` when you need to set options for several reports
// in order to avoid calling `ShowDiff` or `ShowRaw` for each report.
//
//	w := witness.New(witness.ShowDiff, witness.ShowRaw)
//	w.Got(123).Fail(t, "Not expected")
//	w.Got("c").Fail(t, "Expected d")
func New(showFlag ...bool) *Witness {
	var showDiff, showRaw = false, false
	if len(showFlag) >= 2 {
		showDiff = showFlag[0]
		showRaw  = showFlag[1]
	} else if len(showFlag) == 1 {
		showDiff = showFlag[0]
	}
	return &Witness{
		showDiff: showDiff,
		showRaw:  showRaw,
	}
}

// Turn on a flag to show diff details
func (w *Witness) ShowDiff() *Witness {
	w.showDiff = true
	return w
}

// Turn on a flag to show raw values
func (w *Witness) ShowRaw() *Witness {
	w.showRaw = true
	return w
}

// Turn on flags for both showDiff and showRaw
func (w *Witness) ShowAll() *Witness {
	return w.ShowDiff().ShowRaw()
}

// Set test name
func Name(n string) *Witness {
	return &Witness{
		name: n,
	}
}

// Set test name
func (w *Witness) Name(n string) *Witness {
	w.name = n

	return w
}

// Set test name by format
func Namef(format string, a ...any) *Witness {
	return &Witness{
		name: fmt.Sprintf(format, a...),
	}
}

// Set test name by format
func (w *Witness) Namef(format string, a ...any) *Witness {
	w.name = fmt.Sprintf(format, a...)

	return w
}

// Set Got value
func Got(v any) *Witness {
	return &Witness{
		got: obj.NewObject(v),
	}
}

// Set Got value
func (w *Witness) Got(v any) *Witness {
	if w.got != nil && w.got.Touch() {
		panic("Already set Got()")
	}

	w.got = obj.NewObject(v)

	return w
}

// Set Expect value
func Expect(v any) *Witness {
	return &Witness{
		expect: obj.NewObject(v),
	}
}

// Set Expect value
func (w *Witness) Expect(v any) *Witness {
	if w.expect != nil && w.expect.Touch() {
		panic("Already set Expect()")
	}

	w.expect = obj.NewObject(v)

	return w
}

func baseReprot(reason string) *report.Failure {
	return report.NewFailure().
		Trace(strings.Join(trace.Info(), "\n\t")).
		Reason(reason)
}

// Set additional messages to show
func (w *Witness) Message(label string, msg string) *Witness {
	w.messages = append(w.messages, map[string]string{label: msg})

	return w
}

var funcFail = func(t *testing.T, r *report.Failure) {
	t.Helper()
	t.Fail()
	t.Error(r.Put())
}

// Do fail with report
func (w *Witness) Fail(t *testing.T, reason string) {
	t.Helper()
	funcFail(t, w.buildReport(t, reason))
}

var funcFailNow = func(t *testing.T, r *report.Failure) {
	t.Helper()
	t.Fatal(r.Put())
}

// Do fail with report and stop running test right now
func (w *Witness) FailNow(t *testing.T, reason string) {
	t.Helper()
	funcFailNow(t, w.buildReport(t, reason))
}

func (w *Witness) buildReport(t *testing.T, reason string) *report.Failure {
	r := baseReprot(reason).Messages(w.messages)

	if w.name != "" {
		r.Name(strings.Join([]string{t.Name(), w.name}, ", "))
	} else {
		r.Name(t.Name())
	}

	if w.got != nil && w.got.Touch() {
		r.Got(w.got)
	}
	if w.expect != nil && w.expect.Touch() {
		r.Expect(w.expect)
	}

	if w.showDiff &&
		(w.got != nil && w.got.Touch()) && (w.expect != nil && w.expect.Touch()) {
		r.Diff(diff.Diff(w.expect.AsRawValue(), w.got.AsRawValue()))
	}

	if w.showRaw {
		r = setRawForReport(w, r)
	}

	return r
}

const rawDataTemplate = "---\n%s\n---"

func setRawForReport(w *Witness, r *report.Failure) *report.Failure {
	if w.got != nil && w.got.Touch() {
		if w.got.IsStringType() {
			r.RawGot(fmt.Sprintf(rawDataTemplate, w.got.AsRawValue()))
		} else if w.got.IsDumpableRawType() {
			r.RawGot(w.got.AsDumpString())
		}
	}
	if w.expect != nil && w.expect.Touch() {
		if w.expect.IsStringType() {
			r.RawExpect(fmt.Sprintf(rawDataTemplate, w.expect.AsRawValue()))
		} else if w.expect.IsDumpableRawType() {
			r.RawExpect(w.expect.AsDumpString())
		}
	}

	return r
}

// Fail is shortcut method. These are same expression.
//
//	witness.Got(got).Fail(t, reason)
//	witness.Fail(t, reason, got)
//
// Fail with 2 values cases are below
//
//	witness.Got(got).Expect(expect).Fail(t, reason)
//	witness.Fail(t, reason, got, expect)
func Fail(t *testing.T, reason string, got any, expect ...any) {
	if len(expect) == 0 {
		Got(got).Fail(t, reason)
	} else {
		Got(got).Expect(expect[0]).Fail(t, reason)
	}
}

// FailNow is shortcut method. There are same expression.
//
//	witness.Got(got).FailNow(t, reason)
//	witness.FailNow(t, reason, got)
//
// FailNow with 2 values cases are below
//
//	witness.Got(got).Expect(expect).FailNow(t, reason)
//	witness.FailNow(t, reason, got, expect)
func FailNow(t *testing.T, reason string, got any, expect ...any) {
	if len(expect) == 0 {
		Got(got).FailNow(t, reason)
	} else {
		Got(got).Expect(expect[0]).FailNow(t, reason)
	}
}

// Diff is to get a diff string of 2 objects for debugging in test
// Two args should be same type. Otherwise, diff string will be a blank string.
func Diff(a any, b any) string {
	return diff.DiffSimple(a, b)
}

// Dump is to get a dumped string by `spew.Sdump` for debugging in test
func Dump(v any) string {
	return obj.NewObject(v).AsDumpString()
}
