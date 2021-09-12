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

func main() {
	quizFilePathPtr := flag.String("f", "problems.csv", "Quiz problem file path")
	quizTimeLimitPtr := flag.Int("timelimit", 30, "Quiz time limit in seconds")

	questions, err := loadQuiz(*quizFilePathPtr)
	if err != nil {
		panic(err)
	}

	timer := time.NewTimer(time.Duration(*quizTimeLimitPtr) * time.Second)

	ch := make(chan string)
	go run(ch, questions)

	for {
		select {
		case <-timer.C:
			ch <- "TimeOut"
		case msg := <-ch:
			if strings.Compare("Done", msg) == 0 {
				fmt.Println("Done")
				os.Exit(0)
			} else if strings.Compare("TimeOut", msg) == 0 {
				fmt.Printf("Time Out. %d seconds.\n", *quizTimeLimitPtr)
				os.Exit(1)
			}
		}
	}
}

func run(ch chan string, questions []quizItem) {
	var correctCount int
	reader := bufio.NewReader(os.Stdin)
	for _, question := range questions {
		stdinCh := make(chan string)
		go func(ch chan string) {
			text, _ := reader.ReadString('\n')
			ch <- text
		}(stdinCh)

		fmt.Printf("Question#%s %s=", question.index, question.question)

		select {
		case <-ch:
			fmt.Printf("\nYou scored %d out of %d\n", correctCount, len(questions))
			ch <- "TimeOut"
			return
		case answer := <-stdinCh:
			if strings.Compare(question.answer, strings.TrimSuffix(answer, "\n")) == 0 {
				correctCount++
			}
		}
	}

	fmt.Printf("You scored %d out of %d\n", correctCount, len(questions))
	ch <- "Done"
}

func loadQuiz(filePath string) ([]quizItem, error) {
	var quizItems []quizItem
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, err
	}

	for index, record := range records {
		item := quizItem{
			index:    fmt.Sprintf("%d", index+1),
			question: record[0],
			answer:   record[1],
		}
		quizItems = append(quizItems, item)
	}

	return quizItems, nil
}
