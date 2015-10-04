package main

import "log"

// main is the first code run after any init() functions
func main() {
	count := 0
	for _, repo := range getRepositories() {
		for _, course := range getCourses(repo) {
			count += indexCourse(course)
		}
	}
	log.Printf("Indexed %d learning resources.\n", count)
}
