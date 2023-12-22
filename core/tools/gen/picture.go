package gen

import (
	"JishouSchedule/core/tasks/spider"
	"JishouSchedule/core/tools/common"
	"fmt"
	"image/color"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/Fromsko/gouitls/logs"
	"github.com/fogleman/gg"
)

var log = logs.InitLogger()

// ImgOption 包含生成图片所需的各种选项
type ImgOption struct {
	FontName        string // 字体名称
	FontPath        string // 字体路径
	StoragePath     string // 存储路径
	TemplateImgPath string // 模板图片路径
	Size            struct {
		Width  int
		Height int
	}
	Image *gg.Context // 图片对象
}

type ImgOptionModifier func(*ImgOption)

// WithFontOrPath 根据参数类型设置字体名称或字体路径
func WithFontOrPath(fontOrPath string) ImgOptionModifier {
	return func(img *ImgOption) {
		if _, err := os.Stat(fontOrPath); os.IsNotExist(err) {
			img.FontName = fontOrPath
		} else {
			img.FontPath = fontOrPath
		}
	}
}

// WithSize 返回一个设置图片尺寸的 ImgOptionModifier
func WithSize(w, h int) ImgOptionModifier {
	return func(img *ImgOption) {
		img.Size.Width = w
		img.Size.Height = h
	}
}

// WithStoragePath 返回一个设置存储路径的 ImgOptionModifier
func WithStoragePath(path string) ImgOptionModifier {
	return func(img *ImgOption) {
		img.StoragePath = path
	}
}

// InitImg 根据提供的修改函数创建一个新的 ImgOption 实例
func InitImg(modifiers ...ImgOptionModifier) *ImgOption {
	// 默认设置
	op := &ImgOption{
		FontName:        "Deng",
		StoragePath:     "output",
		TemplateImgPath: "init_photo.png",
		Size: struct {
			Width  int
			Height int
		}{
			Width:  2000,
			Height: 1000,
		},
	}

	op.Image = gg.NewContext(
		op.Size.Width,
		op.Size.Height,
	)

	for _, modifier := range modifiers {
		modifier(op)
	}

	op.InitFont()
	return op
}

// InitFont 方法用于初始化字体
func (o *ImgOption) InitFont() {
	// 如果字体路径不为空，则加载指定路径的字体
	if o.FontPath != "" {
		err := o.Image.LoadFontFace(o.FontPath, 30)
		if err != nil {
			fmt.Println("Error loading font:", err)
		}
		return
	}

	// 根据平台设置字体路径
	fontPath := func(platform string) (fontPath string) {
		switch platform {
		case "darwin":
			fontPath = "/Library/Fonts/" + o.FontName + ".ttf"
		case "linux":
			fontPath = "/etc/fonts/" + o.FontName + ".ttf"
		case "windows":
			fontPath = "C:\\Windows\\Fonts\\" + o.FontName + ".ttf"
		default:
			fontPath = "" // 其他平台暂不支持自动加载字体
		}
		o.FontPath = fontPath
		return fontPath
	}(runtime.GOOS)

	// 检查字体文件是否存在
	if !common.ExistFile(fontPath) {
		log.Errorf("字体文件未找到! %s", fontPath)
	}

	// 加载字体
	err := o.Image.LoadFontFace(fontPath, 30)
	if err != nil {
		fmt.Println("Error loading font:", err)
	}
}

// CreateBasePhoto 方法用于创建基础模板图片
func (o *ImgOption) CreateBasePhoto() *gg.Context {
	// 创建画布
	o.Image.SetRGB(1, 1, 1) // 设置背景颜色为白色，不包含透明度
	o.Image.Clear()         // 清除背景

	// 画格子
	o.Image.SetRGB(0, 0, 0) // 设置线条颜色为黑色
	o.Image.DrawLine(250, 0, 250, 1000)
	for i := 500; i <= 1750; i += 250 {
		o.Image.DrawLine(float64(i), 0, float64(i), 890)
	}
	for _, i := range []int{100, 235, 375, 515, 655, 750, 890} {
		o.Image.DrawLine(0, float64(i), 2000, float64(i))
	}
	o.Image.Stroke()

	// 备注
	o.Image.SetRGB(47/255.0, 79/255.0, 79/255.0) // 设置文字颜色
	o.Image.DrawStringAnchored("备注", 95, 950, 0, 0)
	o.Image.SetRGB(0, 0, 0) // 恢复线条颜色为黑色
	o.Image.DrawStringAnchored("晚上", 95, 710, 0, 0)

	// 节次
	o.Image.DrawStringAnchored("第 1-2 节", 70, 160, 0, 0)
	o.Image.DrawStringAnchored("第 3-4 节", 70, 295, 0, 0)
	o.Image.DrawStringAnchored("第 5-6 节", 70, 435, 0, 0)
	o.Image.DrawStringAnchored("第 7-8 节", 70, 575, 0, 0)
	o.Image.DrawStringAnchored("第 9-10-11 节", 38, 800, 0, 0)

	// 时间
	o.Image.DrawStringAnchored("08:00—09:40", 45, 200, 0, 0)
	o.Image.DrawStringAnchored("10:10—11:50", 45, 335, 0, 0)
	o.Image.DrawStringAnchored("14:30-16:10", 45, 475, 0, 0)
	o.Image.DrawStringAnchored("16:20-18:00", 45, 615, 0, 0)
	o.Image.DrawStringAnchored("19:00-20:40", 45, 850, 0, 0)

	// 星期
	days := []string{"星期一", "星期二", "星期三", "星期四", "星期五", "星期六", "星期天"}
	dateNum := []float64{330, 580, 830, 1080, 1330, 1580, 1830}
	o.Image.SetRGB(0, 0, 0) // 恢复文字颜色为黑色
	for i, text := range days {
		o.Image.DrawStringAnchored(text, dateNum[i], 60, 0, 0)
	}

	return o.Image
}

