package validator

import "regexp"

func IsValidUrl(s string) bool {
	pattern := `^https:\/\/[a-zA-Z0-9\-_]+(\.[a-zA-Z0-9\-_]+)+(\/[^\s]*)?$`
	isUrl, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}

	return isUrl
}
