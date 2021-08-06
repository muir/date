
# date - dates as Julian Days with methods to make them easy to work with

[![GoDoc](https://godoc.org/github.com/muir/date?status.png)](https://pkg.go.dev/github.com/muir/date)

Install:

	go get github.com/muir/date

---

## Date

Go has a very nice time package.  Using times for dates is a bit
awkward for a number of reasons.

Using YYYY-MM-DD strings as dates is easy to work read, but hard
to manipulate.

Using Julian Days is efficient and date math is super easy, but they're
not human readable.

This package makes Julian Days human readable by providing a `.String()` 
method and by implementing encoding.TextMarshaler and encoding.TextUnmarshaler
so that dates look great when used as map keys and as values.

Database convience is provided by implementations of sql.Valuer and sql.Scanner 

