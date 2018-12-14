package bows

import (
	"encoding/json"
	"testing"
)

type PayLoad struct {
	User
	Price string `json:"price"`
}

type User struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
}

func TestTm(t *testing.T) {

	byt, _ := json.Marshal(PayLoad{Name: "fhy", Pass: "123", Price: "100"})

	t.Log(string(byt))
}
