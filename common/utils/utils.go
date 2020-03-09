package utils

import "strings"

func IsStrBlank(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
