package kuaishou

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

// KuaiShou 结构体表示快手直播对象
type KuaiShou struct {
	RID string
}

// GetRealURL 方法用于获取直播流媒体地址
func (k *KuaiShou) GetRealURL() (string, error) {
	userAgentList := []string{
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/45.0.2454.101 Safari/537.36",
		"Mozilla/5.0 (Macintosh; U; PPC Mac OS X 10.5; en-US; rv:1.9.2.15) Gecko/20110303 Firefox/3.6.15",
	}

	headers := map[string]string{
		"User-Agent": randomChoice(userAgentList),
		"cookie":     "did=web_d563dca728d28b00336877723e0359ed",
		"Referer":    "https://live.kuaishou.com/u/" + k.RID,
	}

	// 发送 GET 请求
	res, err := makeRequest("https://live.kuaishou.com/u/"+k.RID, headers)
	if err != nil {
		return "", err
	}

	// 正则匹配直播流数据
	liveStreamMatch := regexp.MustCompile(`liveStream":(.*),"author`).FindStringSubmatch(res)
	if len(liveStreamMatch) < 2 {
		return "", fmt.Errorf("直播间不存在或未开播")
	}

	// 解析 JSON 数据
	var liveStreamData map[string]interface{}
	if err := json.Unmarshal([]byte(liveStreamMatch[1]), &liveStreamData); err != nil {
		return "", err
	}

	// 获取最高画质的流媒体地址
	playUrls := liveStreamData["playUrls"].([]interface{})
	hlsPlayUrls := playUrls[0].(map[string]interface{})
	url := hlsPlayUrls["adaptationSet"].(map[string]interface{})["representation"].([]interface{})[len(playUrls)-1].(map[string]interface{})["url"].(string)

	return url, nil
}

// makeRequest 方法用于发送 HTTP 请求
func makeRequest(url string, headers map[string]string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// randomChoice 从字符串切片中随机选择一个元素
func randomChoice(choices []string) string {
	rand.Seed(time.Now().UnixNano())
	return choices[rand.Intn(len(choices))]
}

func KuaiShouStreamUrl(c *cli.Context) error {
	roomId := c.Args().Get(0) // 房间id
	ks := KuaiShou{RID: roomId}
	realURL, err := ks.GetRealURL()
	if err != nil {
		return err
	}
	fmt.Println("快手直播间视频流地址：", realURL)
	return nil
}
