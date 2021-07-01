package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/JimmyZhangJW/biliStreamClient"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"

	"github.com/cutety/danmu-reader/constant"
)

func main() {
	biliClient := biliStreamClient.New()
	biliClient.Connect(613765)

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
	content = biliStreamClient.Sanitize(content)
	encodedVoice, err := biliStreamClient.GetVoiceFromTencentCloud(constant.SecretID, constant.SecretKey, biliStreamClient.DefaultBoyVoice, content)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := base64.StdEncoding.DecodeString(encodedVoice)

	streamer, format, err := wav.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalln(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Minute/50))
	speaker.Play(streamer)


	return nil
}