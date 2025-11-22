# Vocabulary Test Time

Vocabulary Test Time is a desktop quiz app built with [Gio](https://gioui.org/) that helps you practice vocabulary using AI‑generated sentences and multiple‑choice questions.

You enter a small word list, and the app generates quiz questions that use those words in context. It then tracks your score over a fixed number of questions.

---

## Features

- 📝 Enter a custom list of **10 vocabulary words**
- 🧠 AI‑generated sentences that use your words in context
- ❓ Multiple‑choice questions with 4 options each
- ✅ Immediate feedback after each answer
- 📊 Final score out of 15 questions, with percentage
- 🔁 Restart the quiz with the same vocabulary list
- ⚡ Background preloading of upcoming questions for smoother UX

---

## How it Works

1. On launch, the app shows an input screen:
   - Type **exactly 10 words**, separated by commas or new line (press enter).
   - Click **“Start Quiz”**.

2. The app:
   - Builds the **first question** immediately from your word list.
   - Starts a **background preloader** that:
     - Preloads questions into a buffer channel.
     - Generates:
       - Single‑word questions for indices 0–9.
       - Two‑word questions for indices 10–14.
     - Stops after 15 total questions.

3. During the quiz:
   - A sentence is shown plus 4 answer options.
   - Clicking an option:
     - Shows feedback (`Correct!` or the correct answer).
     - After ~2 seconds, advances to the next question.
   - While loading the next question, a **“Loading Sentence…”** screen may appear briefly.

4. After 15 questions:
   - The **results screen** shows:
     - Raw score (e.g. `12 / 15`)
     - Percentage (e.g. `80.0%`)
   - You can click **“Restart Quiz”** to:
     - Reset score and state.
     - Regenerate new questions using the same word bank.

---

## Tech Stack

- Language: **Go**
- UI: **[gioui.org](https://gioui.org)**
- Fonts: **gofont** collection
- Config: **dotenv** via [`github.com/joho/godotenv`](https://github.com/joho/godotenv)
- Randomness: `math/rand` with time‑based seeding
- Concurrency:
  - Background question preloading with goroutines
  - Buffered `chan Question` for question pipeline

---

## Requirements

- Go 1.20+ 
- An AI API key (e.g., OpenAI) exposed via environment variable:
  - `.env` file in the project root containing:

    ```bash
    SecretKey=your_api_key_here
    ```

- OS:
  - Any platform supported by Gio (Linux, macOS, Windows, etc.)

---

## Installation

```bash
git clone https://github.com/Alford05/StudyBuddy.git
cd StudyBuddy

# Copy .env.example to .env if you provide one, or create .env manually
echo "SecretKey=your_api_key_here" > .env

go mod tidy
go run .
```

---

## Usage 

Launch the app
On the input screen:
    Enter exactly 10 words (for now) 
Click Start Quiz
Answer each question by clicking 1 of 4 options
View your score at the end and optionally restart with new questions generated from the same words

---

## Environment Variables 
- 'SecretKey' - API key for the language model the app uses internally to generate questions 
Loaded via:

```GO
    err := godotenv.Load()
    apiKey := os.Getenv("SecretKey")
if SecretKey is missing, the app logs:
    SecretKey variable not set.
```

---

## Architecture Overview

    App struct:
        Holds UI state (State, Theme, WordInput, buttons)
        Quiz data (WordBank, CurrentQuestion, CurrentIndex, Score)
        Feedback (FeedbackMsg, ShowFeedback, ErrorMsg)
        Concurrency primitives:
            QuestionBuffer chan Question
            PreloadDone chan struct{}
            StopPreload chan struct{}
    Screens:
        ScreenInput – enter 10 words and start quiz
        ScreenLoading – transient loading view while waiting for next question
        ScreenQuiz – question, options, feedback
        QuizComplete – results and restart button
    Question generation (external helpers you provide):
        buildSingleWordQuestion(...)
        buildTwoWordQuestion(...)
    Preloading:
        startPreloading runs in a goroutine
        Fills QuestionBuffer up to capacity
        Stops when:
            StopPreload is closed, or
            15 total questions are generated

---

TODO / Future Ideas

    Add settings for:
        Number of questions
        Difficulty level
    Support synonyms / definitions / fill‑in‑the‑blank modes
    Persist scores and word lists
    Package as downloadable binaries for major OSes
