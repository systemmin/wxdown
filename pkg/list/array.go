package list

import "strings"

func IsContain(array []string, content string) bool {
	for _, s := range array {
		if strings.Contains(s, content) {
			return true
		}
	}
	return false
}

func IsExist(array []string, str string) bool {
	for _, s := range array {
		if str == s {
			return true
		}
	}
	return false
}

func IsEmpty(array []string) bool {
	return len(array) == 0
}
