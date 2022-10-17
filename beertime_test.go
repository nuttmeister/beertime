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
	// Test each second from end of beertime until next beertime.
	for i := 0; i < 572399; i++ {
		date, _ := time.Parse(dateFormat, "2020-11-14 00:00")
		date = date.Add(time.Duration(i) * time.Second)
		assert.Equal(t, false, Now(date))
		assert.Equal(t, float64((6*24+15)*3600-i), Duration(date).Seconds())
	}

	// Test each second during beertime.
	for i := 0; i < 32400; i++ {
		date, _ := time.Parse(dateFormat, "2020-11-13 15:00")
		date = date.Add(time.Duration(i) * time.Second)
		assert.Equal(t, true, Now(date))
	}
}