type SaveToDB struct {
	FileName string
	meta     *gg.Context
}

// TODO: 随着 db.go 一同移除
func (s *SaveToDB) Write(p []byte) (n int, err error) { return }

func (s *SaveToDB) save() (err error) {
	fileBase := filepath.Base(s.FileName)
	if err = s.meta.SavePNG(s.FileName); err != nil {
		log.Errorf("保存失败: %s", fileBase)
	} else {
		log.Infof("保存成功: %s", fileBase)
	}
	return err
}

// SaveImg 存储图片
func SaveImg(path string, img *gg.Context) (string, error) {
	stb := &SaveToDB{
		FileName: path,
		meta:     img,
	}
	stb.meta.EncodePNG(stb)
	return stb.FileName, stb.save()
}

// CheckData 检查数据
func (o *ImgOption) CheckData(data map[string]any) (dataClass map[string]any) {
	photoPath := common.GenPath(spider.ImgPath, o.TemplateImgPath)
	if !common.ExistFile(photoPath) {
		img := o.CreateBasePhoto()
		SaveImg(photoPath, img)
	}

	if data != nil {
		dataClass = map[string]any{
			"data":       data,
			"init_photo": photoPath,
		}
		return dataClass
	} else {
		log.Error("传入数据为空")
	}
	return nil
}

// Create 方法用于创建图片
func (o *ImgOption) Create(cnameData map[string]any) (string, *gg.Context) {
	weekPlaceData := []float64{250, 500, 750, 1000, 1250, 1500, 1750}
	hPlaceData := map[string]float64{"第一大节": 100, "第二大节": 235, "第三大节": 375, "第四大节": 515, "第五大节": 750}
	colorList := []color.RGBA{
		{251, 255, 242, 200},
		{192, 192, 192, 200},
		{255, 255, 0, 200},
		{244, 164, 95, 200},
		{127, 255, 0, 200},
		{218, 112, 214, 200},
		{156, 147, 133, 200},
		{186, 164, 48, 200},
		{15, 56, 154, 200},
		{49, 65, 196, 200},
		{153, 51, 250, 200},
		{34, 139, 34, 200},
		{255, 192, 203, 200},
		{255, 127, 80, 200},
		{237, 145, 33, 200},
	}

	// 颜色数据 和 组合数据
	colorData := make(map[string]color.RGBA)
	originData := o.CheckData(cnameData)

	// 截取需要的数据
	cnameData = originData["data"].(map[string]any)
	photoPath := originData["init_photo"].(string)
	cname := cnameData["班级"].(string)
	week := cnameData["周次"].(string)

	// 创建画板
	img, _ := gg.LoadImage(photoPath)
	context := gg.NewContextForImage(img)
	_ = context.LoadFontFace(o.FontPath, 30)

	for weekDay, dayName := range cnameData["课程信息"].(map[string]interface{})["星期数据"].([]interface{}) {
		dayInfo := cnameData["课程信息"].(map[string]interface{})["课程数据"].(map[string]interface{})[dayName.(string)].(map[string]interface{})
		for timeSlot, courseInfo := range dayInfo {
			if courseInfo == "没课哟" {
				continue
			}

			lessonName := courseInfo.(map[string]interface{})["课程名"].(string)
			teacher := courseInfo.(map[string]interface{})["老师"].(string)
			weekInfo := courseInfo.(map[string]interface{})["周次"].(string)
			place := courseInfo.(map[string]interface{})["教室"].(string)

			if _, exists := colorData[lessonName]; !exists {
				colorIndex := rand.Intn(len(colorList))
				colorData[lessonName] = colorList[colorIndex]
				colorList = append(colorList[:colorIndex], colorList[colorIndex+1:]...)
			}

			colorInfo := colorData[lessonName]

			context.SetRGB(float64(colorInfo.R)/255.0, float64(colorInfo.G)/255.0, float64(colorInfo.B)/255.0)
			context.DrawRectangle(weekPlaceData[weekDay], hPlaceData[timeSlot], 249, 140)
			context.FillPreserve()
			context.Stroke()

			context.SetRGB(0, 0, 0)

			match := func(dt string, str string) bool {
				m, _ := regexp.MatchString(str, dt)
				return m
			}
			if match(lessonName, `（网络）`) {
				lessonName = strings.Replace(lessonName, `（网络）`, "网络-", 1)
			}

			context.DrawStringAnchored(lessonName, weekPlaceData[weekDay], hPlaceData[timeSlot]+30, 0, 0)
			context.DrawStringAnchored(teacher, weekPlaceData[weekDay]+5, hPlaceData[timeSlot]+65, 0, 0)
			context.DrawStringAnchored(weekInfo, weekPlaceData[weekDay]+5, hPlaceData[timeSlot]+95, 0, 0)
			context.DrawStringAnchored(place, weekPlaceData[weekDay]+5, hPlaceData[timeSlot]+125, 0, 0)
		}
	}

	context.SetRGB(0, 0, 0)
	context.DrawStringAnchored(cname, 70, 40, 0, 0)
	context.DrawStringAnchored(week, 80, 80, 0, 0)
	context.DrawStringAnchored(cnameData["备注"].(string), 300, 955, 0, 0)

	savePath := common.GenPath(
		spider.OutputPath,
		fmt.Sprintf("%s%s课表.png", cname, week),
	)
	return savePath, context
}
