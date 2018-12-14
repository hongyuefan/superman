package utils

import (
	"testing"
)

func TestParseStringToArry(t *testing.T) {

	s := "okex,huobi,binity"

	t.Log(ParseStringToArry(s, ",")[0])
}
