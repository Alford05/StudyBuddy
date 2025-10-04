package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const GeminiAPI = "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"

func buildQuestion(wordBank []string, apiKey string, rng *rand.Rand) (Question, error) {
	if len(wordBank) != 10 {
		return Question{
			Sentence:     "Please enter at least 10 words.",
			Options:      []string{"_", "_", "_", "_", "_", "_"},
			CorrectIndex: 0,
		}, fmt.Errorf("not enough words in bank")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading hidden file")
	}

	apiSecret := os.Getenv("SecretKey")
	if apiSecret == "" {
		log.Fatal("API_SECRET environment variable is not set")
	}

	correctWord := wordBank[rng.Intn(len(wordBank))]

	// Generate sentence with Gemini
	sentence, err := generateSentence(correctWord, apiSecret)
	if err != nil {
		return Question{
			Sentence:     "Error generating sentence.",
			Options:      []string{"_", "_", "_", "_", "_", "_"},
			CorrectIndex: 0,
		}, err
	}

	blankSentence := replaceFirstWord(sentence, correctWord)
	if blankSentence == sentence {
		blankSentence = strings.Replace(sentence, correctWord, "_______", 1)
	}

	options := buildMultipleChoiceOptions(correctWord, wordBank, rng)

	correctIndex := -1
	for i, opt := range options {
		if opt == correctWord {
			correctIndex = i
			break
		}
	}

	return Question{
		Sentence:     blankSentence,
		Options:      options,
		CorrectIndex: correctIndex,
		AnswerWords:  []string{correctWord},
	}, nil
}

func generateSentence(word string, apiSecret string) (string, error) {
	prompt := fmt.Sprintf("Write an English sentence based on a 7th grade reading level that uses the word \"%s\". Do not quote the word. Output only the sentence.", word)

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}
	bodyBytes, _ := json.Marshal(reqBody)

	url := fmt.Sprintf("%s?key=%s", GeminiAPI, apiSecret)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respData, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(respData))
	}

	var res struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	if len(res.Candidates) == 0 || len(res.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no sentence returned")
	}
	return res.Candidates[0].Content.Parts[0].Text, nil
}

func replaceFirstWord(sentence, word string) string {
	lowerSentence := strings.ToLower(sentence)
	lowerWord := strings.ToLower(word)

	index := strings.Index(lowerSentence, lowerWord)
	if index == -1 {
		return sentence
	}
	return sentence[:index] + "_______" + sentence[index+len(word):]
}

func buildMultipleChoiceOptions(correct string, bank []string, rng *rand.Rand) []string {
	options := make([]string, 0, 4)
	used := map[string]bool{correct: true}

	correctIndex := rng.Intn(4)
	for i := 0; i < 4; i++ {
		if i == correctIndex {
			options = append(options, correct)
		} else {
			var distractor string
			for {
				distractor = bank[rand.Intn(len(bank))]
				if !used[distractor] {
					break
				}
			}
			options = append(options, distractor)
			used[distractor] = true
		}
	}
	return options
}
