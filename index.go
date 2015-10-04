package main

import (
	"log"
	"strconv" // string conversion utilities

	"github.com/shawnmilo/index_lore/elastic" // Elasticsearch driver
)

/*
Functions for indexing into Elasticsearch.
*/

const maxRecs = 250                    // maximum number of resources to index at once
const mappingName = "learningresource" // name of Elasticsearch mapping

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
		request := elastic.NewBulkIndexRequest().Index("haystack").Type(mappingName).Id(id).Doc(doc)
		service.Add(request)
	}
	_, err := service.Do()
	if err != nil {
		log.Fatal("Failed to index resources: ", err)
	}
}

// ensureMapping makes sure the index and mapping exist
func ensureMapping() {
	exists, err := elastic.NewIndexExistsService(es).Index(index).Do()
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err := es.CreateIndex(index).Do()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("created index %s\n", index)
	}

	exists, err = es.TypeExists().Type(mappingName).Do()
	if !exists {
		generatedMapping := getMapping().serializable()
		_, err = elastic.NewPutMappingService(es).Type(mappingName).BodyJson(generatedMapping).Do()
		if err != nil {
			log.Fatalf("failed to create mapping %s: %s", mappingName, err)
		}
		log.Printf("created mapping %s\n", mappingName)

		log.Fatal("quitting here for your enjoyment")
	}
}
