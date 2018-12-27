package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

func getFutureRelease(o *ReleaseOptions, t string) {

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
		fmt.Println("error, try again")
	}
}

// convert release to int
func ReleaseToInt(currentRelease string) ([]int, []string) {

	// parse release digits major.minor.patch by .
	releaseDigits := strings.Split(currentRelease, ".")
	outputRelease := []int{}

	// convert release digits from string to int
	for _, digit := range releaseDigits {
		aux, _ := strconv.Atoi(digit)
		outputRelease = append(outputRelease, aux)
	}

	return outputRelease, releaseDigits
}

// increments the value of the release
func IncrementRelease(currentRelease string, intRelease []int, release int, releaseDigits []string) string {

	// increments the value
	intRelease[release]++
	increment := strconv.Itoa(intRelease[release])
	releaseDigits[release] = increment

	// format the new release
	newRelease := fmt.Sprintf("%s.%s.%s", releaseDigits[0], releaseDigits[1], releaseDigits[2])
	return newRelease
}
