package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"panda/arithmetic"
	"time"
)

type ReqMobile struct {
	Mobile     string `json:"mobile"`
	NationCode string `json:"nationcode"`
}
type ReqMsg struct {
	Ext    string    `json:"ext`
	Extend string    `json:"extend"`
	Params []string  `json:"params"`
	Sig    string    `json:"sig"`
	Tel    ReqMobile `json:"tel"`
	Time   int64     `json:"time"`
	TplId  int       `json:"tpl_id"`
}

type RspMsg struct {
	Result int    `json:"result"`
	ErrMsg string `json:"errmsg"`
	Ext    string `json:"ext"`
	Fee    int    `json:"fee"`
	Sid    string `json:"sid"`
}

func SigMsg(mobile, appKey, sRand, sTime string) string {

	var bySum []byte

	bySum = make([]byte, 32)

	strSha := "appkey=" + appKey + "&random=" + sRand + "&time=" + sTime + "&mobile=" + mobile

	Sum := sha256.Sum256([]byte(strSha))

	copy(bySum[:], Sum[:])

	return hex.EncodeToString(bySum)
}

func SendMsg(appId, appKey, nation, mobile string, params []string, tplId int) (err error) {

	sRand := arithmetic.GetRandLimit(4)

	nowTime := time.Now().Unix()

	sig := SigMsg(mobile, appKey, sRand, fmt.Sprintf("%v", nowTime))

	reqMsg := ReqMsg{
		Params: params,
		Sig:    sig,
		Tel: ReqMobile{
			Mobile:     mobile,
			NationCode: nation,
		},
		Time:  nowTime,
		TplId: tplId,
	}

	return MsgPostReq(appId, sRand, reqMsg)
}

func MsgPostReq(appId, sRand string, reqMsg ReqMsg) (err error) {

	reqUrl := "https://yun.tim.qq.com/v5/tlssmssvr/sendsms?sdkappid=" + appId + "&random=" + sRand

	reqBody, err := json.Marshal(reqMsg)
	if err != nil {
		return
	}

	body, err := Post(reqBody, reqUrl)
	if err != nil {
		return
	}
	buffer := new(bytes.Buffer)

	io.Copy(buffer, body.Body)

	var rspMsg RspMsg

	if err = json.Unmarshal(buffer.Bytes(), &rspMsg); err != nil {
		return
	}

	if 0 == rspMsg.Result {
		return
	}

	err = fmt.Errorf("%v", rspMsg.ErrMsg)

	return
}

func Post(data []byte, url string) (body *http.Response, err error) {

	client := &http.Client{}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	return client.Do(req)
}
