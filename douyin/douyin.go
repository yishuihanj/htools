package douyin

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func DouYinStream(c *cli.Context) error {
	roomId := c.Args().Get(0) // 房间id
	stream := c.Args().Get(1) // flv or hls

	if roomId == "" {
		return errors.New("参数错误")
	}

	if stream != "flv" && stream != "hls" {
		return errors.New("视频流必须是 flv 或者 hls")
	}
	douyinObj := DouYinFlv{Rid: roomId, Stream: stream}
	douyinURL := douyinObj.getDouyinURL()

	fmt.Println("抖音直播间视频流地址：", douyinURL)

	return nil
}

type DouYinFlv struct {
	Rid    string
	Stream string
}

func (d *DouYinFlv) getDouyinURL() string {
	liveURL := fmt.Sprintf("https://live.douyin.com/%s", d.Rid)

	// Send initial request to obtain __ac_nonce
	client := &http.Client{}
	req, err := http.NewRequest("GET", liveURL, nil)
	if err != nil {
		fmt.Println("Error creating initial request:", err)
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	oresp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending initial request:", err)
		return ""
	}
	defer oresp.Body.Close()

	// Extract __ac_nonce from Set-Cookie header
	cookieHeader := oresp.Header.Get("Set-Cookie")
	acNonceMatch := regexp.MustCompile(`(?i)__ac_nonce=(.*?);`).FindStringSubmatch(cookieHeader)
	if len(acNonceMatch) < 2 {
		return ""
	}
	acNonce := acNonceMatch[1]

	// Set __ac_nonce cookie and send another request
	req.Header.Set("Cookie", fmt.Sprintf("__ac_nonce=%s", acNonce))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending second request:", err)
		return ""
	}
	defer resp.Body.Close()

	// Extract ttwid from Set-Cookie header
	ttwidMatch := regexp.MustCompile(`(?i)ttwid=.*?;`).FindStringSubmatch(resp.Header.Get("Set-Cookie"))
	if len(ttwidMatch) < 1 {
		return ""
	}
	ttwid := ttwidMatch[0]

	// Build URL for final request
	url := fmt.Sprintf("https://live.douyin.com/webcast/room/web/enter/?aid=6383&app_name=douyin_web&live_id=1&device_platform=web&language=zh-CN&enter_from=web_live&cookie_enabled=true&screen_width=1728&screen_height=1117&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Chrome&browser_version=116.0.0.0&web_rid=%s", d.Rid)
	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating final request:", err)
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", ttwid)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Host", "live.douyin.com")
	req.Header.Set("Connection", "keep-alive")

	// Send the final request
	ress, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending final request:", err)
		return ""
	}
	defer ress.Body.Close()

	// Parse the JSON response
	body, err := ioutil.ReadAll(ress.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return ""
	}

	status := data["data"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["status"].(float64)

	if status != 2 {
		return ""
	}

	realURL := ""
	streamData := data["data"].(map[string]interface{})["data"].([]interface{})[0].(map[string]interface{})["stream_url"].(map[string]interface{})["live_core_sdk_data"].(map[string]interface{})["pull_data"].(map[string]interface{})["stream_data"].(string)

	var value map[string]interface{}
	if err := json.NewDecoder(strings.NewReader(streamData)).Decode(&value); err != nil {
		fmt.Println("Error decoding stream data:", err)
		return ""
	}

	if d.Stream == "flv" {
		realURL = value["data"].(map[string]interface{})["origin"].(map[string]interface{})["main"].(map[string]interface{})["flv"].(string)
	} else if d.Stream == "hls" {
		realURL = value["data"].(map[string]interface{})["origin"].(map[string]interface{})["main"].(map[string]interface{})["hls"].(string)
	}

	return realURL
}
