package util

import (
	"fmt"
	"regexp"
)

func RegexProcessString(regex *regexp.Regexp, str string) (map[string]string, error) {
	// Verify the match first
	match := regex.FindStringSubmatch(str)
	if match == nil {
		return nil, fmt.Errorf("regex '%s' did not find any match on the given string '%s'", regex, str)
	}

	// Construct a map with the key-value pairs
	result := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("regex '%s' did not find any match on the given string '%s'", regex, str)

	} else {
		return result, nil
	}
}