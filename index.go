package main

import (
	"log"
	"strconv" // string conversion utilities

	"github.com/shawnmilo/index_lore/elastic" // Elasticsearch driver
)

/*
Functions for indexing into Elasticsearch.
*/

const maxRecs = 250 // maximum number of resources to index at once

func indexCourse(course *Course) (count int) {
	// All variables are initialized with their "zero value."
	// For integers, this is zero.
	var lastID, total int
	// Loop over all learning resources obeying the maxRecs global.
	for {
		resources := getLearningResources(course, lastID)
		count := len(resources)
		if count == 0 {
			// We're done with this course.
			break
		}

		total += count
		addTermInfo(resources)
		indexResources(resources)
		lastID = resources[len(resources)-1].id
	}
	return total
}

// indexResources updates Elasticsearch.
func indexResources(resources []*LearningResource) {
	service := elastic.NewBulkService(es)
	for _, resource := range resources {
		id := strconv.Itoa(resource.id)
		doc := resource.toSerializable()
		request := elastic.NewBulkIndexRequest().Index("haystack").Type("learningresource").Id(id).Doc(doc)
		service.Add(request)
	}
	_, err := service.Do()
	if err != nil {
		log.Fatal("Failed to index resources: ", err)
	}
}
