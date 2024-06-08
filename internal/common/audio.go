package common

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// MpCommonMpAudio 定义结构体来表示 XML 节点及其属性
type MpCommonMpAudio struct {
	XMLName           xml.Name `xml:"mp-common-mpaudio"`
	Src               string   `xml:"src,attr"`
	Isaac2            string   `xml:"isaac2,attr"`
	LowSize           string   `xml:"low_size,attr"`
	SourceSize        string   `xml:"source_size,attr"`
	HighSize          string   `xml:"high_size,attr"`
	Name              string   `xml:"name,attr"`
	PlayLength        string   `xml:"play_length,attr"`
	Author            string   `xml:"author,attr"`
	DataTopicID       string   `xml:"data-topic_id,attr"`
	DataTopicName     string   `xml:"data-topic_name,attr"`
	DataPluginName    string   `xml:"data-pluginname,attr"`
	DataTransState    string   `xml:"data-trans_state,attr"`
	DataVerifyState   string   `xml:"data-verify_state,attr"`
	VoiceEncodeFileID string   `xml:"voice_encode_fileid,attr"`
}

type MpCommonMpVoice struct {
	XMLName           xml.Name `xml:"mpvoice"`
	Src               string   `xml:"src,attr"`
	Isaac2            string   `xml:"isaac2,attr"`
	LowSize           string   `xml:"low_size,attr"`
	SourceSize        string   `xml:"source_size,attr"`
	HighSize          string   `xml:"high_size,attr"`
	Name              string   `xml:"name,attr"`
	PlayLength        string   `xml:"play_length,attr"`
	Author            string   `xml:"author,attr"`
	DataTopicID       string   `xml:"data-topic_id,attr"`
	DataTopicName     string   `xml:"data-topic_name,attr"`
	DataPluginName    string   `xml:"data-pluginname,attr"`
	DataTransState    string   `xml:"data-trans_state,attr"`
	DataVerifyState   string   `xml:"data-verify_state,attr"`
	VoiceEncodeFileID string   `xml:"voice_encode_fileid,attr"`
}

// AudioParse 解析音频信息
func AudioParse(xmlData string) (MpCommonMpAudio, error) {
	// 创建一个 MpCommonMpAudio 实例
	var audio MpCommonMpAudio

	// 解析 XML 数据
	err := xml.NewDecoder(strings.NewReader(xmlData)).Decode(&audio)
	if err != nil {
		fmt.Printf("Error decoding XML: %v\n", err)
		return MpCommonMpAudio{}, err
	}
	return audio, nil
}

func VoiceParse(xmlData string) (MpCommonMpVoice, error) {
	// 创建一个 MpCommonMpAudio 实例
	var audio MpCommonMpVoice

	// 解析 XML 数据
	err := xml.NewDecoder(strings.NewReader(xmlData)).Decode(&audio)
	if err != nil {
		fmt.Printf("Error decoding XML: %v\n", err)
		return MpCommonMpVoice{}, err
	}
	return audio, nil
}
