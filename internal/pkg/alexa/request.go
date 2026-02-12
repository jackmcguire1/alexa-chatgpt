package alexa

const (
	// built-in request types
	IntentRequestType       = "IntentRequest"
	LaunchRequestType       = "LaunchRequest"
	SessionEndedRequestType = "SessionEndedRequest"
)

// Request represents the structure of the request sent from Alexa.
type Request struct {
	Version string  `json:"version"`
	Session Session `json:"session"`
	Body    ReqBody `json:"request"`
	Context Context `json:"context"`
}

// Session represents the Alexa skill session.
type Session struct {
	New         bool   `json:"new"`
	SessionID   string `json:"sessionId"`
	Application struct {
		ApplicationID string `json:"applicationId"`
	} `json:"application"`
	Attributes map[string]any `json:"attributes"`
	User       struct {
		UserID      string `json:"userId"`
		AccessToken string `json:"accessToken,omitempty"`
	} `json:"user"`
}

// Context represents the Alexa skill request context.
type Context struct {
	System struct {
		APIAccessToken string `json:"apiAccessToken"`
		Device         struct {
			DeviceID string `json:"deviceId,omitempty"`
		} `json:"device"`
		Application struct {
			ApplicationID string `json:"applicationId,omitempty"`
		} `json:"application"`
	} `json:"System"`
}

// ReqBody is the request body from Alexa.
type ReqBody struct {
	Type        string `json:"type"`
	RequestID   string `json:"requestId"`
	Timestamp   string `json:"timestamp"`
	Locale      string `json:"locale"`
	Intent      Intent `json:"intent"`
	Reason      string `json:"reason,omitempty"`
	DialogState string `json:"dialogState,omitempty"`
}

// Intent is the Alexa skill intent.
type Intent struct {
	Name               string          `json:"name"`
	Slots              map[string]Slot `json:"slots"`
	ConfirmationStatus string          `json:"confirmationStatus"`
}

// Slot is an Alexa skill slot.
type Slot struct {
	Name        string      `json:"name"`
	Value       string      `json:"value"`
	Resolutions Resolutions `json:"resolutions"`
}

type Resolutions struct {
	ResolutionPerAuthority []struct {
		Values []struct {
			Value struct {
				Name string `json:"name"`
				Id   string `json:"id"`
			} `json:"value"`
		} `json:"values"`
	} `json:"resolutionsPerAuthority"`
}
