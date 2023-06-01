package main

import (
	"github.com/getlantern/systray"
	"github.com/getlantern/systray/example/icon"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"zpaul.org/chd/sxc/service"
	"zpaul.org/chd/sxc/util"
)

type Config struct {
	DefaultConfig struct {
		Ignore     []string `yaml:"ignore"`
		Path       string   `yaml:"path"`
		Token      string   `yaml:"token"`
		DatabaseId string   `yaml:"databaseId"`
	} `yaml:"default"`
	Accounts []struct {
		Token  string   `yaml:"token"`
		Name   []string `yaml:"name"`
		Ignore []string `yaml:"ignore"`
	} `yaml:"accounts"`
}

var version string
var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	logger.Infof("Ver: %s", version)
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Bug 小工具")
	systray.SetTooltip("Bug 小工具")
	mQuit := systray.AddMenuItem("关闭", "关闭")
	mQuit.SetIcon(icon.Data)
	c := getConfig()
	if c.Accounts == nil || len(c.Accounts) == 0 {
		systray.Quit()
	}
	for _, account := range c.Accounts {
		for _, name := range account.Name {
			ignore := append(c.DefaultConfig.Ignore, account.Ignore...)
			strName := util.HexToStr(name)
			logger.Infof("Account: %s, Token: %s, ignore: %s", strName, account.Token, ignore)
			service.Listen(filepath.Join(c.DefaultConfig.Path, name), ignore, func(message util.XCAutoLog) {
				message.Account = strName
				service.PushPlusExec(c.DefaultConfig.Token, c.DefaultConfig.DatabaseId, message)
			})
		}
	}
	// 监听退出菜单项的点击事件
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func getConfig() Config {
	path, e := os.Getwd()
	util.CheckError(e)
	file, e := os.ReadFile(filepath.Join(path, "config.yml"))
	util.CheckError(e)
	var c Config
	e = yaml.Unmarshal(file, &c)
	util.CheckError(e)
	return c
}

func onExit() {
	// clean up here
}
