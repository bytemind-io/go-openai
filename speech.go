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
	Model              SpeechModel          `json:"model"`
	Input              string               `json:"input"`
	Voice              SpeechVoice          `json:"voice"`
	ResponseFormat     SpeechResponseFormat `json:"response_format,omitempty"`       // Optional, default to mp3
	Speed              float64              `json:"speed,omitempty"`                 // Optional, default to 1.0
	Language           string               `json:"language,omitempty"`              // todo: 新增
	ReferWavPathGpt    string               `json:"refer_wav_path_gpt,omitempty"`    // 参考音频
	ReferWavPathSovits map[string]FloatFrac `json:"refer_wav_path_sovits,omitempty"` // 融合音频
	Prompt             string               `json:"prompt,omitempty"`                // 音频样本文本
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
