package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
)

type AppState int

const (
	ScreenInput AppState = iota
	ScreenQuiz
)

type Question struct {
	Sentence     string
	Options      []string
	CorrectIndex int
	AnswerWords  []string
}

type App struct {
	State           AppState
	WordInput       widget.Editor
	StartButton     widget.Clickable
	WordBank        []string
	CurrentIndex    int
	CurrentQuestion Question
	OptionButtons   [4]widget.Clickable
	APIKey          string
	RNG             *rand.Rand
}

func main() {
	go func() {
		w := new(app.Window)
		var ops op.Ops

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		apiKey := os.Getenv("SecretKey")

		myApp := &App{
			State:    ScreenInput,
			APIKey:   apiKey,
			RNG:      rng,
			WordBank: []string{},
		}
		myApp.WordInput.SingleLine = true

		for {
			e := w.NextEvent()
			switch e := e.(type) {
			case app.FrameEvent:
				gtx := layout.Context{
					Ops:    &ops,
					Metric: e.Metric,
					Now:    e.Now,
				}
				e.Frame(gtx.Ops)
			case app.DestroyEvent:
				log.Println("Window closed")
				os.Exit(0)
			}
		}

		//for e := range w.Run(func(e system.Event) {
		//	switch e := e.(type) {
		//	case app.FrameEvent:
		//		gtx := layout.Context{
		//			Ops:    &ops,
		//			Metric: e.Metric,
		//			Now:    e.Now,
		//		}
		//		e.Frame(gtx.Ops)
		//	case app.DestroyEvent:
		//		log.Println("Window closaed")
		//		return
		//	}
	}()
	app.Main()
}
