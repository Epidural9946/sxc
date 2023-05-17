package util

import (
	"encoding/hex"
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func init() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	//logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
}

func CheckError(err error) {
	if err != nil {
		logger.Errorln("error: ", err)
	}
}
func CheckErrorF(err error) {
	if err != nil {
		logger.Errorln("error: ", err)
	}
}

func CheckErrorExec(err error, f func()) {
	if err != nil {
		f()
	}
}

func HexToStr(hex1 string) string {
	d, _ := hex.DecodeString(hex1)
	d, _ = gbkToUtf8(d)
	return string(d)
}
