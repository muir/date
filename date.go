package date

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"time"

	"github.com/carlosjhr64/jd"
	"github.com/pkg/errors"
)

// Date represents a date using a Julian Day that has been
// wrapped with methods to make it easy to use and easy to
// work with.
type Date int

func (d Date) String() string {
	return jd.ToDate(int(d))
}

func (d Date) Time() time.Time {
	year, month, day := jd.J2YMD(int(d))
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// Parse works the same as time.Parse but only pays attention
// to 2006 01 02
func Parse(format string, date string) (Date, error) {
	t, err := time.Parse(format, date)
	return DateFromTime(t), err
}

func MustParse(format string, date string) Date {
	d, err := Parse(format, date)
	if err != nil {
		panic("parse date: " + err.Error())
	}
	return d
}

// Format is just time.Format
func (d Date) Format(format string) string {
	return d.Time().Format(format)
}

// DateFromString parses dates in the format YYYY-MM-DD
func FromString(s string) (Date, error) {
	if len(s) != 10 {
		return 0, errors.Errorf("cannot convert '%s' to date", s)
	}
	if s[4] != '-' || s[7] != '-' {
		return 0, errors.Errorf("cannot convert '%s' to date", s)
	}
	y, err := strconv.Atoi(s[0:4])
	if err != nil {
		return 0, errors.Wrapf(err, "cannot convert '%s' to date", s)
	}
	m, err := strconv.Atoi(s[5:7])
	if err != nil {
		return 0, errors.Wrapf(err, "cannot convert '%s' to date", s)
	}
	d, err := strconv.Atoi(s[8:10])
	if err != nil {
		return 0, errors.Wrapf(err, "cannot convert '%s' to date", s)
	}
	return Date(jd.YMD2J(y, m, d)), nil
}

func MustFromString(s string) Date {
	d, err := FromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func DateFromTime(t time.Time) Date {
	return Date(jd.Number(t))
}

// AddDate is just like time.AddDate
func (d Date) AddDate(years int, months int, days int) Date {
	if years == 0 && months == 0 {
		return d + Date(days)
	}
	return DateFromTime(d.Time().AddDate(years, months, days))
}

func (d1 Date) Sub(d2 Date) int {
	return int(d1 - d2)
}

// Scan implements sql.Scanner for database reads
func (d *Date) Scan(src interface{}) error {
	switch t := src.(type) {
	case int64:
		// YYYYMMDD integer
		*d = Date(jd.YMD2J(int(t)/10000, (int(t)%10000)/100, int(t)%100))
		return nil
	case float64:
		// YYYYMMDD integer
		*d = Date(jd.YMD2J(int(t)/10000, (int(t)%10000)/100, int(t)%100))
		return nil
	case []byte:
		j, err := jd.ToNumber(string(t))
		if err != nil {
			return err
		}
		*d = Date(j)
		return nil
	case string:
		j, err := jd.ToNumber(t)
		if err != nil {
			return err
		}
		*d = Date(j)
		return nil
	case time.Time:
		*d = Date(jd.Number(t))
		return nil
	case nil:
		*d = 0
		return nil
	default:
		return fmt.Errorf("Scan: unable to scan type %T into UUID", src)
	}
}

// Value implements sql.Valuer
func (d Date) Value() (driver.Value, error) {
	return d.String(), nil
}

// MarshalText implements encoding.TextMarshaler so that dates look good
// in map keys
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler
func (d *Date) UnmarshalText(text []byte) error {
	var err error
	*d, err = FromString(string(text))
	return err
}
