package api

import "JishouSchedule/core/tools/common"

// InitTemplateMessage 模板数据
func InitTemplateMessage() map[string]any {
	cnameInfo := GetCnameData()
	weather := SearchWeather("吉首")
	onesay := GetEveryDay()
	cnameInfo["Week"] = map[string]string{"value": common.GetWeekly()}
	cnameInfo["City"] = map[string]string{"value": weather.Local}
	cnameInfo["Weather"] = map[string]string{"value": weather.WeatherInfo.Text}
	cnameInfo["Temp"] = map[string]string{"value": weather.WeatherInfo.Temp}
	cnameInfo["Onesay"] = map[string]string{"value": onesay}
	return cnameInfo
}
