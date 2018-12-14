package utils

import (
	"strings"
)

func ParseStringToArry(src, step string) []string {
	return strings.SplitN(src, step, -1)
}
