package service

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/dstotijn/go-notion"
	"io"
	"log"
	"net/http"
	"os"
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
		Timeout:   10 * time.Second,
		Transport: &httpTransport{w: buf},
	}
	client := notion.NewClient(apiKey, notion.WithHTTPClient(httpClient))

	params := notion.CreatePageParams{
		ParentType: notion.ParentTypePage,
		ParentID:   parentPageID,
		Title: []notion.RichText{
			{
				Text: &notion.Text{
					Content: message.Name,
				},
			},
		},
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
			"原经验条": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.BeginExp))},
			"后等级":  notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.EndLevel))},
			"后经验条": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.EndExp))},
			"回城复活": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.Revive1))},
			"苏生复活": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.Revive2))},
			"图鉴激活": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{Content: strings.Join(message.Card, " ")},
					},
				},
			},
			"时长": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(message.TimeCons))},
			"翻牌": notion.DatabasePageProperty{
				RichText: []notion.RichText{
					{
						Text: &notion.Text{Content: strings.Join(message.Book, " ")},
					},
				},
			},
			"MaxLevel": notion.DatabasePageProperty{Number: notion.Float64Ptr(float64(1))},
		},
	}
	_, err := client.CreatePage(ctx, params)
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	decoded := map[string]interface{}{}
	if err := json.NewDecoder(buf).Decode(&decoded); err != nil {
		log.Fatal(err)
	}

	// Pretty print JSON reponse.
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(decoded); err != nil {
		log.Fatal(err)
	}
}
