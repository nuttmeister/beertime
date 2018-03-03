// musslanbeertime returns current date and time in ISO-B33R format.
// and number of nano seconds until ISO-B33R is true (if it's currently false).

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// MusslanBeerTime represents current time according to ISO-B33R and how many hours until next beer time
// if it's currently not beertime. Otherwise it will be 0.
type MusslanBeerTime struct {
	CurrentBeerTime       bool  `json:"currentBeerTime"`
	UniversalBeerConstant int   `json:"universalBeerConstant"`
	NanoToBeerTime        int64 `json:"nanoToBeerTime"`
}

var (
	headers = map[string]string{
		"content-type": "application/json",
	}
)

// GetMusslanBeerTime will return current Musslan Beer Time according to ISO-B33R.
// Returns MusslanBeerTime
func GetMusslanBeerTime(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Grab the Beer Constant from env.
	beerConst, err := getBeerConstant()
	if err != nil {
		log.Printf("%s\n", err.Error())
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}

	// Get current UTC time.
	now := time.Now()
	now = now.UTC()

	// Check if it's beer time!
	beerTime := isItBeerTime(&now, beerConst)

	if beerTime {
		return events.APIGatewayProxyResponse{Body: convertBeerTimeToString(&MusslanBeerTime{CurrentBeerTime: true, UniversalBeerConstant: beerConst, NanoToBeerTime: 0}), StatusCode: 200, Headers: headers}, nil
	}

	// If it's not musslan beer time just return false.
	return events.APIGatewayProxyResponse{Body: convertBeerTimeToString(&MusslanBeerTime{CurrentBeerTime: false, UniversalBeerConstant: beerConst, NanoToBeerTime: getHoursUntilBeerTime(&now, beerConst)}), StatusCode: 200, Headers: headers}, nil
}

// getBeerConstant gets the current beer constant. beer constant should be set to either 0 or 1 depending on what
// weeks beer time should be on. This can be seen as the space time constant.
// In case of an error we will return 0 and the error.
// 0 = beer time is on even weeks.
// 1 = beer time is on odd weeks.
func getBeerConstant() (int, error) {
	beerConst, err := strconv.Atoi(os.Getenv("BEERCONST"))
	if err != nil {
		return 0, err
	}

	switch {
	case beerConst > 1:
		return 0, fmt.Errorf("Beer Constant was set, but it was larger than 1 which is not permitted")
	case beerConst < 0:
		return 0, fmt.Errorf("Beer Constant was set, but it was smaller than 0 which is not permitted")
	}

	return beerConst, nil
}

// isItBeerWeek returns true if current week is a beer week.
func isItBeerWeek(week int, beerConst int) bool {
	if week%2 == beerConst {
		return true
	}
	return false
}

// isItBeerTime will check if it's currently beer time or not. Return true for when it's currently beer time.
func isItBeerTime(now *time.Time, beerConst int) bool {
	_, week := now.ISOWeek()

	if isItBeerWeek(week, beerConst) {
		switch {
		case now.Weekday() == time.Friday && now.Hour() > 13:
			return true
		case now.Weekday() == time.Saturday && now.Hour() < 6:
			return true
		}
	}

	return false
}

// getHoursUntilBeerTime will return number of nano seconds to next beer time.
func getHoursUntilBeerTime(now *time.Time, beerConst int) int64 {
	next := *now

	// Get which week next beertime will be in.
	days := getNumberOfDaysUntilNextBeerTime(now, beerConst)
	remaining := getRemainingNanosecsOfDay(now)

	// Add the static hour of beer time start (14:00 UTC).
	// Add number of full days until next beer time.
	// It will not include the current day and friday, since those should be calculated with a duration instead.
	// Add the reamning time of the current day.
	next = next.Add(time.Duration(14) * time.Hour)
	next = next.Add(time.Duration(days*24) * time.Hour)
	next = next.Add(remaining)

	return next.Sub(*now).Nanoseconds()
}

// getRemainingNanosecsOfDay returns the number of nano seconds remaining for the current day.
func getRemainingNanosecsOfDay(now *time.Time) time.Duration {
	dur := time.Duration(now.Hour()) * time.Hour
	dur += time.Duration(now.Minute()) * time.Minute
	dur += time.Duration(now.Second()) * time.Second
	dur += time.Duration(now.Nanosecond()) + time.Nanosecond
	return (time.Duration(24) * time.Hour) - dur
}

// getNumberOfDaysUntilNextBeerTime returns the number of full days until next beer time.
func getNumberOfDaysUntilNextBeerTime(now *time.Time, beerConst int) int {
	_, week := now.ISOWeek()
	if isItBeerWeek(week, beerConst) {
		return 7 + getNumberOfDaysUntilFriday(now)
	}

	return getNumberOfDaysUntilFriday(now)
}

// getNumberOfDaysUntilFriday returns the number of days until friday from current day.
func getNumberOfDaysUntilFriday(now *time.Time) int {
	days := 0

	switch now.Weekday() {
	case time.Monday:
		days = 3
	case time.Tuesday:
		days = 2
	case time.Wednesday:
		days = 1
	case time.Thursday:
		days = 0
	case time.Friday:
		days = -1
	case time.Saturday:
		days = 5
	case time.Sunday:
		days = 4
	}

	return days
}

// convertBeerTimeToString will json marshal the MusslanBeerTime struct into a JSON string.
func convertBeerTimeToString(beertime *MusslanBeerTime) string {
	bytes, _ := json.Marshal(beertime)
	return string(bytes)
}

func main() {
	lambda.Start(GetMusslanBeerTime)
}
