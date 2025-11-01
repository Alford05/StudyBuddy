package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (a *App) Layout(gtx layout.Context) layout.Dimensions {
	switch a.State {
	case ScreenInput:
		return a.layoutInputScreen(gtx)
	case ScreenLoading:
		return a.layoutLoadingScreen(gtx)
	case ScreenQuiz:
		return a.layoutQuizScreen(gtx)
	case QuizComplete:
		return a.layoutQuizCompleteScreen(gtx)
	default:
		return layout.Dimensions{}
	}
}

func (a *App) layoutInputScreen(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(material.H6(a.Theme, "Enter 10 vocabulary words:").Layout),
			layout.Rigid(material.Editor(a.Theme, &a.WordInput, "e.g. innovate, challenge, create").Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(a.Theme, &a.StartButton, "Start Quiz")
				for a.StartButton.Clicked(gtx) {
					words := parseWords(a.WordInput.Text())
					if len(words) != 10 {
						a.ErrorMsg = "Please enter 10 words."
					} else {
						a.ErrorMsg = ""
						a.WordBank = words
						a.CurrentIndex = 0

						question, err := buildSingleWordQuestion(a.WordBank, a.CurrentIndex, a.APIKey, a.RNG)
						if err != nil {
							a.ErrorMsg = "Error generating first question."
							break
						}
						a.CurrentQuestion = question
						go a.startPreloading()
						a.State = ScreenQuiz
						a.Window.Invalidate()
					}
				}
				return btn.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if a.ErrorMsg != "" {
					return material.Body2(a.Theme, a.ErrorMsg).Layout(gtx)
				}
				return layout.Dimensions{}
			}),
		)
	})
}

func (a *App) layoutLoadingScreen(gtx layout.Context) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.H6(a.Theme, "Loading Sentence.....").Layout(gtx)
	})
}
