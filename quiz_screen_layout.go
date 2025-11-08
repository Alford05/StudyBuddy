package main

import (
	"fmt"
	"log"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (a *App) layoutQuizScreen(gtx layout.Context) layout.Dimensions {
	if a.CurrentQuestion.Sentence == "" {
		return material.Body1(a.Theme, "No question loaded").Layout(gtx)
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(material.H5(a.Theme, fmt.Sprintf("Question %d", a.CurrentIndex+1)).Layout),
			layout.Rigid(material.Body1(a.Theme, a.CurrentQuestion.Sentence).Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return a.layoutOption(gtx, 0)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return a.layoutOption(gtx, 1)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return a.layoutOption(gtx, 2)
						})
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							return a.layoutOption(gtx, 3)
						})
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if a.ShowFeedback && a.FeedbackMsg != "" {
					return layout.Inset{Top: unit.Dp(16)}.Layout(gtx,
						material.Body2(a.Theme, a.FeedbackMsg).Layout,
					)
				}
				return layout.Dimensions{}
			}),
		)
	})
}

func (a *App) layoutOption(gtx layout.Context, i int) layout.Dimensions {
	if i >= len(a.CurrentQuestion.Options) {
		return layout.Dimensions{}
	}

	btn := material.Button(a.Theme, &a.OptionButtons[i], a.CurrentQuestion.Options[i])
	if a.OptionButtons[i].Clicked(gtx) && !a.ShowFeedback {
		log.Printf("Selected: %s\nCorrect Answer:%s", a.CurrentQuestion.Options[i], a.CurrentQuestion.Options[a.CurrentQuestion.CorrectIndex])
		if i == a.CurrentQuestion.CorrectIndex {
			a.Score++
			a.FeedbackMsg = "✅ Correct!"
			fmt.Println("✅ Correct!")
		} else {
			a.FeedbackMsg = fmt.Sprintf("❌ Wrong. The correct answer is: %s", a.CurrentQuestion.Options[a.CurrentQuestion.CorrectIndex])
			fmt.Println("❌ Wrong.")
		}
		a.ShowFeedback = true
		a.Window.Invalidate()

		go func() {
			time.Sleep(2 * time.Second)
			a.CurrentIndex++
			a.FeedbackMsg = ""
			a.ShowFeedback = false

			if a.CurrentIndex >= 15 {
				a.State = QuizComplete
				a.Window.Invalidate()
				return
			}
			a.State = ScreenLoading
			a.Window.Invalidate() // update screen
			a.loadNextQuestion()
		}()
	}
	return btn.Layout(gtx)
}
