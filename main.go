package main

import (
	"fmt"
	"final-project/scheduler"
)


func main() {

	//welcome message to the program
	fmt.Println("Welcome to the Worker Pool Jobs Scheduler!\n")

	//prompts the user to select a number of workers
	fmt.Printf("Enter how many workers you want to use: ")
	var numOfWorkers int
	fmt.Scan(&numOfWorkers)

	//creates a new scheduler with the inputed number of workers
	s := scheduler.NewScheduler(numOfWorkers)
	s.StartWorkers()

	//prompts the user to enter computational equations and creates a job object for each one
	fmt.Println("\nNow you can enter as many equations as you want, for example: 1+2-9*8\nWhen you are done type 'DONE' \nThe program only works for the following operators: '+', '-', '*', '/'\n")
	var curEquation string
	var curJobID int = 1
	for {
		fmt.Printf("Enter Equation Here: ")
		fmt.Scan(&curEquation)
		if curEquation == "DONE"{
			break
		}
		s.CreateJob(curEquation, curJobID)
		curJobID++
	}

	fmt.Println("\nThanks for inputing Equations! Now the workers will begin completing the jobs.\nNOTE: Each job duration is determined by the number of operators, each operator takes one extra second to compute.\n")

	//calls the scheduler to begin solving the equations using the workers
	s.RunWorkers()
	s.Stop()
}