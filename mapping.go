package main

const (
	// All strings are analyzed by default, so no "Analyzed" constant is used.
	NoIndex   = "no"           // only stored and retrieved; not used for searches
	NoAnalyze = "not_analyzed" // only exact matches can be made

	FloatType   = "float"
	StringType  = "string"
	IntegerType = "integer"
)

// FieldOptions are options such as "index": "not_analyzed" or "type": "string"
type FieldOptions map[string]string

//Mapping is a container for Fields
type Mapping map[string]map[string]FieldOptions

// NewMapping returns a pointer to an empty Mapping, with
// just the mapping name set
func NewMapping(name string) *Mapping {
	return &Mapping{"properties": map[string]FieldOptions{}}
}

// AddField adds a field to the mapping.
func (m Mapping) AddField(name, index, fieldType, nullValue string) {
	opts := FieldOptions{}
	if index != "" {
		opts["index"] = index
	}
	if fieldType != "" {
		opts["type"] = fieldType
	}
	if nullValue != "" {
		opts["nullValue"] = nullValue
	}
	m["properties"][name] = opts
}

// AddStoredString adds a string field to the index which will
// not be searchable, but will be returned with search results
func (m *Mapping) AddStoredString(name string) {
	m.AddField(name, NoIndex, StringType, "")
}

// AddExactString adds a string field to the index which will be
// searchable only by an exact match
func (m *Mapping) AddExactString(name string) {
	m.AddField(name, NoAnalyze, StringType, "")
}

// AddSearchableString adds a string field to the index which will be
// searchable by partial matches
func (m *Mapping) AddSearchableString(name string) {
	m.AddField(name, "", StringType, "")
}

// AddFloat adds a float field to the index.
func (m *Mapping) AddFloat(name string) {
	m.AddField(name, "", FloatType, "")
}

// AddInteger adds an integer field to the index.
func (m *Mapping) AddInteger(name string) {
	m.AddField(name, "", IntegerType, "")
}

// serializable returns a map[string]interface{} as required by
// the elastic package for sending a mapping defition to Elasticsearch.
func (m *Mapping) serializable() map[string]interface{} {

	return map[string]interface{}{mappingName: m}

}

func getMapping() *Mapping {

	m := NewMapping(mappingName)

	storedStrings := []string{"content_xml", "description_path", "preview_url"}
	exactStrings := []string{"course", "run", "repository", "resource_type", "titlesort"}
	searchableStrings := []string{"description", "content_stripped", "title"}
	floatFields := []string{"xa_avg_grade", "xa_histogram_grade"}
	intFields := []string{"xa_nr_views", "xa_nr_attempts"}

	for _, name := range storedStrings {
		m.AddStoredString(name)
	}
	for _, name := range exactStrings {
		m.AddExactString(name)
	}
	for _, name := range searchableStrings {
		m.AddSearchableString(name)
	}
	for _, name := range floatFields {
		m.AddFloat(name)
	}
	for _, name := range intFields {
		m.AddInteger(name)
	}

	for _, name := range getVocabularyNames() {
		m.AddExactString(name)
	}

	return m
}

/*

This is an example of the JSON created by a Mapping type.

The format is specificied by Elasticsearch:
https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-put-mapping.html

"learningresource": {
    "properties": {
        "run": {
            "index": "not_analyzed",
            "type": "string"
        },
        "description": {
            "type": "string"
        },
        "repository": {
            "index": "not_analyzed",
            "type": "string"
        },
        "xa_histogram_grade": {
            "type": "float"
        },
        "course": {
            "index": "not_analyzed",
            "type": "string"
        },
        "content_stripped": {
            "type": "string"
        },
        "xa_avg_grade": {
            "type": "float"
        },
        "title": {
            "type": "string"
        },
        "description_path": {
            "index": "no",
            "type": "string"
        },
        "content_xml": {
            "index": "no",
            "type": "string"
        },
        "xa_nr_views": {
            "type": "integer"
        },
        "preview_url": {
            "index": "no",
            "type": "string"
        },
        "xa_nr_attempts": {
            "type": "integer"
        },
        "titlesort": {
            "index": "not_analyzed",
            "type": "string"
        },
        "resource_type": {
            "index": "not_analyzed",
            "type": "string"
        }
    }
}
*/
