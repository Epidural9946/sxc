package util

import (
	"bytes"
	"encoding/json"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type xcLog struct {
	Name       string           `json:"name"`
	BeginTime  string           `json:"begintime"`
	EndTime    string           `json:"endtime"`
	EndMsg     string           `json:"endmsg"`
	BeginLevel int              `json:"beginlevel"`
	BeginExp   int              `json:"beginexp"`
	EndLevel   int              `json:"endlevel"`
	EndExp     int              `json:"endexp"`
	Revive     []map[string]int `json:"revive"`
	Item1      []map[string]int `json:"item1"`
	Item2      []map[string]int `json:"item2"`
	Msg        []string         `json:"msg"`
	Card       []string         `json:"card"`
	Collect    []string         `json:"collect"`
}

type XCAutoLog struct {
	Account     string //账号
	Name        string //名称
	TimeCons    int64  //耗时  ParseLog [10000000 = 1s]
	BeginLevel  int
	EndLevel    int
	BeginExp    float64
	EndExp      float64
	Revive1     int            //苏生
	Revive2     int            //城市
	Msg         string         //消息
	Acquisition map[string]int //获得品
	Consumables map[string]int //消耗品
	Card        []string       //翻牌
	Book        []string       //图鉴
}

func (x *XCAutoLog) ToBooks() []string {
	l := make([]string, 0)
	for _, item := range x.Book {
		split := strings.Split(item, "成功激活图鉴 : ")
		l = append(l, split[1])
	}
	return l
}
func (x *XCAutoLog) ToCards() []string {
	l := make([]string, 0)
	for _, item := range x.Card {
		split := strings.Split(item, "领奖奖励:")
		if len(split) == 2 {
			l = append(l, split[1])
		}
	}
	return l
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
	CheckError(err)
	if err != nil {
		return XCAutoLog{}
	}
	return log.getXCAutoLog()
}

func (log *xcLog) getRevive(i int) int {
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
	if i == 1 {
		return s
	} else {
		return c
	}
}
func (log *xcLog) getXCAutoLog() XCAutoLog {
	return XCAutoLog{
		Name:        log.Name,
		Msg:         log.EndMsg,
		BeginLevel:  log.BeginLevel,
		EndLevel:    log.EndLevel,
		BeginExp:    float64(log.BeginExp) / 100,
		EndExp:      float64(log.EndExp) / 100,
		TimeCons:    log.getTimeCons(),
		Revive1:     log.getRevive(1),
		Revive2:     log.getRevive(2),
		Acquisition: changeStruct(log.Item1),
		Consumables: changeStruct(log.Item2),
		Card:        log.Card,
		Book:        log.Collect,
	}
}

func (log *xcLog) getTimeCons() int64 {
	begin, _ := strconv.ParseInt(log.BeginTime, 10, 0)
	end, _ := strconv.ParseInt(log.EndTime, 10, 0)
	return (end - begin) / 10000000 / 60
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
