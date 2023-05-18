package service

import (
	"embed"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"text/template"
	"zpaul.org/chd/sxc/util"
)

type data string

func (d *data) Write(p []byte) (n int, err error) {
	var s1 = string(p)
	d2 := data(s1)
	*d = *d + d2
	return 0, nil
}

type PushPlus struct {
	Token       string `json:"token"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Template    string `json:"template"`
	Channel     string `json:"channel"`
	Webhook     string `json:"webhook"`
	CallbackUrl string `json:"callbackUrl"`
	Timestamp   string `json:"timestamp"`
}

//go:embed pushplus.html
var htmlContent embed.FS

func PushPlusExec(token string, message util.XCAutoLog) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(err) // 将 interface{} 转型为具体类型。
		}
	}()
	var d data = ""
	pushPlusT, err := htmlContent.ReadFile("pushplus.html")
	util.CheckError(err)
	t, err := template.New("CHD").Parse(string(pushPlusT))
	util.CheckError(err)
	err = t.ExecuteTemplate(&d, "CHD", message)
	util.CheckError(err)
	body := PushPlus{Token: token, Title: message.Name, Content: string(d), Template: "html", Channel: "wechat"}
	marshal, _ := json.Marshal(body)
	request, _ := http.NewRequest("POST", "https://www.pushplus.plus/send", strings.NewReader(string(marshal)))
	request.Header.Add("Content-Type", "application/json")
	do, _ := http.DefaultClient.Do(request)
	rBody, _ := io.ReadAll(do.Body)
	logger.Println(string(rBody))
}
