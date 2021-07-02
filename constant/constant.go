package constant

import "github.com/JimmyZhangJW/biliStreamClient"

const (
	DanmuMsg = "DANMU_MSG"
	SendGift = "SEND_GIFT"
	ComboSend = "COMBO_SEND"
)

// 参考文档 https://cloud.tencent.com/document/product/1073/37995
var (
	IntellectualFemaleVoice = biliStreamClient.VoiceConfig{
		Endpoint:  "tts.tencentcloudapi.com",
		Region:    "ap-shanghai",
		VoiceCode: 101009,
	}

	CantoneseFemaleVoice = biliStreamClient.VoiceConfig{
		Endpoint:  "tts.tencentcloudapi.com",
		Region:    "ap-shanghai",
		VoiceCode: 101019,
	}

	DefaultBoyVoice = biliStreamClient.VoiceConfig{
		Endpoint:  "tts.tencentcloudapi.com",
		Region:    "ap-shanghai",
		VoiceCode: 101015,
	}

	DefaultGirlVoice = biliStreamClient.VoiceConfig{
		Endpoint:  "tts.tencentcloudapi.com",
		Region:    "ap-shanghai",
		VoiceCode: 101016,
	}

	VoiceTypeIntellectualFemale = 1

	VoiceTypeCantoneseFemale = 2

	VoiceTypeDefaultBoy = 3

	VoiceTypeDefaultGirl = 4
)