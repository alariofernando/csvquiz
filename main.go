package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

type quiz struct {
	questions []problem
	score     int
	total     int
}

var file string
var limit int

func main() {
	quizz := newQuiz(file)
	c1 := make(chan quiz)
	timeoutFlag := time.After(time.Duration(limit) * time.Second)
	go func() {
		quizz.start()
		quizz.results()
		c1 <- quizz
	}()
	select {
	case <-c1:
	case <-timeoutFlag:
		fmt.Println("TIMEOUT!!!")
		quizz.results()
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func newQuiz(path string) quiz {
	file, err := os.Open(path)
	check(err)
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	check(err)
	var questions []problem
	var quiz quiz
	for _, val := range data {
		question := problem{question: val[0], answer: val[1]}
		question.cleanUp()
		questions = append(questions, question)
		quiz.total++
	}
	quiz.questions = questions
	return quiz
}

func (q *quiz) start() {
	for _, val := range q.questions {
		problem, answer := val.question, val.answer
		fmt.Printf("Question: %v \n", problem)
		var ans string
		fmt.Printf("Answer:  ")
		fmt.Scanf("%s \n", &ans)
		if ans == answer {
			q.score++
		}
	}
}

func (q *quiz) results() {
	fmt.Printf("You score %v out of %v !!", q.score, q.total)
}
func (p *problem) cleanUp() {
	p.question = clean(p.question)
	p.answer = clean(p.answer)
}

func clean(s string) string {
	s = strings.TrimPrefix(s, " ")
	s = strings.TrimSuffix(s, " ")
	s = strings.ToLower(s)
	return s
}

func init() {
	flag.StringVar(&file, "file", "problems.csv", "A file with problems and Answers in CSV format.")
	flag.IntVar(&limit, "timeout", 30, "Timeout time in seconds for the quiz to end.")
	flag.Parse()
}
