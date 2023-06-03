package service

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
	"zpaul.org/chd/sxc/util"
)

var c = make([]string, 0)
var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func Listen(path string, ignore []string, exec func(message util.XCAutoLog)) {
	sc := make(chan string)
	go addAccountRootDirWatch(path, sc)
	go addAccountDateDirWatch(sc, ignore, exec)
}

// addDirWatch 监听账号下的 Log 根目录监听
func addAccountRootDirWatch(path string, sc chan<- string) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	util.CheckErrorF(err)
	err = watcher.Add(path)
	util.CheckErrorF(err)
	path = filepath.Join(path, time.Now().Format("2006-01-02"))
	logger.Info(path)
	if _, err2 := os.Stat(path); err2 == nil {
		sc <- path
	}
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) && !contains(c, event.Name) {
				info, err := os.Stat(event.Name)
				util.CheckError(err)
				if info.IsDir() {
					sc <- event.Name
					logger.Println("Add New Paths {}", event.Name)
				}

			}
		case err, ok := <-watcher.Errors:
			if ok {
				util.CheckError(err)
			}
		}
	}
}

func addAccountDateDirWatch(sc <-chan string, ignore []string, exec func(message util.XCAutoLog)) {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	defer watcher.Close()
	util.CheckErrorF(err)

	go func(watcher *fsnotify.Watcher) {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write && !contains(c, event.Name) {
					autoLog := util.ParseAutoLog(event.Name)
					if contains(ignore, autoLog.Name) {
						continue
					}
					c = append(c, event.Name)
					exec(autoLog)
				}
			case _err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Println("error:", _err)
			}
		}
	}(watcher)
	for s := range sc {
		err = watcher.Add(s)
		util.CheckError(err)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
