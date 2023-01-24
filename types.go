package main

type Rate struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type QAndA struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type FullFeedback struct {
	Title     string  `json:"title"`
	Rate      Rate    `json:"rate"`
	Questions []QAndA `json:"questions"`
}

type Feedback struct {
	Channel string       `json:"channel"`
	Full    FullFeedback `json:"feedback"`
}
