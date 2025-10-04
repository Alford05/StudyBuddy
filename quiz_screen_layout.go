package main

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (a *App) layoutQuizScreen(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if a.CurrentQuestion.Sentence == "" && len(a.WordBank) > 0 {
		q, err := buildQuestion(a.WordBank, a.APIKey, a.RNG)
		if err != nil {
			fmt.Println("Error generating question:", err)
			return layout.Dimensions{}
		}
		a.CurrentQuestion = q
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceEvenly,
		}.Layout(gtx,
			layout.Rigid(material.H5(th, fmt.Sprintf("Question %d", a.CurrentIndex+1)).Layout),
			layout.Rigid(material.Body1(th, a.CurrentQuestion.Sentence).Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.layoutOption(gtx, th, 0)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.layoutOption(gtx, th, 1)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.layoutOption(gtx, th, 2)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return a.layoutOption(gtx, th, 3)
					}),
				)
			}),
		)
	})
}

func (a *App) layoutOption(gtx layout.Context, th *material.Theme, i int) layout.Dimensions {
	if i >= len(a.CurrentQuestion.Options) {
		return material.Body2(th, "Preparing...").Layout(gtx)
	}

	btn := material.Button(th, &a.OptionButtons[i], a.CurrentQuestion.Options[i])
	if a.OptionButtons[i].Clicked(gtx) {
		if i == a.CurrentQuestion.CorrectIndex {
			fmt.Println("✅ Correct!")
		} else {
			fmt.Println("❌ Wrong.")
		}
		a.CurrentIndex++
		q, err := buildQuestion(a.WordBank, a.APIKey, a.RNG)
		if err != nil {
			fmt.Println("Error generating next question:", err)
			return layout.Dimensions{}
		}
		a.CurrentQuestion = q
	}
	return btn.Layout(gtx)
}

//TODO
// Animate the feedback (Correct! vs Wrong) in the UI
// Slight delay before moving toteh next question
// Keep score and add a quiz complete screen
