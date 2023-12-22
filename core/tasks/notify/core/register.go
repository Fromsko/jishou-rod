package core

import (
	"JishouSchedule/core/tasks/notify"
	"JishouSchedule/core/tools/common"
	"JishouSchedule/core/tools/config"
)

type (
	Server struct {
		*Service
		TemplateID         string
		RegisterServerList []RegisterServer
	}
	RegisterServer struct {
		NickName string // 订阅者 别名
		UserID   string // 订阅者 ID
	}
	Option func(*Server)
)

func (receiver RegisterServer) Notify(schedule Service) {
	schedule.Template.ToUser = receiver.UserID
	// 推送
	schedule.SendMsg(schedule.TempInfo, func(resp string) {
		config.Log.Info("❤️ 成功推送给 => " + receiver.NickName)
	})
}

func WithTemplateID(templateID string) Option {
	return func(server *Server) {
		server.TemplateID = templateID
	}
}

func WithFlower(flowers ...string) Option {
	return func(server *Server) {
		if server.AccessToken == "" {
			server.GetToken()
		}

		if len(flowers) == 0 {
			flowerList, _ := server.GetFlowerList()
			flowers = flowerList.Data.OpenID
		}

		for _, flower := range flowers {
			flower := RegisterServer{
				NickName: flower,
				UserID:   flower,
			}
			server.RegisterServerList = append(
				server.RegisterServerList, flower,
			)
		}
	}
}

func NewRegister(opts ...Option) (server *Server, observer *ScheduleObserver) {
	server = &Server{Service: new(Service)}
	observer = &ScheduleObserver{}

	for _, opt := range opts {
		opt(server)
	}

	server.Service.Template = &TemplateMessage{
		TemplateID: server.TemplateID,
		Data:       make(map[string]any),
		Url:        notify.CnameImage + common.GetWeek(36),
	}

	for _, user := range server.RegisterServerList {
		observer.Subscribe(user)
	}

	return server, observer
}
