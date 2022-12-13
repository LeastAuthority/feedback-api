package main

type rate struct {
	typ string `json:"type"`
	val string `json:"value"`
}

type qAndA struct {
	question string `json:"question"`
	answer   string `json:"answer"`
}

type fullFeedback struct {
	title string          `json:"title"`
	rate  rate            `json:"rate"`
	questions []qAndA     `json:"questions"`
}

type feedback struct {
	channel  string               `json:"channel"`
	full     fullFeedback         `json:"feedback"`
}
