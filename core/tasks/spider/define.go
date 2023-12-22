package spider

import (
	"JishouSchedule/core/tools/common"
)

const (
	// LoginPageURL 登录页面
	LoginPageURL = "https://webvpn.jsu.edu.cn/https/77726476706e69737468656265737421e0f6528f693a7b45300d8db9d6562d/#/UserLogin?sn=ELNmlR8NQXCY1q-kHEz7xA&client_id=lFkRfDkwSW6z3IDAZpZo3g&redirect_uri=https%3A%2F%2Fwebvpn.jsu.edu.cn%2Flogin%3Foauth_login%3Dtrue"

	// InputLogin 输入组件
	InputLogin    = "#pane-loginPwd > div > div.login-frame-left-below-div1-user > div > div > div > input"
	InputLoginTwo = "#pane-loginPwd > div > div.login-frame-left-below-div1-password > div > div > div > input"
	LoginButton   = "#login > div.login-frame > div > div > div.login-frame-left-below-div2 > div > button"

	// LoginStatus 登录状态
	LoginStatus = "body > div.el-message"

	// PeopleInfo 主页面信息
	PeopleInfo  = "#app > div > div:nth-child(2) > div:nth-child(1) > div > div.w50.mr20.pdt10 > div > div > div.fsb.mt10 > div"
	WeatherPage = "#app > div > div:nth-child(2) > div:nth-child(1) > div > div.w50.mr20.pdt10 > div > div > div.fsb.mt10 > iframe"
)

var (
	OutputPath    = common.WorkPath("cache", "output")
	DataPath      = common.WorkPath("cache", "data")
	ImgPath       = common.WorkPath("cache", "img")
	ImgHomePage   = common.GenPath(ImgPath, "home-page.png")
	ImgFailed     = common.GenPath(ImgPath, "Login-Failed.png")
	ImgSuccess    = common.GenPath(ImgPath, "Login-Success.png")
	ImgWeather    = common.GenPath(ImgPath, "Weather.png")
	ImgPeopleInfo = common.GenPath(ImgPath, "PeopleInfo.png")
	ImgSchedule   = common.GenPath(ImgPath, "Schedule.png")
)
