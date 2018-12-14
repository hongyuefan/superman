package protocol

type OkexUserParam struct {
	ApiKey     string `json:"api_key"`
	SecretKey  string `json:"secret_key"`
	PassPhrase string `json:"pass_phrase"`
}
