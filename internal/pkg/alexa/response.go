package alexa

const MAX_IMAGE_SIZE = int64(500 * 1024) // 500KB

// NewResponse builds a simple response. The session can
// be optionally ended by setting 'endSession' to true.
func NewResponse(title, text string, endSession bool) Response {
	r := Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &Payload{
				Type:  "PlainText",
				Text:  text,
				Title: title,
			},
			Card: &Payload{
				Type:    "Standard",
				Title:   title,
				Content: text,
				Text:    text,
			},
			ShouldEndSession: endSession,
		},
	}
	return r
}

func NewImageResponse(title, text string, imageSmallUrl string, imageLargeUrl string, endSession bool) Response {
	r := Response{
		Version: "1.0",
		Body: ResBody{
			OutputSpeech: &Payload{
				Type:  "PlainText",
				Text:  text,
				Title: title,
			},
			Card: &Payload{
				Type:    "Standard",
				Title:   title,
				Content: text,
				Text:    text,
				Image: Image{
					SmallImageURL: imageSmallUrl,
					LargeImageURL: imageLargeUrl,
				},
			},
			ShouldEndSession: endSession,
		},
	}
	return r
}

// Response is the response back to the Alexa speech service.
type Response struct {
	Version           string                 `json:"version"`
	SessionAttributes map[string]interface{} `json:"sessionAttributes,omitempty"`
	Body              ResBody                `json:"response"`
}

// ResBody is the actual body of the response.
type ResBody struct {
	OutputSpeech     *Payload  `json:"outputSpeech,omitempty"`
	Card             *Payload  `json:"card,omitempty"`
	Reprompt         *Reprompt `json:"reprompt,omitempty"`
	ShouldEndSession bool      `json:"shouldEndSession"`
}

type Reprompt struct {
	OutputSpeech Payload `json:"outputSpeech,omitempty"`
}

type Image struct {
	SmallImageURL string `json:"smallImageUrl,omitempty"`
	LargeImageURL string `json:"largeImageUrl,omitempty"`
}

type Payload struct {
	Type    string `json:"type,omitempty"`
	Title   string `json:"title,omitempty"`
	Text    string `json:"text,omitempty"`
	SSML    string `json:"ssml,omitempty"`
	Content string `json:"content,omitempty"`
	Image   Image  `json:"image,omitempty"`
}
