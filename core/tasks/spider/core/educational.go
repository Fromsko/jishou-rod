package core

import (
	"JishouSchedule/core/tasks/spider"
	"fmt"

	rod "github.com/Fromsko/rodPro"
)

// NextPage 教务系统页面
func (web *WebObject) NextPage() {
	log.Info("正在寻找教务系统入口...")
	search, _ := web.Page.Search("教务系统（师生入口）")
	search.First.MustClick()

	// 新页面
	web.Page.MustWaitLoad()
	pages, _ := web.Browser.Pages()
	web.Page = pages.Last()
	web.DownPage = pages.First()
	web.DownPage.MustScreenshot(spider.ImgSchedule)
}

func (web *WebObject) Extract(CallBack func(downloadPage *rod.Page, selectName string)) {
	log.Info("正在提取数据")
	web.DownPage.MustElementR("a", "学期理论课表").MustClick()
	target := web.DownPage.MustWaitStable()
	targetInfo, _ := target.Info()

	log.Info("标题: ", targetInfo.Title)

	for i := 1; i <= 20; i++ {
		selectElem := target.MustElement("#zc")
		// 使用正则表达式选择包含特定周数的选项
		selectName := fmt.Sprintf(`第%d周`, i)
		_ = selectElem.Select(
			[]string{selectName},
			true,
			rod.SelectorTypeRegex,
		)
		target.MustWaitLoad()
		CallBack(target, selectName)
	}

	log.Info("准备关闭开启的页面!")
	defer web.DownPage.MustClose()
}

// Html 提取源码
func (web *WebObject) Html(downloadPage *rod.Page) (html string) {
	element, err := downloadPage.ElementX("/html/body/div[4]/div[2]/form[2]")
	if err != nil {
		log.Errorf("没找到课表数据: %s", err)
		return
	}

	html, err = element.HTML()
	if err != nil {
		log.Error("查找失败!")
	}

	return html
}
