package main

/*
The functions in this file set up the connections to the
PosgreSQL and Elasticsearch servers and make sure they're good.
*/

import (
	"database/sql" // perform SQL calls
	"log"          // logging package
	"os"           // for getting environment variables

	_ "github.com/lib/pq"                     // PostgreSQL driver
	"github.com/shawnmilo/index_lore/elastic" // Elasticsearch driver
)

var db *sql.DB         // global variable for PostgreSQL DB connection
var es *elastic.Client // global variable for Elasticsearch connection
var index string       // name of Elasticsearch index.

// init always runs first when a program is executed -- before main(). This
// is the place to do any setup.
func init() {
	// Set things up here. If any of this fails, the program will quit.
	connectToPostgres()
	connectToElasticsearch()
	populateLearningResourceLookup()
}

// connectToPostgres sets the global db variable to a PostgreSQL connection.
func connectToPostgres() {
	url := os.Getenv("DATABASE_URL")
	if os.Getenv("LORE_DB_DISABLE_SSL") == "True" {
		url += "?sslmode=disable"
	}
	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Unable to connect to the databate at URL '%s': %s\n", url, err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Unable to communicate with the databate at URL '%s': %s\n", url, err)
	}
}

// connectToElasticsearch sets the global es variable to an Elasticsearch connection.
func connectToElasticsearch() {
	var err error
	esURL := os.Getenv("HAYSTACK_URL")
	es, err = elastic.NewClient(elastic.SetURL(esURL))
	if err != nil {
		log.Fatalf("Unable to connect to Elasticsearch at URL %s: %s\n", esURL, err)
	}
	// Ping() ignores the client's URL and defaults to 127.0.0.1, so
	// it must be specified here.
	_, _, err = es.Ping().URL(esURL).Do()
	if err != nil {
		log.Fatal("Unable to ping Elasticsearch: ", err)
	}
	index = os.Getenv("HAYSTACK_INDEX")
	if index == "" {
		index = "haystack"
	}
}
