package utils

import "strings"

func IsBlankStr(str string) bool {
	return len(strings.TrimSpace(str)) == 0
}
