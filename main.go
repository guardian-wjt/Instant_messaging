package main

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zserge/lorca"
	"os"
	"os/signal"
	"syscall"
)

//go:embed server/frontend/dist/*
var FS embed.FS

func main() { // main函数是程序的入口点
	go func() { // gin 协程
		gin.SetMode(gin.DebugMode)
		router := gin.Default()
		router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		router.Run(":8080") // listen and serve on 0.0.0.0:8080
	}()

	// 初始化Lorca用户界面
	var ui lorca.UI
	var err error
	// 使用Lorca创建一个新的UI实例，指定初始页面为百度首页
	// 参数包括页面大小和Chrome浏览器的启动参数
	ui, err = lorca.New("http://127.0.0.1:8080/", "", 800, 600, "--disable-sync", "--disable-translate", "--remote-allow-origins=*")
	if err != nil {
		// 如果创建UI时发生错误，打印错误信息并退出程序
		fmt.Println("Failed to create UI:", err)
		os.Exit(1)
	}

	// 创建一个用于接收系统信号的通道
	chSignals := make(chan os.Signal, 1)
	// 注册SIGINT和SIGTERM信号，以便在接收到这些信号时能够优雅地关闭UI
	signal.Notify(chSignals, syscall.SIGINT, syscall.SIGTERM)

	// 使用select语句监听UI的完成事件和系统信号
	select {
	case <-ui.Done():
		// 如果UI完成关闭，则正常退出程序
	case <-chSignals:
		// 如果接收到系统信号，则关闭UI后退出程序
	}
	// 显式关闭UI，以确保所有资源被正确释放
	ui.Close()
}
