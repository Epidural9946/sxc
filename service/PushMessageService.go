package service

import (
	"bytes"
	"context"
	"github.com/dstotijn/go-notion"
	"io"
	"net/http"
	"strings"
	"time"
	"zpaul.org/chd/sxc/util"
)

type httpTransport struct {
	w io.Writer
}

// RoundTrip implements http.RoundTripper. It multiplexes the read HTTP response
// data to an io.Writer for debugging.
func (t *httpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	res.Body = io.NopCloser(io.TeeReader(res.Body, t.w))

	return res, nil
}
func PushPlusExec(apiKey string, parentPageID string, message util.XCAutoLog) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn(err) // 将 interface{} 转型为具体类型。
		}
	}()
	ctx := context.Background()
	buf := &bytes.Buffer{}
	httpClient := &http.Client{
		Timeout:   20 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))
	params := notion.CreatePageParams{
		ParentType: notion.ParentTypePage,
		ParentID:   parentPageID,
		Title:      []notion.RichText{{Text: &notion.Text{Content: message.Name}}},
		DatabasePageProperties: &notion.DatabasePageProperties{
			"Name": notion.DatabasePageProperty{
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: message.Name,
						},
					},
				},
			},
			"状态": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{Content: message.Msg},
					},
				},
			},
			"角色名称": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{Content: message.Account},
					},
				},
			},
			"原等级":  notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.BeginLevel))},
			"原经验条": notion.DatabasePageProperty{Number: notion.Float64Ptr(message.BeginExp)},
			"后等级":  notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.EndLevel))},
			"后经验条": notion.DatabasePageProperty{Number: notion.Float64Ptr(message.EndExp)},
			"回城复活": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.Revive2))},
			"苏生复活": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.Revive1))},
			"图鉴激活": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{Content: strings.Join(message.ToBooks(), " | ")},
					},
				},
			},
			"时长": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.TimeCons))},
			"翻牌": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{Content: strings.Join(message.ToCards(), " | ")},
					},
				},
			},
		},
		Children: []notion.Block{
			notion.Heading2Block{RichText: []notion.RichText{{Text: &notion.Text{Content: "产出"}}}},
			notion.TableBlock{
				TableWidth:      2,
				HasColumnHeader: true,
				Children:        util.HandleDataToTable(message.Acquisition),
			},
			notion.Heading2Block{RichText: []notion.RichText{{Text: &notion.Text{Content: "消耗"}}}},
			notion.TableBlock{
				TableWidth:      2,
				HasColumnHeader: true,
				Children:        util.HandleDataToTable(message.Consumables),
			},
			notion.Heading2Block{RichText: []notion.RichText{{Text: &notion.Text{Content: "日志"}}}},
		},
	}
	_l := util.HandleDataToNumList(message.Log)
	for _, text := range _l {
		params.Children = append(params.Children, text)
	}
	_, err := client.CreatePage(ctx, params)
	util.CheckErrorExec(err, func(err error) {
		logger.Warnf("Message: %s", message.Name)
		logger.Warnf("Error  : %s", err)
	})
	if err == nil {
		logger.Infof("成功发送：%s", message.Name)
	}

}
