package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/logger"
	"web_app/routers"
	"web_app/settings"

	"go.uber.org/zap"
)

//@title web_app
//@version 0.0.1
//@description Go Web 开发通用脚手架模板
//@termsOfService http://swagger.io/terms/
//
//@contact.name author：@ouzhsh
//@contact.url http://www.swagger.io/support
//@contact.email support@swagger.io
//
//@license.name Apache 2.0
//@license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
//@host 127.0.0.1:8080
//@BasePath /
func main() {
	// 1. 加载配置文件
	// 通过命令行输入配置文件位置
	if len(os.Args) < 2 {
		panic("请输入执行程序的配置文件...")
	}
	// os.Args[1]获取输入的文件 ./conf/config.yaml
	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Println("init settings failed, err:", err)
		return
	}
	fmt.Println("init settings success...")

	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Println("init logger failed, err:", err)
		return
	}
	fmt.Println("init logger success...")

	// 3. 初始化MySQL
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init mysql failed, err:", err)
		return
	}
	fmt.Println("init mysql success...")
	defer mysql.Close() // 程序退出关闭数据库连接

	// 4. 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init redis failed, err:", err)
		return
	}
	fmt.Println("init redis success...")
	defer redis.Close()

	// 5. 注册路由
	r := routers.Setup(settings.Conf.Mode)

	// 6.启动服务 (优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	// 开启一个goroutine启动服务
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error(fmt.Sprintf("listen: %s", err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error(fmt.Sprintf("Server Shutdown: %s", err))
	}
	zap.L().Info("Server exiting")
}
