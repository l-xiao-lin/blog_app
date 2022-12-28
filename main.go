package main

import (
	"blog_app/controller"
	"blog_app/dao/mysql"
	"blog_app/dao/redis"
	"blog_app/logger"
	"blog_app/pkg/snowflake"
	"blog_app/routers"
	"blog_app/settings"
	"flag"
	"fmt"
)

func main() {
	var filePath string
	flag.StringVar(&filePath, "f", "./conf/config.yaml", "配置文件路径")
	flag.Parse()
	fmt.Println(filePath)
	//1、加载配置文件
	if err := settings.Init(filePath); err != nil {
		fmt.Printf("load config failed,err:%v\n", err)
		return
	}
	//2、日志配置
	err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode)
	if err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
	}

	//3、加载mysql
	if err = mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	defer mysql.Close() //程序退出关闭数据库连接

	//4、加载redis
	if err = redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	//5、初始化gin框架内置的校验器使用的翻译器
	if err = controller.InitTrans("zh"); err != nil {
		fmt.Printf("InitTrans failed,err:%v\n", err)
		return

	}
	//6、加载雪花算法
	if err = snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineId); err != nil {
		fmt.Printf("snowflake.Init failed,err:%v\n", err)
		return
	}

	//7、注册路由
	r := routers.SetupRouter()

	r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
}
