package witness

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/bayashi/witness/report"
)

// Global variables to check result
var res string
var failnow bool

func stub() {
	res = ""
	funcFail = func(t *testing.T, r *report.Failure) {
		res = r.Put()
		failnow = false
	}
	funcFailNow = func(t *testing.T, r *report.Failure) {
		funcFail(t, r)
		failnow = true
	}
}

func TestError(t *testing.T) {
	stub()

	err := fmt.Errorf("error example %d", 123)
	Got(err).Fail(t, "oops")

    // Fail reason:    oops
    // Type:           Got:*errors.errorString
    // Actually got:   error example 123

	gotTypeRe := regexp.MustCompile(`Type:\s*\tGot:\*errors\.errorString`)
	if gotTypeRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", gotTypeRe.String(), res)
	}

	gotRe := regexp.MustCompile(`Actually got:\s*\t\(\*errors\.errorString\)\([0-9a-fx]+\)\(error example 123\)`)
	if gotRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", gotRe.String(), res)
	}
}

func TestFailGot(t *testing.T) {
	stub()

	got := "got string"
	reason := "failure reason"
	Got(got).Name("Gotcha").Fail(t, reason)

	// Fail reason:    failure reason
	// Type:           Got:string
	// Actually got:   "got string"

	if !strings.Contains(res, "Test name:") {
		t.Errorf("Expected to be contained the string `Test name:`, but not: %q", res)
	}
	if !strings.Contains(res, "TestFailGot, Gotcha") {
		t.Errorf("Expected to be contained the string `Test name:`, but not: %q", res)
	}

	if !strings.Contains(res, "Trace:") {
		t.Errorf("Expected to be contained the string `Trace:`, but not: %q", res)
	}

	if !strings.Contains(res, "Fail reason:") {
		t.Errorf("Expected to be contained the string `Fail reason:`, but not: %q", res)
	}

	if !strings.Contains(res, reason) {
		t.Errorf("Expected to be contained the string `%s`, but not: %q", reason, res)
	}

	if !strings.Contains(res, "Got:string") {
		t.Errorf("Expected to be contained type, but not: %q", res)
	}

	gotRe := regexp.MustCompile(fmt.Sprintf("Actually got:\\s*\\t%q", got))
	if gotRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", gotRe.String(), res)
	}
}

func TestFailGotExpect(t *testing.T) {
	stub()

	got := "got string"
	expect := "expect string"
	reason := "failure reason"
	Got(got).Expect(expect).Fail(t, reason)

	// Fail reason:    failure reason
	// Type:           Expect:string, Got:string
	// Expected:       "expect string"
	// Actually got:   "got string"

	if !strings.Contains(res, "Test name:") {
		t.Errorf("Expected to be contained the string `Test name:`, but not: %q", res)
	}
	if !strings.Contains(res, "Trace:") {
		t.Errorf("Expected to be contained the string `Trace:`, but not: %q", res)
	}

	if !strings.Contains(res, "Fail reason:") {
		t.Errorf("Expected to be contained the string `Fail reason:`, but not: %q", res)
	}

	if !strings.Contains(res, reason) {
		t.Errorf("Expected to be contained the string `%s`, but not: %q", reason, res)
	}

	if !strings.Contains(res, "Expect:string, Got:string") {
		t.Errorf("Expected to be contained types, but not: %q", res)
	}

	gotRe := regexp.MustCompile(fmt.Sprintf("Actually got:\\s*\\t%q", got))
	if gotRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", gotRe.String(), res)
	}

	expectRe := regexp.MustCompile(fmt.Sprintf("Expected:\\s*\\t%q", expect))
	if expectRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", expectRe.String(), res)
	}
}

func TestFailWithDiff(t *testing.T) {
	stub()

	w := New(ShowDiff, NotShowRaw)

	got := "a\nb\nc"
	expect := "a\nd\nc"
	reason := "not same string"
	w.Got(got).Expect(expect).Fail(t, reason)

	// Fail reason:    not same string
	// Type:           Expect:string, Got:string
	// Expected:       "a\nd\nc"
	// Actually got:   "a\nb\nc"
	// Diff details:   --- Expected
	//                 +++ Actually got
	//                 @@ -1,3 +1,3 @@
	//                  a
	//                 -d
	//                 +b
	//                  c

	if !strings.Contains(res, "Diff details:") {
		t.Errorf("Expected to be contained the string `Diff details:`, but not: %q", res)
	}

	diffRe := regexp.MustCompile("\\s*a\n\\s*-d\n\\s*\\+b\n\\s*c")
	if diffRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", diffRe.String(), res)
	}
}

func TestFailWithRawData(t *testing.T) {
	stub()

	w := New(NotShowDiff, ShowRaw)

	got := "a\nb\nc"
	expect := "a\nd\nc"
	reason := "not same string"
	w.Got(got).Expect(expect).Fail(t, reason)

	// Fail reason:    not same string
	// Type:           Expect:string, Got:string
	// Expected:       "a\nd\nc"
	// Actually got:   "a\nb\nc"
	// Raw Expect:     ---
	//                 a
	//                 d
	//                 c
	//                 ---
	// Raw Got:        ---
	//                 a
	//                 b
	//                 c
	//                 ---

	rawExpectRe := regexp.MustCompile("Raw Expect:\\s*\t---\n\\s*a\n\\s*d")
	if rawExpectRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", rawExpectRe.String(), res)
	}

	rawGotRe := regexp.MustCompile("Raw Got:\\s*\t---\n\\s*a\n\\s*b")
	if rawGotRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", rawGotRe.String(), res)
	}
}

func TestFailWithAdditionalMessage(t *testing.T) {
	stub()

	g := "a\nb\nc"
	e := "a\nd\nc"

	Got(g).Expect(e).Message("Example Label", "Some info").Fail(t, "Not same")

	// Fail reason:    Not same
	// Type:           Expect:string, Got:string
	// Expected:       "a\nd\nc"
	// Actually got:   "a\nb\nc"
	// Example Label:  Some info

	infoRe := regexp.MustCompile("Example Label:\\s*\tSome info\n")
	if infoRe.FindStringSubmatch(res) == nil {
		t.Errorf("Not matched the regexp `%s` for %q", infoRe.String(), res)
	}
}
