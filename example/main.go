package main

import (
	"log"
	"os"

	"github.com/hiroakis/jr"
)

func main() {
	f, err := os.Open("../Jobfile.example.yml")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	jobs, err := jr.LoadJob(f)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, job := range jobs.Jobs {
		if err := job.Run(); err != nil {
			log.Fatal(err)
			return
		}
	}
}
