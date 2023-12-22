package web

import (
	"JishouSchedule/core/tasks"
	"JishouSchedule/core/tools/common"
	"JishouSchedule/core/tools/config"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

var (
	router *gin.Engine
	server *http.Server
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	config.Log.Info("程序启动成功🚀")
	config.Log.Info("当前版本: " + tasks.Version)
	config.Log.Info("项目地址: https://github.com/Fromsko/Jishouschedule")
}

// AutoTask 自动任务
func AutoTask(Timer string, Task func()) {
	c := cron.New(cron.WithSeconds())

	// 每天早晨7:00
	if _, err := c.AddFunc(Timer, Task); err != nil {
		config.Log.Debugf("添加任务时出错：%v", err)
		return
	}

	c.Start()
}

// StartServer 启动Web服务
func StartServer(port string) {
	// 创建Gin路由
	router = gin.Default()
	router.NoRoute(func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 404, "msg": "页面不存在"})
	})
	router.Use(Cors())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"code": 200, "msg": "主页"})
	})
	api := router.Group("/api/v1")
	{
		api.GET("/get_cname_data", getCnameData)
		api.GET("/get_cname_table", getCnameTable)
	}

	server = &http.Server{
		Addr:    port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			config.Log.Errorf("Failed to start server: %s", err)
		}
	}()

	config.Log.Infof("成功启动服务 => http://localhost%s", port)
}

// RestartServer 重启服务
func RestartServer(port string) {
	// 关闭当前服务器
	if server != nil {
		if err := server.Shutdown(context.TODO()); err != nil {
			fmt.Printf("Failed to shutdown server: %s", err)
		}
	}
	// 启动新的服务器
	StartServer(port)
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func getCnameData(c *gin.Context) {
	weekStr := c.Query("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week parameter"})
		return
	}

	// 读取指定目录下的JSON文件
	data, err := readJSONData(week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read JSON data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func getCnameTable(c *gin.Context) {
	weekStr := c.Query("week")
	week, err := strconv.Atoi(weekStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid week parameter"})
		return
	}

	// 生成图片并返回
	imagePath, err := generateImage(week)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate image"})
		return
	}

	c.File(imagePath)
}

func readJSONData(week int) (result map[string]any, err error) {
	search := fmt.Sprintf("第%d周", week)
	err = common.ReadFilesWithCallback(
		common.GenPath("cache", "data"),
		search,
		func(filePath string) error {
			content, _ := os.ReadFile(filePath)
			_ = json.Unmarshal(content, &result)

			if result == nil {
				result = map[string]any{
					"search": search,
					"error":  "No data found for this week.",
				}
			}

			return nil
		},
	)

	return result, err
}

func generateImage(week int) (imagePath string, err error) {
	search := fmt.Sprintf("第%d周", week)
	err = common.ReadFilesWithCallback(
		common.GenPath("cache", "output"),
		search,
		func(filePath string) (err error) {
			if filePath == "" {
				return errors.New("没找到")
			}
			imagePath = filePath
			return nil
		},
	)
	return imagePath, err
}
