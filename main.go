package main

import (
	"log"
	"time"
)

// main is the first code run after any init() functions
func main() {
	start := time.Now()
	count := 0
	for _, repo := range getRepositories() {
		for _, course := range getCourses(repo) {
			count += indexCourse(course)
		}
	}
	duration := time.Now().Sub(start)
	log.Printf("Indexed %d learning resources in %0.2f seconds.\n", count, duration.Seconds())
}
