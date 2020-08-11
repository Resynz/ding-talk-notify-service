/**
 * @Author: Resynz
 * @Date: 2020/8/4 10:07
 */
package structs

type DingTalkConfig struct {
	Token       string   `json:"token"`
	AesKey      string   `json:"aes_key"`
	Key         string   `json:"key"`
	CallBackTag []string `json:"call_back_tag"`
	Url         string   `json:"url"`
}

type RegisterMap map[string]string

type DingTalkResponse struct {
	MsgSignature string `json:"msg_signature"`
	Encrypt      string `json:"encrypt"`
	TimeStamp    string `json:"timeStamp"` // 这里钉钉要求返回的格式也太丑陋了!
	Nonce        string `json:"nonce"`
}

type Register struct {
	InstanceId string `json:"instance_id"`
	NotifyUrl  string `json:"notify_url"`
}

type EventType string

type NotifyTask struct {
	NotifyUrl string `json:"notify_url"`
	Body      []byte `json:"body"`
}
