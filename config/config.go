package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Xiangze-Li/nga-auto-poster/utils"
	"github.com/robfig/cron/v3"
	"gopkg.in/yaml.v3"
)

func LoadConfig(fileName string) Config {
	var config Config

	bConfig, err := os.ReadFile(fileName)
	utils.ExitOnError(err)
	err = yaml.Unmarshal(bConfig, &config)
	utils.ExitOnError(err)

	config.Split = strings.TrimSpace(config.Split)
	config.Cron = strings.TrimSpace(config.Cron)

	if !config.Verify() {
		utils.ExitOnError(fmt.Errorf("配置有误"))
	}

	sp := strings.Split(config.CookieString, ";")
	config.Cookies = make(map[string]string)
	for _, s := range sp {
		s = strings.TrimSpace(s)
		if len(s) == 0 {
			continue
		}
		sp2 := strings.Split(s, "=")
		if len(sp2) != 2 {
			continue
		}
		config.Cookies[strings.TrimSpace(sp2[0])] = strings.TrimSpace(sp2[1])
	}

	return config
}

type Config struct {
	CookieString string            `yaml:"cookie_string"`
	Cookies      map[string]string `yaml:"-"`
	Fid          int               `yaml:"fid"`
	Tid          int               `yaml:"tid"`
	ContentFile  string            `yaml:"content_file"`
	Split        string            `yaml:"split"`
	Cron         string            `yaml:"cron"`
}

func (v Config) Verify() bool {
	const rw = 0600

	if len(v.CookieString) == 0 ||
		v.Fid <= 0 ||
		v.Tid <= 0 ||
		len(v.ContentFile) == 0 ||
		len(v.Split) == 0 ||
		len(v.Cron) == 0 {
		log.Println("配置字段不完整")
		return false
	}

	fi, err := os.Stat(v.ContentFile)
	if err != nil {
		log.Printf("打开内容文件(%s)失败: %v\n", v.ContentFile, err)
		return false
	}

	if fi.Mode()&rw != rw {
		log.Println("内容文件权限不足")
		return false
	}

	_, err = cron.ParseStandard(v.Cron)
	if err != nil {
		log.Printf("cron表达式(%s)不合法: %v\n", v.Cron, err)
		return false
	}

	return true
}
