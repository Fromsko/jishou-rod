package core

import (
	"fmt"
	"os"
	"strings"

	rod "github.com/Fromsko/rodPro"
	"github.com/Fromsko/rodPro/lib/launcher"
	"github.com/PuerkitoBio/goquery"
)

type ConfigOption func(*WebObject)
type PaserOption func() (doc *goquery.Document, err error)

// WebObject 浏览器对象
type WebObject struct {
	baseUrl  string
	wsUrl    string
	Page     *rod.Page
	DownPage *rod.Page
	Browser  *rod.Browser
}

// CnameObject 课程总数据
type CnameObject struct {
	CnameResult  map[string]any `json:"课程信息"`
	CnameSpecial string         `json:"备注"`
	Cname        string         `json:"班级"`
	Weekly       string         `json:"周次"`
}

func WithWebSoket(ws string) ConfigOption {
	return func(web *WebObject) {
		web.wsUrl = ws
	}
}

func WithBaseUrl(url string) ConfigOption {
	return func(web *WebObject) {
		web.baseUrl = url
	}
}

// InitWeb 初始化浏览器
func InitWeb(opts ...ConfigOption) (Web *WebObject) {
	var (
		u string
		b *rod.Browser
		w = &WebObject{
			Browser: b,
		}
	)

	for _, opt := range opts {
		opt(w)
	}

	if w.wsUrl != "" {
		b = rod.New().ControlURL(w.wsUrl).MustConnect()
	} else {
		if path, exists := launcher.LookPath(); exists {
			u = launcher.New().Bin(path).MustLaunch()
		}
		log.Infof("链接成功 %s", u)
		b = rod.New().ControlURL(u).MustConnect()
	}
	w.Page = b.MustPage()
	w.Browser = b
	return w
}

func InitCname(option PaserOption) (cname *CnameObject, doc *goquery.Document) {
	doc, err := option()
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}

	return &CnameObject{
		CnameSpecial: "",
		CnameResult:  map[string]any{},
	}, doc
}

func ReadHTML(content string) PaserOption {
	reader := strings.NewReader(content)

	return func() (doc *goquery.Document, err error) {
		doc, err = goquery.NewDocumentFromReader(reader)
		if err != nil {
			return nil, fmt.Errorf("解析 HTML 文档时出错: %s", err)
		}
		return doc, nil
	}
}

func ReadFile(fileName string) PaserOption {
	html, err := os.ReadFile(fileName)
	if err != nil {
		log.Errorf("无法读取HTML文件: %s", err)
	}

	return func() (doc *goquery.Document, err error) {
		return ReadHTML(string(html))()
	}
}
