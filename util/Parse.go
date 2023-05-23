package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strconv"
	"time"
)

type xcLog struct {
	Name      string           `json:"name"`
	BeginTime string           `json:"begintime"`
	EndTime   string           `json:"endtime"`
	EndMsg    string           `json:"endmsg"`
	Revive    []map[string]int `json:"revive"`
	Item1     []map[string]int `json:"item1"`
	Item2     []map[string]int `json:"item2"`
	Msg       []string         `json:"msg"`
	Card      []string         `json:"card"`
	Collect   []string         `json:"collect"`
}

type XCAutoLog struct {
	Account     string         //账号
	Name        string         //名称
	TimeCons    int64          //耗时  ParseLog [10000000 = 1s]
	Revive      string         //死亡次数
	Msg         string         //消息
	Acquisition map[string]int //获得品
	Consumables map[string]int //消耗品
	Card        []string       //翻牌
	Book        []string       //图鉴
}

func ParseAutoLog(path string) XCAutoLog {
	time.Sleep(2 * time.Millisecond)
	data, err := os.ReadFile(path)
	CheckError(err)
	jsonByte, err := gbkToUtf8(data)
	CheckError(err)
	return parseNewVerContent(string(jsonByte))
}
func parseNewVerContent(content string) XCAutoLog {
	log := xcLog{}
	err := json.Unmarshal([]byte(content), &log)
	autoLog := XCAutoLog{}
	CheckError(err)
	if err != nil {
		return autoLog
	}
	begin, _ := strconv.ParseInt(log.BeginTime, 10, 0)
	end, _ := strconv.ParseInt(log.EndTime, 10, 0)
	autoLog.Name = log.Name
	autoLog.TimeCons = (end - begin) / 10000000 / 60
	s := 0
	c := 0
	for _, item := range log.Revive {
		for _, i := range item {
			if i == 2 {
				s++
			} else if i == 0 {
				c++
			}
		}
	}
	autoLog.Revive = fmt.Sprintf("%v/%v", c, s)
	autoLog.Msg = log.EndMsg
	autoLog.Acquisition = changeStruct(log.Item1)
	autoLog.Consumables = changeStruct(log.Item2)
	autoLog.Card = log.Card
	autoLog.Book = log.Collect
	return autoLog
}

func changeStruct(data []map[string]int) map[string]int {
	m := make(map[string]int)
	// 遍历切片中的每个 map
	for _, mp := range data {
		// 遍历 map 中的键值对，插入到新的 map 中
		for k, v := range mp {
			m[k] = v
		}
	}
	return m
}

// gbkToUtf8 GBK 转 UTF-8
func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}
