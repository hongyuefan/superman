package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bitly/go-simplejson"
)

type PayLoad struct {
	User
	Price string `json:"price"`
}

type User struct {
	Name  string `json:"name"`
	Pass  string `json:"pass"`
	Emain string
}

func TestTm(t *testing.T) {

	byt, _ := json.Marshal(PayLoad{User: User{Name: "fhy", Pass: "123"}, Price: "100"})

	user := new(User)

	json.Unmarshal(byt, user)

	t.Log(string(byt), user)
}

func Test_simplejs(t *testing.T) {
	ss := `[{"data":{"result":"false","error_code":"20104"},"channel":"ok_sub_futurusd_ltc_ticker_this_week"}]`
	js, err := simplejson.NewJson([]byte(ss))
	if err != nil {
		fmt.Println("error:", err)
	}
	arr, err := js.Array()
	if err != nil {
		fmt.Println("error:", err)
	}
	for i := 0; i < len(arr); i++ {
		fmt.Println("data:", arr[i])
	}
}
