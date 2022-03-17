package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main(){
	csvfilename := flag.String("csv", "problems.csv", "adding flag")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvfilename)
	if err != nil{
		exit(fmt.Sprintf("Failed to open csv %s", *csvfilename))
	}
	
	r := csv.NewReader(file)
	lines, err := r.ReadAll()
	if err != nil{
		exit("Failed to parse the CSV file")
	}
	problems := parseLines(lines)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)
	correct := 0

	for i, p := range problems{
		fmt.Printf("Problem #%d: %s = \n", i+1, p.q)
		answerCh := make(chan string)
		go func(){
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select{
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
}

func parseLines(lines [][]string) []problem{
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem {
			q: line[0], 
			a: line[1],
		}
	}
	return ret
}

type problem struct{
	q string
	a string
}

func exit(msg string){
	fmt.Print(msg)
	os.Exit(1)
}