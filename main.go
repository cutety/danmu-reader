package main

import (
	"fmt"
	"log"

	"github.com/JimmyZhangJW/biliStreamClient"

	"github.com/cutety/danmu-reader/constant"
)

func main() {
	biliClient := biliStreamClient.New()
	biliClient.Connect(1017)

	for {
		packBody := <- biliClient.Ch
		switch packBody.Cmd {
		case constant.DanmuMsg:
			danmu, err := packBody.ParseDanmu()
			if err != nil {
				log.Fatalln(err)
			}
			processDanmuMessage(&danmu)
		case constant.SendGift:
			//danmu, err := packBody.ParseGift()
			//if err != nil {
			//	log.Fatalln(err)
			//}
			//
			//biliStreamClient.GetVoiceFromTencentCloud()
		}

	}

}

func processDanmuMessage(danmu *biliStreamClient.DanmuMsg) error {
	sender := danmu.Name
	msg := danmu.Message

	content := fmt.Sprintf("%s说  %s", sender, msg)
	log.Printf("MSG:%s", content)
	res, err := biliStreamClient.GetVoiceFromTencentCloud(constant.SecretID, constant.SecretKey, biliStreamClient.DefaultBoyVoice, content)
	if err != nil {
		log.Fatalln(err)
	}



	return nil
}