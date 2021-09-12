package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type quizItem struct {
	index    string
	question string
	answer   string
}

func main() {
	quizFilePath := "problems.csv"
	if len(os.Args) > 1 {
		if strings.Compare(os.Args[1], "-f") != 0 {
			fmt.Println("quiz -f <problems.csv>")
			os.Exit(1)
		} else {
			quizFilePath = os.Args[2]
		}
	}

	questions, err := loadQuiz(quizFilePath)
	if err != nil {
		panic(err)
	}

	var correctCount int
	reader := bufio.NewReader(os.Stdin)
	for _, question := range questions {
		fmt.Printf("Question#%s %s=", question.index, question.question)
		text, _ := reader.ReadString('\n')
		if strings.Compare(question.answer, strings.TrimSuffix(text, "\n")) == 0 {
			correctCount++
		} else {
		}
	}

	fmt.Printf("You scored %d out of %d\n", correctCount, len(questions))

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
