package api

import (
	"JishouSchedule/core/tools/config"
	"net/http"

	"github.com/Fromsko/gouitls/knet"
	"github.com/tidwall/gjson"
)

// GetEveryDay 获取每日一句
func GetEveryDay() string {
	var equiangular string
	Spider := knet.SendRequest{
		FetchURL: "http://open.iciba.com/dsapi/?date",
	}
	Spider.Send(func(resp []byte, cookies []*http.Cookie, err error) {
		if err != nil {
			config.Log.Error("获取每日一句失败")
			equiangular = "千里之堤, 始于足下。"
			return
		}
		equiangular = gjson.Get(string(resp), "note").String()
	})
	return equiangular
}
