package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Patient struct {
	MRNumber string
	Gender   string
}

type PatientQueue struct {
	Queue           []Patient
	RoundRobinMode  bool
	CurrentPosition int
}

func isPatientExist(patients []Patient, mrNumber string) bool {
	for i := 0; i < len(patients); i++ {
		if patients[i].MRNumber == mrNumber {
			return true
		}
	}
	return false
}

func IsNumberValid(mrNumber string) bool {
	if strings.HasPrefix(mrNumber, "MR") && len(mrNumber) == 6 {
		digit := mrNumber[len(mrNumber)-4:]
		_, err := strconv.Atoi(digit)
		if err != nil {
			fmt.Println("Invalid Patient MR Number")
			return false
		}
		return true
	}
	return false
}

func (pq *PatientQueue) Enqueue(patient Patient) {
	if IsNumberValid(patient.MRNumber) {
		if !isPatientExist(pq.Queue, patient.MRNumber) {
			pq.Queue = append(pq.Queue, patient)
			fmt.Printf("Patient %s added to the queue.\n", patient.MRNumber)
			fmt.Println(pq.Queue)
		} else {
			fmt.Printf("Patient %s is already exists in the queue.\n", patient.MRNumber)
			fmt.Println(pq.Queue)
		}
	} else {
		fmt.Println("Invalid Patient MR Number")
	}
}

func (pq *PatientQueue) Dequeue() *Patient {
	if len(pq.Queue) > 0 {
		patient := pq.Queue[0]
		pq.Queue = pq.Queue[1:]
		return &patient
	}
	return nil
}

func (pq *PatientQueue) ChangeMode(mode string) {
	if mode == "ROUNDROBIN" {
		pq.RoundRobinMode = true
		pq.Queue = toggleRoundRobin(pq.Queue, true)
		fmt.Println("Queue ordering changed to Round Robin Gender.")
		fmt.Println(pq.Queue)
	} else if mode == "DEFAULT" {
		pq.RoundRobinMode = false
		pq.Queue = toggleRoundRobin(pq.Queue, false)
		pq.CurrentPosition = 0
		fmt.Println("Queue ordering changed to Default (FIFO).")
		fmt.Println(pq.Queue)
	}
}

func roundRobinByGender(input []Patient, enableRoundRobin bool) []Patient {
	if !enableRoundRobin {
		return input
	}

	maleQueue := make([]Patient, 0)
	femaleQueue := make([]Patient, 0)
	output := make([]Patient, 0)

	for _, patient := range input {
		if patient.Gender == "M" {
			maleQueue = append(maleQueue, patient)
		} else if patient.Gender == "F" {
			femaleQueue = append(femaleQueue, patient)
		}
	}

	for len(maleQueue) > 0 || len(femaleQueue) > 0 {
		if len(maleQueue) > 0 {
			output = append(output, maleQueue[0])
			maleQueue = maleQueue[1:]
		}
		if len(femaleQueue) > 0 {
			output = append(output, femaleQueue[0])
			femaleQueue = femaleQueue[1:]
		}
	}

	return output
}

func toggleRoundRobin(input []Patient, enableRoundRobin bool) []Patient {
	return roundRobinByGender(input, enableRoundRobin)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pq := PatientQueue{}

	for {
		fmt.Print("Enter command: ")
		scanner.Scan()
		command := scanner.Text()

		switch {
		case strings.HasPrefix(command, "IN"):
			inputs := strings.Fields(command)
			if len(inputs) != 3 {
				fmt.Println("Invalid input. Please try again.")
				continue
			}
			mrNumber := inputs[1]
			gender := inputs[2]
			patient := Patient{MRNumber: mrNumber, Gender: gender}
			pq.Enqueue(patient)

		case command == "OUT":
			patient := pq.Dequeue()
			if patient != nil {
				fmt.Printf("Patient %s dispatched from the queue.\n", patient.MRNumber)
				fmt.Println(pq.Queue)
			} else {
				fmt.Println("No patients in the queue.")
			}

		case command == "ROUNDROBIN":
			pq.ChangeMode("ROUNDROBIN")

		case command == "DEFAULT":
			pq.ChangeMode("DEFAULT")

		case command == "EXIT":
			fmt.Println("Exiting the application. Goodbye!")
			return

		default:
			fmt.Println("Invalid command. Please try again.")
		}
	}
}
