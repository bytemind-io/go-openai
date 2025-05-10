package openai

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"
)

type SpeechModel string

const (
	TTSModel1      SpeechModel = "tts-1"
	TTSModel1HD    SpeechModel = "tts-1-hd"
	TTSModelCanary SpeechModel = "canary-tts"
)

type SpeechVoice string

const (
	VoiceAlloy   SpeechVoice = "alloy"
	VoiceEcho    SpeechVoice = "echo"
	VoiceFable   SpeechVoice = "fable"
	VoiceOnyx    SpeechVoice = "onyx"
	VoiceNova    SpeechVoice = "nova"
	VoiceShimmer SpeechVoice = "shimmer"
)

type SpeechResponseFormat string

const (
	SpeechResponseFormatMp3  SpeechResponseFormat = "mp3"
	SpeechResponseFormatOpus SpeechResponseFormat = "opus"
	SpeechResponseFormatAac  SpeechResponseFormat = "aac"
	SpeechResponseFormatFlac SpeechResponseFormat = "flac"
	SpeechResponseFormatWav  SpeechResponseFormat = "wav"
	SpeechResponseFormatPcm  SpeechResponseFormat = "pcm"
)

type FloatFrac float64

func (f FloatFrac) MarshalJSON() ([]byte, error) {
	n := float64(f)
	if math.IsInf(n, 0) || math.IsNaN(n) {
		return nil, errors.New("unsupported number")
	}
	prec := -1
	if math.Trunc(n) == n {
		prec = 1 // Force ".0" for integers.
	}
	return strconv.AppendFloat(nil, n, 'f', prec, 64), nil
}

type CreateSpeechRequest struct {
	Model             SpeechModel          `json:"model"`
	Input             string               `json:"input"`
	Voice             SpeechVoice          `json:"voice"`
	ResponseFormat    SpeechResponseFormat `json:"response_format,omitempty"`     // Optional, default to mp3
	Stream            bool                 `json:"stream,omitempty"`              // Optional, default to false
	Speed             FloatFrac            `json:"speed,omitempty"`               // Optional, default to 1.0 [0.5-2.0]
	Language          string               `json:"language,omitempty"`            // 音频语言：zh
	Volume            FloatFrac            `json:"volume,omitempty"`              // 音频：音量【0 -10】，默认1
	Pitch             int                  `json:"pitch,omitempty"`               // 音频：语调【-12， 12】，默认0
	Bitrate           int                  `json:"bitrate,omitempty"`             // 音频码率： Optional, default to 128000
	SampleRate        int                  `json:"sample_rate,omitempty"`         // 音频采样率： Optional, default to 24000
	Channel           int                  `json:"channel,omitempty"`             // 音频声道数： Optional, default to 1
	ReferenceVoiceWav string               `json:"reference_voice_wav,omitempty"` // 参考音频路径
	TimberWeights     map[string]FloatFrac `json:"timber_weights,omitempty"`      // 融合音色权重列表
}

func (c *Client) CreateSpeech(ctx context.Context, request CreateSpeechRequest) (response RawResponse, err error) {
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL("/audio/speech", withModel(string(request.Model))),
		withBody(request),
		withContentType("application/json"),
	)
	if err != nil {
		return
	}

	return c.sendRequestRaw(req)
}
