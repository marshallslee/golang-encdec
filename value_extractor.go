package main

import (
	"regexp"
	"strings"
)

// ExtractNameValue extracts the value of __name__ parameter from the input string.
func ExtractNameValue(input string) string {
	re := regexp.MustCompile(`__name__="(.*?)"`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// ExtractValueByParameter extracts the value of a given parameter from the input string.
// It handles both = and =~ operators and strips leading/trailing .* patterns.
func ExtractValueByParameter(input, parameterName string) string {
	re := regexp.MustCompile(parameterName + `=~?"(.*?)"`)
	match := re.FindStringSubmatch(input)
	if len(match) > 1 {
		value := match[1]
		value = strings.TrimPrefix(value, ".*")
		value = strings.TrimSuffix(value, ".*")
		return value
	}
	return ""
}
