package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/corrots/go-crontab/common"
	"github.com/corrots/go-crontab/master/handler"
	"github.com/corrots/go-crontab/master/manager"
	"github.com/gin-gonic/gin"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	conf = flag.String("c", "./conf.yaml", "specific the service config file")
)

func main() {
	flag.Parse()
	initSystem()
	if err := common.InitConfig(conf); err != nil {
		log.Println(err)
		return
	}
	if viper.GetBool("app.debug") {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		gin.DefaultWriter = io.MultiWriter(os.Stdout)
	}

	if err := manager.InitJobManager(); err != nil {
		log.Printf("init job manager err: %v\n", err)
		return
	}

	r := gin.Default()
	index := r.Group("/index")
	{
		r.LoadHTMLGlob("master/static/views/*")
		index.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index.tpl", gin.H{
				"title": "Cron Job Management",
			})
		})
	}

	job := r.Group("/job")
	{
		job.GET("/load/:jobName", handler.Load)
		job.GET("/list", handler.List)
		job.POST("/save", handler.Save)
		job.DELETE("/:jobName", handler.Delete)
		job.POST("/kill/:jobName", handler.Kill)
	}
	r.Run(viper.GetString("api.port"))
}

func initSystem() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
