# witness

<a href="https://github.com/bayashi/witness/blob/main/LICENSE" title="witness License"><img src="https://img.shields.io/badge/LICENSE-MIT-GREEN.png" alt="MIT License"></a>
<a href="https://github.com/bayashi/witness/actions" title="witness CI"><img src="https://github.com/bayashi/witness/workflows/main/badge.svg" alt="witness CI"></a>
<a href="https://goreportcard.com/report/github.com/bayashi/witness" title="witness report card" target="_blank"><img src="https://goreportcard.com/badge/github.com/bayashi/witness" alt="witness report card"></a>
<a href="https://pkg.go.dev/github.com/bayashi/witness" title="Go witness package reference" target="_blank"><img src="https://pkg.go.dev/badge/github.com/bayashi/witness.svg" alt="Go Reference: witness"></a>

`witness` is a test helper to make an evident report on a fail of your test.

## Usage

Simple case.

```go
package main

import (
    "testing"

    w "github.com/bayashi/witness"
)

func TestExample(t *testing.T) {
    g := "a\nb\nc"
    e := "a\nd\nc"

    if g != e {
        w.Got(g).Expect(e).Fail(t, "Not same")
    }
}
```

below result will be shown:

```go
Test name:      TestExample
Trace:          /home/usr/go/src/github.com/bayashi/witness/witness_test.go:14
Fail reason:    Not same
Type:           Expect:string, Got:string
Expected:       "a\nd\nc"
Actually got:   "a\nb\nc"
```

It's able to show more additional info:

```go
w.Got(g).Expect(e).ShowMore().Fail(t, "Not same")
```

And then,

```go
Test name:      TestExample
Trace:          /home/usr/go/src/github.com/bayashi/witness/witness_test.go:14
Fail reason:    Not same
Type:           Expect:string, Got:string
Expected:       "a\nd\nc"
Actually got:   "a\nb\nc"
Diff details:   --- Expected
                +++ Actually got
                @@ -1,3 +1,3 @@
                 a
                -d
                +b
                 c
Raw Expect:     ---
                a
                d
                c
                ---
Raw Got:        ---
                a
                b
                c
                ---
```

So easy and fair, isn't this?

You shouldn't need to spend your time anymore to show fail report.

And, there is a shortcut function, direct `Fail`, instead of builder interface. Below lines are same.

```go
    w.Got(g).Expect(e).Fail(t, "Not same")
    w.Fail(t, "Not same", g, e)
```

See [Witness Package reference](https://pkg.go.dev/github.com/bayashi/witness) for more details.

## Installation

```cmd
go get github.com/bayashi/witness
```

## License

MIT License

## Author

Dai Okabayashi: https://github.com/bayashi

## See Also

https://github.com/bayashi/actually

## Special Thanks To

https://github.com/stretchr/testify
