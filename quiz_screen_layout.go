package main

import (
	"fmt"
	"log"
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func (a *App) layoutQuizScreen(gtx layout.Context) layout.Dimensions {
	if a.CurrentQuestion.Sentence == "" {
		return material.Body1(a.Theme, "No question loaded").Layout(gtx)
	}
	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			title := material.H5(a.Theme, fmt.Sprintf("Question %d", a.CurrentIndex+1))
			return layout.Inset{Bottom: unit.Dp(12)}.Layout(gtx, title.Layout)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			lbl := material.Label(a.Theme, unit.Sp(20), a.CurrentQuestion.Sentence)
			lbl.Alignment = text.Start
			lbl.MaxLines = 0
			lbl.Font.Weight = font.SemiBold
			return layout.Inset{
				Bottom: unit.Dp(16),
				Left:   unit.Dp(80),
				Right:  unit.Dp(20),
			}.Layout(gtx, lbl.Layout)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				a.layoutOptionInset(0),
				a.layoutOptionInset(1),
				a.layoutOptionInset(2),
				a.layoutOptionInset(3),
			)
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if !a.ShowFeedback || a.FeedbackMsg == "" {
				return layout.Dimensions{}
			}
			feedback := material.Body2(a.Theme, a.FeedbackMsg)
			return layout.Inset{
				Top: unit.Dp(20),
			}.Layout(gtx, feedback.Layout)
		}),
	)
}

func (a *App) layoutOptionInset(index int) layout.FlexChild {
	return layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{
			Top: unit.Dp(8),
		}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			maxwidth := gtx.Dp(400)
			if gtx.Constraints.Max.X > maxwidth {
				gtx.Constraints.Max.X = maxwidth
			}
			return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return a.layoutOption(gtx, index)
			})
		})
	})
}

func (a *App) layoutOption(gtx layout.Context, i int) layout.Dimensions {
	if i >= len(a.CurrentQuestion.Options) {
		return layout.Dimensions{}
	}

	btn := material.Button(a.Theme, &a.OptionButtons[i], a.CurrentQuestion.Options[i])
	btn.TextSize = unit.Sp(18)
	btn.Inset = layout.Inset{
		Top:    unit.Dp(12),
		Bottom: unit.Dp(12),
		Left:   unit.Dp(35),
		Right:  unit.Dp(35),
	}
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
