package web

import (
	"strings"
)

func titler(str string) string {
	str = strings.ToLower(strings.Replace(str, "_", " ", -1))
	str = strings.Replace(str, ".", " ", -1)
	return strings.Title(str)
}
