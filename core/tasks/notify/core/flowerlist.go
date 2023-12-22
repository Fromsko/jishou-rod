package core

import (
	"JishouSchedule/core/tasks/notify"
	"JishouSchedule/core/tools/config"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Fromsko/gouitls/knet"
	"github.com/tidwall/gjson"
)

type (
	Flower struct {
		Total int        `json:"total"` // 关注个数
		Count int        `json:"count"` // 拉取个数
		Data  FlowerData `json:"data"`  // 详细数据
	}
	FlowerData struct {
		OpenID     []string `json:"openid"`
		NextOpenID string   `json:"next_openid"` // 列表最后一个ID
	}
)

func (s *Service) GetFlowerList() (flowers *Flower, err error) {
	Spider := knet.SendRequest{
		FetchURL: fmt.Sprintf(
			notify.FlowerList,
			s.AccessToken,
		),
	}
	Spider.Send(func(resp []byte, cookies []*http.Cookie, err error) {
		statusCode := gjson.GetBytes(resp, "errcode").Int()
		if statusCode == 40013 || err != nil {
			config.Log.Error("无效AppID")
			return
		}
		flowers = new(Flower)
		_ = json.Unmarshal(resp, flowers)
	})
	return flowers, err
}
