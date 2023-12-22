package core

import (
	"JishouSchedule/core/tasks/spider"
	"JishouSchedule/core/tools/config"
)

var log = config.Log

// Login 登录
func (web *WebObject) Login() bool {
	loginPage := web.Page

	log.Info("进入登录页面")
	// 导航到目标页面
	loginPage.MustNavigate(spider.LoginPageURL)

	// 等待页面加载完毕
	loginPage.MustWaitLoad()
	loginPage.MustWaitStable().MustScreenshot(spider.ImgHomePage)

	// 填写登录信息(账号|密码)
	loginPage.MustElement(spider.InputLogin).MustInput(config.Conifg.GetString("UserName"))
	loginPage.MustElement(spider.InputLoginTwo).MustInput(config.Conifg.GetString("PassWord"))

	// 点击立即登录按钮
	loginPage.MustElement(spider.LoginButton).MustClick()

	// 检验是否登录成功
	status, _ := loginPage.Element(spider.LoginStatus)
	if text, _ := status.Text(); text != "登录成功" {
		log.Error("登录失败: ", text)
		loginPage.MustScreenshot(spider.ImgFailed)
		return false
	}

	// 获取基本数据
	log.Info("登录成功")
	loginPage.MustElement(spider.WeatherPage).MustWaitStable().MustScreenshot(spider.ImgWeather)
	loginPage.MustElement(spider.PeopleInfo).MustWaitStable().MustScreenshot(spider.ImgPeopleInfo)
	loginPage.MustScreenshot(spider.ImgSuccess)
	return true
}

// Logout 退出
func (web *WebObject) Logout() {
	log.Info("正在退出登录...")

	// 找到下拉框元素
	search, _ := web.Page.Search("设置")
	search.First.MustClick()

	// 退出
	search, _ = web.Page.Search("退出")
	parent, _ := search.First.Parent()
	parentText, _ := parent.Text()
	parent.MustClick()

	log.Infof("成功%s登录!", parentText)

	defer web.Browser.MustClose()
}
