package main

import (
	"log"
	"time"
)

func (a *App) startPreloading() {
	go func() {
		defer close(a.PreloadDone)
		for {
			if len(a.QuestionBuffer) >= cap(a.QuestionBuffer) {
				time.Sleep(50 * time.Millisecond)
				continue
			}
			preloadIndex := a.CurrentIndex + 1
			var question Question
			var err error
			time.Sleep(3 * time.Second)

			if preloadIndex < 10 {
				question, err = buildSingleWordQuestion(a.WordBank, preloadIndex, a.APIKey, a.RNG)
			} else {
				question, err = buildTwoWordQuestion(a.WordBank, a.APIKey, a.RNG)
			}
			if err != nil {
				log.Println("Error preloading question:", err)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			a.QuestionBuffer <- question
		}
	}()
}

func (a *App) loadNextQuestion() {
	for len(a.QuestionBuffer) == 0 {
		time.Sleep(3 * time.Second)
	}
	a.CurrentQuestion = <-a.QuestionBuffer
	a.ErrorMsg = ""
	a.State = ScreenQuiz
	a.Window.Invalidate()
}

/* Helper function to generate question based on current index
func (a *App) generateQuestion() (Question, error) {
	if a.CurrentIndex < 10 {
		return buildSingleWordQuestion(a.WordBank, a.CurrentIndex, a.APIKey, a.RNG)
	}
	return buildTwoWordQuestion(a.WordBank, a.APIKey, a.RNG)
}
*/
