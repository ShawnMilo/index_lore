package main

var resourceTypeQuery = `
    SELECT id, name
    FROM learningresources_learningresourcetype
    ORDER BY id`

var repositoryQuery = `
    SELECT id, name, slug
    FROM learningresources_repository
    ORDER BY id`

var repositoryVocabularyQuery = `
	SELECT id, name
    FROM taxonomy_vocabulary
    WHERE repository_id=$1`

var vocabularyQuery = `
	SELECT name
    FROM taxonomy_vocabulary`

var termQuery = `
    SELECT id, label
    FROM taxonomy_term
    WHERE vocabulary_id=$1`

var courseQuery = `
	SELECT id, org, run, course_number
    FROM learningresources_course
    WHERE repository_id=$1`

var learningResourceQuery = `
	SELECT id, title, description, content_xml, learning_resource_type_id,
           xa_nr_views, xa_nr_attempts, xa_avg_grade
    FROM learningresources_learningresource
    WHERE course_id = $1 AND id > $2
    ORDER BY id
    LIMIT $3`

var termMappingQuery = `
    SELECT term_id, learningresource_id
    FROM taxonomy_term_learning_resources
    WHERE learningresource_id >= $1 AND learningresource_id <= $2`
