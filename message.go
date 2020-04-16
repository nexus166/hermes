package hermes

import (
	"bytes"
	"encoding/json"
	"time"
)

type AttachmentFields struct {
	Short bool   `json:"short,omitempty"`
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
}

type TemplateAttachments struct {
	AudioURL          string             `json:"audio_url,omitempty"`
	AuthorIcon        string             `json:"author_icon,omitempty"`
	AuthorLink        string             `json:"author_link,omitempty"`
	AuthorName        string             `json:"author_name,omitempty"`
	Collapsed         bool               `json:"collapsed,omitempty"`
	Color             string             `json:"color,omitempty"`
	Fields            []AttachmentFields `json:"fields,omitempty"`
	ImageURL          string             `json:"image_url,omitempty"`
	MessageLink       string             `json:"message_link,omitempty"`
	Text              string             `json:"text,omitempty"`
	ThumbURL          string             `json:"thumb_url,omitempty"`
	Title             string             `json:"title,omitempty"`
	TitleLink         string             `json:"title_link,omitempty"`
	TitleLinkDownload bool               `json:"title_link_download,omitempty"`
	Timestamp         time.Time          `json:"ts,omitempty"`
	VideoURL          string             `json:"video_url,omitempty"`
}

type POSTTemplate struct {
	Alias       string                `json:"alias,omitempty"`
	Avatar      string                `json:"avatar,omitempty"`
	Channel     string                `json:"channel,omitempty"`
	RoomID      string                `json:"roomId,omitempty"`
	Text        string                `json:"text"`
	Emoji       string                `json:"emoji,omitempty"`
	Attachments []TemplateAttachments `json:"attachments,omitempty"`
}

func MakePayload(jsonStruct POSTTemplate) (*bytes.Buffer, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(&jsonStruct); err != nil {
		return nil, err
	}
	return b, nil
}
