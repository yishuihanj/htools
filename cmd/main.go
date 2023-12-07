package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
	"ytools/douyin"
	"ytools/kuaishou"
)

func main() {
	app := cli.NewApp()
	app.Name = "htools"
	app.Commands = commands()
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("run failed!", err)
	}
}

func commands() []*cli.Command {
	return []*cli.Command{
		{
			Name:      "douyinstream",
			Usage:     "-douyinstream rommId flv/hls",
			ArgsUsage: "根据直播间ID获取视频流",
			Action:    douyin.DouYinStream,
		},
		{
			Name:      "kuaishoustream",
			Usage:     "-kuaishoustream rommId",
			ArgsUsage: "根据直播间ID获取视频流",
			Action:    kuaishou.KuaiShouStreamUrl,
		},
	}
}
