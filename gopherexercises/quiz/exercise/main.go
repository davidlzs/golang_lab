package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type quizItem struct {
	index    string
	question string
	answer   string
}

type quiz struct {
	filePath     string
	questions    []quizItem
	duration     time.Duration
	correctCount int
}

func main() {
	quizFilePathPtr := flag.String("f", "problems.csv", "Quiz problem file path")
	quizTimeLimitPtr := flag.Int("timelimit", 30, "Quiz time limit in seconds")
	flag.Parse()

	quiz := quiz{
		filePath: *quizFilePathPtr,
		duration: time.Duration(*quizTimeLimitPtr) * time.Second,
	}

	(&quiz).Start()
}

func (quiz *quiz) Start() {
	quiz.loadQuiz(quiz.filePath)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Presss 'Enter' to start the quiz")
	reader.ReadString('\n')

	timer := time.NewTimer(quiz.duration)

	ch := make(chan string)
	go quiz.run(ch)

	for {
		select {
		case <-timer.C:
			ch <- "TimeOut"
		case msg := <-ch:
			if strings.Compare("Done", msg) == 0 {
				fmt.Println("Done")
				os.Exit(0)
			} else if strings.Compare("TimeOut", msg) == 0 {
				fmt.Printf("Time Out. %v.\n", quiz.duration)
				os.Exit(1)
			}
		}
	}
}

func (quiz *quiz) run(ch chan string) {
	stdinCh := make(chan string)
	defer func() {
		close(ch)
		close(stdinCh)
	}()
	reader := bufio.NewReader(os.Stdin)
	for _, question := range quiz.questions {

		go func(ch chan string) {
			text, _ := reader.ReadString('\n')
			ch <- text
		}(stdinCh)

		fmt.Printf("Question#%s %s=", question.index, question.question)

		select {
		case <-ch:
			fmt.Printf("\nYou scored %d out of %d\n", quiz.correctCount, len(quiz.questions))
			ch <- "TimeOut"
			return
		case answer := <-stdinCh:
			if strings.Compare(question.answer, strings.TrimSuffix(answer, "\n")) == 0 {
				quiz.correctCount++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", quiz.correctCount, len(quiz.questions))
	ch <- "Done"
}

func (quiz *quiz) loadQuiz(filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	for index, record := range records {
		item := quizItem{
			index:    fmt.Sprintf("%d", index+1),
			question: record[0],
			answer:   record[1],
		}
		quiz.questions = append(quiz.questions, item)
	}
}
