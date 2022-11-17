package main

type rating struct {
	value uint
}

type qa struct {
	question string
	answer string
}

type feedback struct {
	title string
	rate  rating
	qas   []qa
}

type FullFeedback struct {
	channel string
	fb      feedback
	
}
