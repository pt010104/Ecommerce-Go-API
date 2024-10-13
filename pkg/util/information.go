package util

import (
	"regexp"
)

func IsValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\d{10}$`)
	return phoneRegex.MatchString(phone)
}
