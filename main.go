package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Job struct {
	id       int
	randomno int
}

type Result struct {
	job         Job
	finalresult int
	worker      int
}

var jobs = make(chan Job, 10)
var results = make(chan Result, 10)

func digits(number int) int {
	sum := 0
	no := number
	for no != 0 {
		digit := no % 10
		sum += digit
		no /= 10
	}
	return sum
}

func worker(i int, wg *sync.WaitGroup) {
	for job := range jobs {
		output := Result{job, digits(job.randomno), i}
		time.Sleep(1 * time.Second)
		results <- output
	}
	wg.Done()
}

func createWorker(noOfWorker int) {
	var wg sync.WaitGroup
	for i := 0; i < noOfWorker; i++ {
		wg.Add(1)
		go worker(i, &wg)
	}
	wg.Wait()
	close(results)
}

func inputJob(noOfJobs int) {
	for i := 0; i < noOfJobs; i++ {
		randnum := rand.Intn(999)
		job := Job{i, randnum}
		jobs <- job
	}
	close(jobs)
}

func result(done chan bool) {
	for result := range results {
		fmt.Printf("Worker id %d, Job id %d, input random no %d , sum of digits %d\n", result.worker, result.job.id, result.job.randomno, result.finalresult)
	}
	done <- true
}

func main() {
	starttime := time.Now()
	noOfJobs := 100
	go inputJob(noOfJobs)
	done := make(chan bool)
	go result(done)
	noOfWorker := 10
	createWorker(noOfWorker)
	fmt.Println(<-done)
	endTime := time.Now()
	diff := endTime.Sub(starttime)
	fmt.Println("total time taken", diff.Seconds(), "seconds")
}
