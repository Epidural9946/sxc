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

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	systray.Run(onReady, onExit)
	//type Config struct {
	//	DefaultConfig struct {
	//		Ignore []string `yaml:"ignore"`
	//		Path   string   `yaml:"path"`
	//	} `yaml:"default"`
	//	Accounts []struct {
	//		Token  string   `yaml:"token"`
	//		Name   []string `yaml:"name"`
	//		Ignore []string `yaml:"ignore"`
	//	} `yaml:"accounts"`
	//}
	//
	//path, err := os.Getwd()
	//util.CheckError(err)
	//file, err := os.ReadFile(filepath.Join(path, "1config.yml"))
	//util.CheckError(err)
	//var c Config
	//err = yaml.Unmarshal(file, &c)
	//util.CheckError(err)
	//logger.Println(c)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("Bug 小工具")
	systray.SetTooltip("Bug 小工具")
	mQuit := systray.AddMenuItem("关闭", "关闭")
	mQuit.SetIcon(icon.Data)
	token, path := getConfig()
	service.Token = token
	service.Listen(path, service.PushPlusExec)
	// 监听退出菜单项的点击事件
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func getConfig() (string, string) {
	type Config struct {
		Token string `yaml:"token"`
		Path  string `yaml:"path"`
	}
	path, err := os.Getwd()
	util.CheckError(err)
	file, err := os.ReadFile(filepath.Join(path, "config.yml"))
	util.CheckError(err)
	var c Config
	err = yaml.Unmarshal(file, &c)
	util.CheckError(err)
	return c.Token, c.Path
}

func onExit() {
	// clean up here
}
