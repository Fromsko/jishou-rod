# **课表绘制(Go-后端)**

<div align="center">

📅 **吉首大学个人课表绘制**

🛠️ 简化课程表查询和生成课表图片的流程

---

</div>

## 📑 功能特点

- 自动化课程表数据爬取 - `rod`
- 可自定义课表样式和颜色 - `gg`
- 支持定时任务，定期 [12/h] 更新课表 - `cron`
- 提供简单的接口，方便集成到其他应用中 - `gin`

---

`配套应用`

- [前端应用](https://github.com/Fromsko/JishouSchedule/tree/main/frontend)：前端页面数据展示
- [微信推送](https://github.com/Fromsko/JishouSchedule/tree/main/notify)：推送数据到微信测试号

## 📦 安装

首先，确保您的 `Go` 版本为 **1.18** 或更高。然后，执行以下步骤来安装项目：

## 🚀 快速开始

1. 填写配置文件 `./config.json`

   ```json
   {
     "username": "xxxx",
     "password": "xxxx"
   }
   ```

2. 构建应用程序 `main.py`

   ```bash
   cd GoSchedule
   go mod tidy
   go build -ldflags "-w -s" -o Jishou-Rod.exe .
   ```

3. 启动程序

   ```bash
   Jishou-Rod.exe
   ```

4. 访问项目文档：

   - 访问 `http://host:port` 查看主页

5. 绘制个人课表：
   - 初次运行会自动获取一次数据
   - 后续则会自动执行定时任务

## 🔗 示例请求

### 获取课表数据

- **请求地址:** `http://host:prot/api/v1/get_cname_data?week=1`
- **请求方式:** `GET`
- **返回格式:** `json`

### 获取课表图片

- **请求地址:** `http://host:prot/api/v1/get_cname_table?week=1`
- **请求方式:** `GET`
- **返回格式:** `image`

## 🙏 鸣谢

> 欢迎提交 Pull Requests 和 Issues！

感谢以下开源项目，它们为本项目的开发提供了重要支持：

- [Colly](https://github.com/gocolly/colly): 🕸️ 用于网页爬取的 Go 语言库。
- [Cobra](https://github.com/spf13/cobra): 🐍 用于创建强大的命令行界面的 Go 语言库。
- [ants](https://github.com/panjf2000/ants): 🐜 一个高性能和低成本的 Go 语言协程池。
- [progressbar](https://github.com/schollz/progressbar): 📊 用于 Go 语言中的进度条显示。
- [gg](github.com/fogleman/gg): 🖼️ 用于图像处理的 Go 库。
- [cron](github.com/robfig/cron/v3): 🕒 用于定时任务的 Go 库。
- [Rod](https://github.com/go-rod/rod): 🌐 用于自动化网页爬取的工具。
- [Gin](https://fastapi.tiangolo.com/): 🚀 用于构建高性能 API 的 Web 框架。

## ©️ 许可

此项目遵循 [MIT 许可](LICENSE)。
