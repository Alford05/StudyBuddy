# Vocabulary Test Time

Vocabulary Test Time is a desktop quiz app built with [Gio](https://gioui.org/) that helps you practice vocabulary using AI‑generated sentences and multiple‑choice fill-in-the-blanks.

You enter a small word list, and the app generates quiz questions that use those words in context. It then tracks your score over a fixed number of questions.

---

## Motivation

Vocabulary Test Time is designed to make vocabulary practice engaging and effective. It tackles the challenge of learning new words in context by leveraging AI to generate unique sentences and multiple-choice quizzes. This was created initially to assist my daughter with her vocabulary homework.Whether you're a student preparing for tests or just looking to expand your lexicon, this app provides a dynamic and personalized way to reinforce your understanding of new words.

---

## Features

- 📝 Enter a custom list of **10 vocabulary words**
- 🧠 AI‑generated sentences that use your words in context
- ❓ Multiple‑choice questions with 4 options each
- ✅ Immediate feedback after each answer
- 📊 Final score out of 15 questions, with percentage
- 🔁 Restart the quiz with the same vocabulary list, but new questions 
- ⚡ Background preloading of upcoming questions for smoother UX

---

## Quick Start

Follow these steps to get Vocabulary Test Time up and running on your local machine.

---

### Requirements

- Go 1.20+ 
- An AI API key (e.g., OpenAI) exposed via environment variable:
  - `.env` file in the project root containing:

    ```bash
    SecretKey=your_api_key_here
    ```

- OS:
  - Any platform supported by Gio (Linux, macOS, Windows, etc.)

---

### Installation

```bash
git clone https://github.com/Alford05/StudyBuddy.git
cd StudyBuddy

# Copy .env.example to .env if you provide one, or create .env manually
echo "SecretKey=your_api_key_here" > .env

go mod tidy
go run .
```

---

### Environment Variables 
- 'SecretKey' - API key for the language model the app uses internally to generate questions 
Loaded via:

```GO
    err := godotenv.Load()
    apiKey := os.Getenv("SecretKey")
```
if SecretKey is missing, the app logs:
    SecretKey variable not set.

---

## Usage 

Launch the app
On the input screen:
    Enter exactly 10 words (for now) 
Click Start Quiz
Answer each question by clicking 1 of 4 options
View your score at the end and optionally restart with new questions generated from the same words

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

## Contributing

### Suggesting Enhancements

For new feature ideas or significant changes, it's often a good idea to open an issue first to discuss it. This allows us to align on the approach before you invest time in coding.

### Development Setup

To get your development environment ready:

1.  **Fork** the [Vocabulary Test Time repository](https://github.com/Alford05/StudyBuddy) to your GitHub account.
2.  **Clone** your forked repository:
    ```bash
    git clone https://github.com/YOUR_USERNAME/StudyBuddy.git
    cd StudyBuddy
    ```
3.  Ensure you have [Go 1.20+](https://go.dev/doc/install) installed.
4.  Set up your `SecretKey` environment variable as described in the [Quick Start](#quick-start) section.
5.  Install dependencies:
    ```bash
    go mod tidy
    ```
6.  Run the application to test your setup:
    ```bash
    go run .
    ```

### Making Changes

1.  Create a new branch for your feature or bug fix:
    ```bash
    git checkout -b feature/your-feature-name
    # OR
    git checkout -b bugfix/fix-description
    ```
2.  Make your changes to the codebase.
3.  Please ensure your code follows the existing style and conventions.
4.  If applicable, add or update tests for your changes. (If you have tests, you'd add instructions here on how to run them, e.g., `go test ./...`)
5.  Commit your changes with a clear and descriptive commit message.

### Submitting a Pull Request

1.  Push your branch to your forked repository:
    ```bash
    git push origin feature/your-feature-name
    ```
2.  Open a Pull Request from your forked repository to the `main` branch of the original [Vocabulary Test Time repository](https://github.com/Alford05/StudyBuddy).
3.  Provide a clear description of your changes in the Pull Request. Refer to any related issues by their number (e.g., `Fixes #123`).

Thank you for your contributions!

---

## TODO / Future Ideas

    Add settings for:
        Number of questions
        Difficulty level
    Support synonyms / definitions / fill‑in‑the‑blank modes
    Persist scores and word lists
    Package as downloadable binaries for major OSes
