package obj

import (
	"fmt"
	"testing"
)

type Example struct {
	id   int
	name string
}

func TestTruncate(t *testing.T) {
	o := NewObjectWithMaxLen(&Example{id: 12, name: "John Doe"}, 10)

	tts := []struct {
		format string
		expect string
	}{
		{format: "%v", expect: "&{12 John <... truncated>"},
		{format: "%+v", expect: "&{id:12 na<... truncated>"},
		{format: "%#v", expect: "&obj.Examp<... truncated>"},
		{format: "%s", expect: "&{%!s(int=<... truncated>"},
		{format: "%Y", expect: "*obj.Example"}, // The type of RawValue()
	}
	for _, tt := range tts {
		if got := fmt.Sprintf(tt.format, o); got != tt.expect {
			t.Errorf("got:%s != expect:%s", got, tt.expect)
		}
	}
}

func TestRawValue(t *testing.T) {
	o := NewObject("John Doe")
	if o.AsRawValue() != "John Doe" {
		t.Error("RawValue() was wrong")
	}
}

func TestIsStringType(t *testing.T) {
	if o := NewObject("John Doe"); !o.IsStringType() {
		t.Error("IsStringType() was wrong")
	}
	if o := NewObject(7); o.IsStringType() {
		t.Error("IsStringType() was wrong")
	}
}

func TestIsDumpableRawType(t *testing.T) {
	if o := NewObject([]int{1, 2}); !o.IsDumpableRawType() {
		t.Error("IsDumpableRawType() was wrong")
	}
	if o := NewObject(7); o.IsDumpableRawType() {
		t.Error("IsDumpableRawType() was wrong")
	}
}

func TestDump(t *testing.T) {
	o := NewObject(123)
	if o.AsDumpString() != "(int) 123\n" {
		t.Error("Dump() was wrong")
	}
}
