package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/JimmyZhangJW/biliStreamClient"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/spf13/viper"

	"github.com/cutety/danmu-reader/constant"
)

func init() {
	viper.AddConfigPath("config")
	viper.SetConfigName("config-dev")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("no such config file")
		} else {
			log.Println("read config file error")
		}
		log.Fatalln(err)
	}
}

var (
	roomId    = flag.Int("r", 0, "直播间ID")
	voiceType = flag.Int("v", 0, "声音类型")
	voice     = biliStreamClient.VoiceConfig{}
)

func main() {
	flag.Parse()
	if *roomId == 0 {
		var inputRoomId int
		fmt.Println("请输入直播间号：")
		fmt.Scanln(&inputRoomId)
		roomId = &inputRoomId
	}
	biliClient := biliStreamClient.New()
	biliClient.Connect(*roomId)

	if *voiceType == 0 {
		var inputVoiceType int
		fmt.Println("请输入你想听的声音类型：")
		fmt.Println("1：知性女声")
		fmt.Println("2：粤语女声")
		fmt.Println("3：女童")
		fmt.Println("4：男童")
		fmt.Scanln(&inputVoiceType)
		voice = chooseVoiceType(inputVoiceType)
	}
	voice = chooseVoiceType(*voiceType)

	for {
		packBody := <-biliClient.Ch
		switch packBody.Cmd {
		case constant.DanmuMsg:
			danmu, err := packBody.ParseDanmu()
			if err != nil {
				log.Fatalln(err)
			}
			processDanmuMessage(&danmu)
		case constant.SendGift:
			gift, err := packBody.ParseGift()
			if err != nil {
				log.Fatalln(err)
			}
			processGiftMessage(&gift)
		}

	}

}

func chooseVoiceType(voiceType int) biliStreamClient.VoiceConfig {
	var voice biliStreamClient.VoiceConfig
	switch voiceType {
	case constant.VoiceTypeIntellectualFemale:
		voice = constant.IntellectualFemaleVoice
	case constant.VoiceTypeCantoneseFemale:
		voice = constant.CantoneseFemaleVoice
	case constant.VoiceTypeDefaultBoy:
		voice = constant.DefaultGirlVoice
	case constant.VoiceTypeDefaultGirl:
		voice = constant.DefaultBoyVoice
	default:
		voice = constant.DefaultGirlVoice
	}

	return voice
}

func processGiftMessage(gift *biliStreamClient.Gift) {
	sender := gift.Sender.Name
	msg := gift.GiftName

	content := fmt.Sprintf("%s送了%s", sender, msg)
	log.Printf("%s", content)
	voice, err := parseTextToVoice(content)
	if err != nil {
		log.Fatal(err)
	}

	err = playVoice(voice)
	if err != nil {
		log.Fatal(err)
	}
}

func parseTextToVoice(content string) ([]byte, error) {
	SID := viper.GetString("tencent.secretID")
	SK := viper.GetString("tencent.secretKey")
	encodedVoice, err := biliStreamClient.GetVoiceFromTencentCloud(SID, SK, voice, content)
	if err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(encodedVoice)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func playVoice(data []byte) error {
	streamer, format, err := wav.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(streamer)

	time.Sleep(time.Second * 1)
	return nil
}

func processDanmuMessage(danmu *biliStreamClient.DanmuMsg) error {
	sender := danmu.Name
	msg := danmu.Message

	content := fmt.Sprintf("%s说  %s", sender, msg)
	log.Printf("%s", content)
	content = biliStreamClient.Sanitize(content)
	encodedVoice, err := biliStreamClient.GetVoiceFromTencentCloud(viper.GetString("tencent.secretID"), viper.GetString("tencent.secretKey"), voice, content)
	if err != nil {
		log.Fatalln(err)
	}

	data, err := base64.StdEncoding.DecodeString(encodedVoice)
	if err != nil {
		log.Fatalln(err)
	}

	streamer, format, err := wav.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatalln(err)
	}

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/15))
	speaker.Play(streamer)

	time.Sleep(time.Second * 2)

	return nil
}
