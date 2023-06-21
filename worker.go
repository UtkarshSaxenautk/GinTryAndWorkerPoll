package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	ID   int
	Func func()
}

func worker(id int, tasks <-chan Task, wg *sync.WaitGroup) {
	for task := range tasks {
		fmt.Println("Worker ", id, " started task ", task.ID)
		task.Func() // Execute the task function
		fmt.Println("Worker ", id, " finished task ", task.ID)
		wg.Done()
	}
}

func taskFunction(id int) func() {
	taskFunc := func() {
		fmt.Println("Started Executing task function with id: ", id, "will take ", id*2, " seconds")
		time.Sleep(time.Duration(id*2) * time.Second)
	}
	return taskFunc
}

func main() {
	var numTasks, numWorkers int

	fmt.Print("Enter the number of tasks: ")
	_, err := fmt.Scanln(&numTasks)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	fmt.Print("Enter the number of workers in the worker pool: ")
	_, err = fmt.Scanln(&numWorkers)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	tasks := make(chan Task, numTasks)
	var wg sync.WaitGroup

	startTime := time.Now()

	// start worker initially channel is empty
	for i := 1; i <= numWorkers; i++ {
		go worker(i, tasks, &wg)
	}

	// Defining the task functions with their id

	// Add tasks to the channel
	for i := 1; i <= numTasks; i++ {
		wg.Add(1)
		tasks <- Task{ID: i, Func: taskFunction(i)}
	}

	close(tasks)

	// Wait for all tasks to be completed
	wg.Wait()

	// Calculate total time taken
	totalTime := time.Since(startTime)

	fmt.Println("All tasks completed. Total time taken : ", totalTime)
}
