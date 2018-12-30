package database

import (
	"testing"
)

func TestRegist(t *testing.T) {
	RegistDB()
}

func TestAdd(t *testing.T) {

	id, err := AddKLine_5Min(&KLine_5Min{Open: "1", High: "2", Low: "0", Close: "1", Deal: "10", Time: "11"})

	t.Log(id, err)

	id, err = SetKLine_5MinByTime(&KLine_5Min{Open: "1111", High: "2", Low: "0", Close: "1", Deal: "10", Time: "11"})

	t.Log(id, err)

	id, err = SetKLine_5MinByTime(&KLine_5Min{Open: "1", High: "2", Low: "0", Close: "1", Deal: "10", Time: "22"})

	t.Log(id, err)
}
