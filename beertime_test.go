package beertime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	dateFormat = "2006-01-02 15:04"
)

func TestBeerTime(t *testing.T) {

	// Test for beertime on friday 16:00. Right week.
	date, _ := time.Parse(dateFormat, "2020-11-20 16:00")
	assert.Equal(t, true, Now(date))
	assert.Equal(t, float64(0), Duration(date).Hours())

	// Test for beertime on friday 23:59. Right week.
	date, _ = time.Parse(dateFormat, "2020-11-20 23:59")
	assert.Equal(t, true, Now(date))
	assert.Equal(t, float64(0), Duration(date).Hours())

	// Test for beertime on moday before beertime week.
	date, _ = time.Parse(dateFormat, "2020-11-16 11:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(4*24+5), Duration(date).Hours())

	// Test for beertime on saturday 00:00. Right week directly after beertime.
	date, _ = time.Parse(dateFormat, "2020-11-21 00:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(14*24-8), Duration(date).Hours())

	// Test for beertime on saturday 16:00. Right week directly after beertime.
	date, _ = time.Parse(dateFormat, "2020-11-21 16:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(13*24), Duration(date).Hours())

	// Test for beertime on thursday 16:00. But wrong week.
	date, _ = time.Parse(dateFormat, "2020-11-12 16:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(8*24), Duration(date).Hours())

	// Test for beertime on friday 16:00. But wrong week.
	date, _ = time.Parse(dateFormat, "2020-11-13 16:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(7*24), Duration(date).Hours())

	// Test for beertime on thursday 16:00. But wrong week.
	date, _ = time.Parse(dateFormat, "2020-11-14 00:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(6*24+16), Duration(date).Hours())
}
