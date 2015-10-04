package main

/*
This file contains the interactions with the LORE PostgreSQL database.
*/

import (
	"bytes"
	"log"
	"strings"
)

// This is global because it's not specific to any repository. It allows
// quick retrieval of the name by primary key.
var resourceTypeLookup map[int]string

// populateLearningResourceLookup populates global variable resourceTypeLookup
func populateLearningResourceLookup() {
	// This must be initialized; it was only declared above.
	resourceTypeLookup = make(map[int]string)

	rows, err := db.Query(resourceTypeQuery)
	if err != nil {
		log.Fatal("Unable to query learning resource types:", err)
	}
	defer rows.Close()

	// Next() moves the cursor forward, returns true if row found, else false.
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		resourceTypeLookup[id] = name
	}
}

func getRepositories() []*Repository {
	rows, err := db.Query(repositoryQuery)
	if err != nil {
		log.Fatal("Unable to query repositories: ", err)
	}
	defer rows.Close()

	repos := make([]*Repository, 0, 10)

	// Next() moves the cursor forward, returns true if row found, else false.
	for rows.Next() {
		repo := Repository{}
		rows.Scan(&repo.id, &repo.name, &repo.slug)
		getVocabularies(&repo)
		repos = append(repos, &repo)
	}
	return repos
}

type Repository struct {
	id           int
	name         string
	slug         string
	vocabularies map[int]string
	terms        map[int]*Term
}

// Declaring a String() method on any type will
// cause that method to be used any time an
// instance of the type is printed.
func (r Repository) String() string {
	return r.name
}

// getVocabularies gets all vocabularies and terms for a Repository.
func getVocabularies(repo *Repository) {
	// get vocabularies
	rows, err := db.Query(vocabularyQuery, repo.id)
	if err != nil {
		log.Fatal("Unable to query vocabularies: ", err)
	}
	defer rows.Close()

	vocabs := make(map[int]string)

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		vocabs[id] = name
	}
	repo.vocabularies = vocabs

	// for each Vocabulary, get terms
	terms := make(map[int]*Term)
	for i := range repo.vocabularies {
		rows, err = db.Query(termQuery, i)
		if err != nil {
			log.Fatal("Unable to query terms: ", err)
		}
		defer rows.Close()
		for rows.Next() {
			var id int
			term := Term{}
			rows.Scan(&id, &term.label)
			// We're iterating over the vocabulary IDs, so we don't
			// need to pull the term's vocabulary_id from the database; it's "i."
			term.vocabulary_id = i
			terms[id] = &term
		}
	}
	repo.terms = terms
}

type Term struct {
	label         string
	vocabulary_id int
}

func getCourses(repo *Repository) []*Course {
	rows, err := db.Query(courseQuery, repo.id)
	if err != nil {
		log.Fatal("Unable to query courses: ", err)
	}
	defer rows.Close()

	courses := make([]*Course, 0, 10)

	for rows.Next() {
		c := Course{}
		rows.Scan(&c.id, &c.org, &c.run, &c.course_number)
		c.repository = repo
		courses = append(courses, &c)
	}
	return courses
}

type Course struct {
	id            int
	org           string
	run           string
	course_number string
	repository    *Repository
}

type LearningResource struct {
	id              int
	course          *Course
	run             string
	repository      *Repository
	title           string
	description     string
	descriptionPath string
	contentXML      string
	xaNumViews      int
	xaNumAttempts   int
	xaAvgGrade      float32
	previewURL      string
	resourceTypeId  int
	terms           termSet
}

func (lr *LearningResource) toSerializable() map[string]interface{} {

	dump := map[string]interface{}{
		"id":               lr.id,
		"_id":              lr.id,
		"course":           lr.course.course_number,
		"run":              lr.course.run,
		"repository":       lr.course.repository.slug,
		"resource_type":    resourceTypeLookup[lr.resourceTypeId],
		"title":            lr.title,
		"description":      lr.description,
		"description_path": lr.descriptionPath,
		"content_xml":      lr.contentXML,
		"xa_nr_views":      lr.xaNumViews,
		"xa_nr_attempts":   lr.xaNumAttempts,
		"xa_avg_grade":     lr.xaAvgGrade,
		"preview_url":      lr.previewURL,
	}
	dump["content_stripped"] = stripXML(lr.contentXML)

	titleSort := strings.TrimSpace(lr.title)
	if titleSort == "" {
		titleSort = "1"
	} else {
		titleSort = "0" + titleSort
	}
	dump["titlesort"] = titleSort

	for vocab_id, term_ids := range lr.terms {
		names := make([]string, 0, len(term_ids))
		for _, id := range term_ids {
			names = append(names, lr.course.repository.terms[id].label)
		}
		dump[lr.course.repository.vocabularies[vocab_id]] = names
	}

	return dump

}

type termSet map[int][]int

func getLearningResources(course *Course, lastID int) []*LearningResource {
	rows, err := db.Query(learningResourceQuery, course.id, lastID, maxRecs)
	if err != nil {
		log.Fatal("Unable to query learning resources:", err)
	}
	defer rows.Close()
	resources := make([]*LearningResource, 0, 100)
	for rows.Next() {
		resource := LearningResource{}
		rows.Scan(
			&resource.id, &resource.title, &resource.description, &resource.contentXML,
			&resource.resourceTypeId, &resource.xaNumViews, &resource.xaNumAttempts,
			&resource.xaAvgGrade,
		)
		resource.run = course.run
		resource.course = course
		resource.repository = course.repository
		resources = append(resources, &resource)
	}

	return resources
}

func addTermInfo(resources []*LearningResource) {
	firstId := resources[0].id
	lastId := resources[len(resources)-1].id
	repo := resources[0].repository
	rows, err := db.Query(termMappingQuery, firstId, lastId)
	if err != nil {
		log.Fatal("Unable to query learning resource term crosswalk:", err)
	}
	defer rows.Close()

	type term struct{ term_id, vocab_id int }

	ids := make(map[int][]term)

	for rows.Next() {
		var learningresource_id int
		var t term
		rows.Scan(&t.term_id, &learningresource_id)
		t.vocab_id = repo.terms[t.term_id].vocabulary_id
		ids[learningresource_id] = append(ids[learningresource_id], t)
	}

	for _, lr := range resources {
		terms := make(termSet)
		for vocab_id := range lr.repository.vocabularies {
			// Always have the key even if there are no terms,
			// for indexing "empty" values.
			terms[vocab_id] = []int{}
		}
		for _, t := range ids[lr.id] {
			terms[t.vocab_id] = append(terms[t.vocab_id], t.term_id)
		}
		lr.terms = terms
	}
}

// stripXML removes XML tags from a string. It is copied and very slightly modified
// from https://github.com/mholt/caddy/blob/master/middleware/context.go#L136
// Thanks to @mholt for sharing this in gophers.slack.com!
func stripXML(s string) string {
	var buf bytes.Buffer
	var inTag, inQuotes bool
	var tagStart int
	for i, ch := range s {
		if inTag {
			if ch == '>' && !inQuotes {
				inTag = false
			} else if ch == '<' && !inQuotes {
				// false start
				buf.WriteString(s[tagStart:i])
				tagStart = i
			} else if ch == '"' {
				inQuotes = !inQuotes
			}
			continue
		}
		if ch == '<' {
			inTag = true
			tagStart = i
			continue
		}
		buf.WriteRune(ch)
	}
	if inTag {
		// false start
		buf.WriteString(s[tagStart:])
		inTag = false
	}
	return buf.String()
}
