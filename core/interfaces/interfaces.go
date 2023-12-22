package interfaces

import (
	"JishouSchedule/core/tasks"
	"JishouSchedule/core/tasks/notify"
	"JishouSchedule/core/tasks/notify/api"
	notifyCore "JishouSchedule/core/tasks/notify/core"

	"JishouSchedule/core/tasks/spider"
	spiderCore "JishouSchedule/core/tasks/spider/core"
	"JishouSchedule/core/tools/config"
	"JishouSchedule/core/tools/gen"

	ginWeb "JishouSchedule/core/web"

	"context"
	"time"

	rod "github.com/Fromsko/rodPro"
)

func PushWechat() {
	server, observer := notifyCore.NewRegister(
		notifyCore.WithFlower(),
		notifyCore.WithTemplateID(notify.TemplateID),
	)

	// 获取数据
	server.TempInfo = api.InitTemplateMessage()
	//// 推送任务
	observer.PushSchedule(*server.Service)
}

func SpiderRod() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	done := make(chan bool)

	go func() {
		if spiderHandler(ctx) {
			done <- true
		}
	}()

	select {
	case <-done:
		config.Log.Info("数据更新成功!")
	case <-ctx.Done():
		config.Log.Warningf("更新数据失败!")
	}

	defer ginWeb.RestartServer(tasks.DefaultPort)
}

func spiderHandler(ctx context.Context) bool {
	web := spiderCore.InitWeb(
		spiderCore.WithBaseUrl(
			spider.LoginPageURL,
		),
	)

	Img := gen.InitImg()

	classTable := func(downloadPage *rod.Page, selectName string) {
		cname, doc := spiderCore.InitCname(
			spiderCore.ReadHTML(
				web.Html(downloadPage),
			),
		)

		cname.Resolve(doc)

		result := cname.WriteFile(
			selectName,
			tasks.ClassName,
		)

		gen.SaveImg(Img.Create(result))
	}

	defer func() {
		if err := recover(); err != nil {
			config.Log.Errorf("教务系统正在维护! (%s)", err)
		}
	}()

	select {
	case <-ctx.Done():
		web.Browser.MustClose()
		return false
	default:
		if loginStatus := web.Login(); loginStatus {
			web.NextPage()
			web.Extract(classTable)
			web.Logout()
			return true
		}
		return false
	}
}
