package core

import (
	"JishouSchedule/core/tasks/spider"
	"JishouSchedule/core/tools/common"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Course 结构表示每一节课的信息
type Course struct {
	CourseName string `json:"课程名"`
	Teacher    string `json:"老师"`
	Weeks      string `json:"周次"`
	Classroom  string `json:"教室"`
}

// 初始化星期和节次数据切片
var weekdays = []string{"星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期日"}
var sections = []string{"第一大节", "第二大节", "第三大节", "第四大节", "第五大节"}

func (c *CnameObject) Resolve(doc *goquery.Document) {
	// 初始化课程表
	courseTable := make(map[string]map[string]any)
	selectRegex := "#kbtable tbody tr:nth-child(%d) td:nth-child(%d)"

	for i, weekday := range weekdays {
		courseTable[weekday] = make(map[string]any)
		for j, section := range sections {
			courseTable[weekday][section] = "没课哟"

			// 获取课程信息
			cell := doc.Find(fmt.Sprintf(selectRegex, j+2, i+2))

			// 如果课程信息非空，则解析课程信息并填充到课程表中
			if strings.TrimSpace(cell.Text()) != "" {
				course, err := c.ParseCourseHTML(cell)
				if err == nil {
					courseTable[weekday][section] = course
				}
			}
		}
	}

	c.CnameSpecial = doc.Find(
		fmt.Sprintf(selectRegex, 7, 2),
	).Text()

	c.CnameResult = map[string]any{
		"课程数据": courseTable,
		"星期数据": weekdays,
		"节次数据": sections,
	}
}

// ParseCourseHTML parseCourseHTML 解析课程信息字符串
func (c *CnameObject) ParseCourseHTML(courseInfo *goquery.Selection) (Course, error) {
	var course Course

	html, _ := courseInfo.Html()

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return course, err
	}

	doc.Find("br").ReplaceWithHtml("|")

	courseText := doc.Find("div").Last().Text()
	courseData := strings.Split(courseText, "|")

	course.CourseName = courseData[0]
	course.Teacher = doc.Find("font[title='老师']").First().Text()
	course.Weeks = doc.Find("font[title='周次(节次)']").First().Text()
	course.Classroom = doc.Find("font[title='教室']").First().Text()

	return course, nil
}

// writeJSONToFile 将 JSON 数据写入文件
func (c *CnameObject) writeJSONToFile(fileName string, data []byte) error {
	saveFile := common.GenPath(
		spider.DataPath,
		fileName+".json",
	)

	file, err := os.Create(saveFile)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	return nil
}

func (c *CnameObject) Marshal(selectName string) (result *map[string]any) {
	c.Weekly = selectName
	marshal, _ := json.Marshal(c)
	_ = json.Unmarshal(marshal, &result)
	return result
}

func (c *CnameObject) WriteFile(fileName string, cname string) (result map[string]any) {
	c.Cname = cname
	c.Weekly = fileName

	marshal, err := json.Marshal(c)
	if err != nil {
		return
	}

	if err = c.writeJSONToFile(fileName, marshal); err != nil {
		fmt.Println("写入 JSON 文件时出错:", err)
		return
	}

	_ = json.Unmarshal(marshal, &result)
	return result
}
