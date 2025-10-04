package main

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/widget/material"
)

func (a *App) Layout(gtx layout.Context, th *material.Theme) {
	switch a.State {
	case ScreenInput:
		a.layoutInputScreen(gtx, th)
	case ScreenQuiz:
		a.layoutQuizScreen(gtx, th)
	}
}

func (a *App) layoutInputScreen(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{
			Axis:    layout.Vertical,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(material.H6(th, "Enter your vocabulary words (comma-separated):").Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return material.Editor(th, &a.WordInput, "e.g. innovate, challenge, create").Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				btn := material.Button(th, &a.StartButton, "Start")
				for a.StartButton.Clicked(gtx) {
					a.WordBank = parseWords(a.WordInput.Text())
					if len(a.WordBank) > 0 {
						a.State = ScreenQuiz
					}
				}
				return btn.Layout(gtx)
			}),
		)
	})
}

func parseWords(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, r := range raw {
		w := strings.TrimSpace(r)
		if w != "" {
			out = append(out, w)
		}
	}
	return out
}

// TODO
// Add new line for list format when vocab words get entered
