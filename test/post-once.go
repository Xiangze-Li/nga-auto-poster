package main

import (
	"log"

	"github.com/Xiangze-Li/nga-auto-poster/config"
	"github.com/Xiangze-Li/nga-auto-poster/poster"
	"github.com/Xiangze-Li/nga-auto-poster/utils"
)

const initMsgTemplate = `自动回帖启动
目标帖子: https://bbs.nga.cn/read.php?tid=%d
内容文件: %s
分隔符  : %s
定时器  : %s
`

func main() {
	config := config.LoadConfig("config.yaml")

	log.Printf(
		initMsgTemplate,
		config.Tid, config.ContentFile, config.Split, config.Cron,
	)
	err := poster.PostReply(config)
	utils.ExitOnError(err)
}
