package finder

import (
	"errors"
	"regexp"
)

func CutJsonp(regex, jsonp string) (string, error) {
	infoRegex := regexp.MustCompile(regex)
	slices := infoRegex.FindStringSubmatch(jsonp)
	if len(slices) < 1 {
		return "", errors.New("invalid jsonp")
	}
	return slices[0], nil
}
