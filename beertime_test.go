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

	// Test for beertime on monday before beertime week.
	date, _ = time.Parse(dateFormat, "2020-11-16 11:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(4*24+4), Duration(date).Hours())

	// Test for beertime on saturday 15:00.
	date, _ = time.Parse(dateFormat, "2020-11-21 15:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(6*24), Duration(date).Hours())

	// Test for beertime on thursday 15:00.
	date, _ = time.Parse(dateFormat, "2020-11-12 15:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(1*24), Duration(date).Hours())

	// Test for beertime on friday 15:00.
	date, _ = time.Parse(dateFormat, "2020-11-13 15:00")
	assert.Equal(t, true, Now(date))
	assert.Equal(t, float64(0), Duration(date).Hours())

	// Test for beertime on saturday 00:00.
	date, _ = time.Parse(dateFormat, "2020-11-14 00:00")
	assert.Equal(t, false, Now(date))
	assert.Equal(t, float64(6*24+15), Duration(date).Hours())
}
