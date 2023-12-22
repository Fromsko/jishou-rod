package main

import (
	"JishouSchedule/core/interfaces"
	"JishouSchedule/core/tasks"
	"JishouSchedule/core/tools/common"
	"JishouSchedule/core/web"
)

func main() {
	if common.FristRun() {
		interfaces.SpiderRod()
	} else {
		web.StartServer(tasks.DefaultPort)
	}

	interfaces.PushWechat()
	go web.AutoTask("0 0 */12 * * ?", interfaces.SpiderRod)
	go web.AutoTask("0 0 7 * * ?", interfaces.PushWechat)
	select {}
}
