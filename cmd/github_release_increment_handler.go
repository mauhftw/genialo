package cmd

import (
	// System
	"fmt"
	"os"
	"strconv"
	"strings"

	// 3rd Party
	log "github.com/sirupsen/logrus"
)

// Get future release
func GetFutureRelease(o *ReleaseOptions, t string) {

	release := map[string]int{
		"major": 0,
		"minor": 1,
		"patch": 2,
	}

	currentRelease := o.FutureRelease
	intRelease, releaseDigits := ReleaseToInt(currentRelease)

	switch release[t] {
	case 0:
		newRelease := IncrementRelease(currentRelease, intRelease, release[t], releaseDigits)
		o.FutureRelease = newRelease
	case 1:
		newRelease := IncrementRelease(currentRelease, intRelease, release[t], releaseDigits)
		o.FutureRelease = newRelease
	case 2:
		newRelease := IncrementRelease(currentRelease, intRelease, release[t], releaseDigits)
		o.FutureRelease = newRelease
	default:
		log.Info("Wrong option, Try again")
	}
}

// Convert release to int
func ReleaseToInt(currentRelease string) ([]int, []string) {

	// Parse release digits major.minor.patch by .
	releaseDigits := strings.Split(currentRelease, ".")

	// Check that follows x.x.x pattern
	if len(releaseDigits) > 3 {
		log.Errorf("Your tag %s, does not follow the semver pattern x.x.x", currentRelease)
		os.Exit(1)
	}

	outputRelease := []int{}

	// Convert release digits from string to int
	for _, digit := range releaseDigits {
		aux, err := strconv.Atoi(digit)
		if err != nil {
			log.Errorf("\t %v", err)
			log.Fatalf("\t Release %v not supported", currentRelease)
		}
		outputRelease = append(outputRelease, aux)
	}

	return outputRelease, releaseDigits
}

// Increments the value of the release
func IncrementRelease(currentRelease string, intRelease []int, release int, releaseDigits []string) string {

	// Increments the value
	intRelease[release]++
	increment := strconv.Itoa(intRelease[release])
	releaseDigits[release] = increment

	// Format the new release
	newRelease := fmt.Sprintf("%s.%s.%s", releaseDigits[0], releaseDigits[1], releaseDigits[2])
	log.Infof("\t New release is: %v", newRelease)

	return newRelease
}
