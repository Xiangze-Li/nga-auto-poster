package main

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/Xiangze-Li/nga-auto-poster/config"
	"github.com/Xiangze-Li/nga-auto-poster/poster"
	"github.com/Xiangze-Li/nga-auto-poster/utils"
	rl "github.com/lestrrat-go/file-rotatelogs"
	"github.com/robfig/cron/v3"
)

const initMsgTemplate = `自动回帖启动
目标帖子: https://bbs.nga.cn/read.php?tid=%d
内容文件: %s
分隔符  : %s
定时器  : %s
`

func init() {
	rOut, err := rl.New(
		"nga-auto-poster.%Y_%m_%d.log",
		rl.WithRotationTime(time.Hour*24),
		rl.WithMaxAge(-1),
		rl.WithRotationCount(3),
	)
	utils.ExitOnError(err)
	log.SetFlags(log.Ldate | log.Ltime)
	log.SetOutput(io.MultiWriter(os.Stderr, rOut))
}

func main() {
	config := config.LoadConfig("config.yaml")

	c := cron.New()
	_, err := c.AddFunc(config.Cron, func() {
		err := poster.PostReply(config)
		if err != nil {
			log.Println(err)
		}
	})
	utils.ExitOnError(err)

	log.Printf(
		initMsgTemplate,
		config.Tid, config.ContentFile, config.Split, config.Cron,
	)

	c.Run()
}
