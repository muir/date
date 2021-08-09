package date_test

import (
	"testing"
	"time"

	"github.com/muir/date"
	"github.com/stretchr/testify/assert"
)

func TestBasics(t *testing.T) {
	d := date.MustFromString("2010-11-12")
	assert.Equal(t, "2010-11-12", d.String(), "string")
	assert.Equal(t, "2010-11-13", d.AddDate(0, 0, 1).String(), "add day")
	assert.Equal(t, "2013-09-13", d.AddDate(3, -2, 1).String(), "add ymd")
	assert.Equal(t, 1, d.Sub(d.AddDate(0, 0, -1)), "one day diff")
	assert.Equal(t, d, date.FromJD(d.JD()), "round trip to int")
	assert.False(t, d.IsZero(), "is zero")

	mt, err := d.MarshalText()
	assert.NoError(t, err, "marshal")
	assert.Equal(t, []byte("2010-11-12"), mt, "marshal")

	d2, err := date.Parse("01/02/06", "08/06/12")
	if assert.NoError(t, err, "parse") {
		assert.Equal(t, "2012-08-06", d2.String(), "string2")
	}

	d3 := date.MustParse("2006/01/02", "2013/09/04")
	assert.Equal(t, "2013-09-04", d3.String())

	assert.Equal(t, "9/4/13", d3.Format("1/2/06"), "format")

	assert.Panics(t, func() { _ = date.MustParse("01/02/06", "2009-03-04") }, "must not parse")

	assert.Panics(t, func() { _ = date.MustFromString("01/02/06") }, "must not datefromstring")

	now := time.Now()
	assert.Equal(t, now.Format("2006-01-02"), date.FromTime(now).String(), "todays date")

	var dp date.Date
	if assert.NoError(t, (&dp).Scan(int64(20210304)), "scan int") {
		assert.Equal(t, "2021-03-04", dp.String(), "scanned int date")
	}
	if assert.NoError(t, (&dp).Scan(float64(20220714.2)), "scan float") {
		assert.Equal(t, "2022-07-14", dp.String(), "scanned float date")
	}
	if assert.NoError(t, (&dp).Scan("2024-11-02"), "scan string") {
		assert.Equal(t, "2024-11-02", dp.String(), "scanned string date")
	}
	if assert.NoError(t, (&dp).Scan(now), "scan time") {
		assert.Equal(t, now.Format("2006-01-02"), dp.String(), "scanned time date")
	}
	if assert.NoError(t, (&dp).Scan(nil), "scan nil") {
		assert.Equal(t, date.FromJD(0), dp, "scanned nil date")
	}
	if assert.NoError(t, (&dp).Scan([]byte("2024-11-12")), "scan bytes") {
		assert.Equal(t, "2024-11-12", dp.String(), "scanned bytes date")
	}

	v, err := dp.Value()
	if assert.NoError(t, err, "value") {
		assert.Equal(t, "2024-11-12", v, "value")

		var u date.Date
		if assert.NoError(t, (&u).UnmarshalText([]byte("2024-11-12")), "unmarshal") {
			assert.Equal(t, u.String(), "2024-11-12", "unmarshal")
		}
	}
}

func TestZero(t *testing.T) {
	zd := date.Zero
	assert.Equal(t, zd, date.FromJD(0), "from jd")
	assert.Equal(t, 0, zd.JD(), "to jd")
	assert.Equal(t, zd, date.FromTime(time.Time{}), "from time")
	assert.True(t, zd.Time().IsZero(), "to time")
	assert.Equal(t, zd, date.MustFromString(""), "from string")
	assert.Equal(t, "", zd.String())
	assert.Equal(t, zd.Format("2006/01/02"), time.Time{}.Format("2006/01/02"), "format")
	zdv, err := zd.Value()
	assert.NoError(t, err, "value error")
	assert.Nil(t, zdv, "zero value")
	assert.True(t, zd.IsZero(), "is zero")
}
