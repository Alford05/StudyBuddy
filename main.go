package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type AppState int

const (
	ScreenInput AppState = iota
	ScreenLoading
	ScreenQuiz
	QuizComplete
)

type Question struct {
	Sentence     string
	Options      []string
	CorrectIndex int
	AnswerWords  []string
}

type App struct {
	State           AppState
	Window          *app.Window
	Theme           *material.Theme
	Ops             op.Ops
	WordInput       widget.Editor
	StartButton     widget.Clickable
	WordBank        []string
	CurrentIndex    int
	CurrentQuestion Question
	OptionButtons   [4]widget.Clickable
	APIKey          string
	RNG             *rand.Rand
	ErrorMsg        string

	Score         int
	RestartButton widget.Clickable

	FeedbackMsg  string
	ShowFeedback bool

	QuestionBuffer chan Question
	PreloadDone    chan struct{}
}

func main() {
	go func() {
		var w app.Window
		w.Option(app.Title("Vocabulary Test Time!"))
		w.Option(app.Size(unit.Dp(900), unit.Dp(600)))
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))

		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading hidden file")
		}
		apiKey := os.Getenv("SecretKey")
		if apiKey == "" {
			log.Println("SecretKey variable not set.")
		}

		var ops op.Ops
		w.Invalidate()

		theme := material.NewTheme()
		theme.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))

		myApp := &App{
			State:          ScreenInput,
			Theme:          theme,
			APIKey:         apiKey,
			RNG:            rng,
			WordBank:       []string{},
			Window:         &w,
			QuestionBuffer: make(chan Question, 1),
			PreloadDone:    make(chan struct{}),
		}
		myApp.WordInput.SingleLine = false

		for {
			e := w.Event()
			switch e := e.(type) {
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				layout.Flex{}.Layout(gtx)
				if myApp.Theme == nil {
					myApp.Theme = material.NewTheme()
				}

				myApp.Layout(gtx)
				e.Frame(gtx.Ops)
			case app.DestroyEvent:
				log.Println("Window closed")
				os.Exit(0)
			}
		}
	}()
	app.Main()
}
