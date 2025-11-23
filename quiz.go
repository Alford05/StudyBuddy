package main

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (a *App) layoutQuizCompleteScreen(gtx layout.Context) layout.Dimensions {
	percent := float64(a.Score) / 15.0 * 100.0
	msg := fmt.Sprintf("Quiz complete!\nYour score: %d / 15\n(%.1f%%)", a.Score, percent)

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(gtx,
			layout.Rigid(material.Body1(a.Theme, msg).Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.Theme, &a.RestartButton, "Restart Quiz")
				if a.RestartButton.Clicked(gtx) {
					select {
					case <-a.StopPreload:
					default:
						close(a.StopPreload)
					}
					a.CurrentIndex = 0
					a.Score = 0
					a.ShowFeedback = false
					a.FeedbackMsg = ""
					a.ErrorMsg = ""

					a.QuestionBuffer = make(chan Question, 1)
					a.PreloadDone = make(chan struct{})
					a.StopPreload = make(chan struct{})

					question, err := buildSingleWordQuestion(a.WordBank, a.CurrentIndex, a.APIKey, a.RNG)
					if err != nil {
						a.ErrorMsg = "Error generating first question."
						return btn.Layout(gtx)
					}
					a.CurrentQuestion = question

					go a.startPreloading()
					a.State = ScreenLoading
					a.Window.Invalidate()

					go a.loadNextQuestion()
				}
				return btn.Layout(gtx)
			}),
		)
	})
}
