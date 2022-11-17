package main

type rating struct {
	ty string     `json:"type"`
	value uint    `json:"value"`
}

type qa struct {
	question string    `json:"question"`
	answer string      `json:"answer"`
}

type feedback struct {
	title string    `json:"title"`
	rate  rating    `json:"rate"`
	qas   []qa      `json:"questions"`
}

type FullFeedback struct {
	channel string    `json:"channel"`
	fb      feedback  `json:"feedback"`
}
