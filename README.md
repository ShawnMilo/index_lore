# lore_index

This application indexes records for the [LORE](https://github.com/mitodl/lore) system from its PostgreSQL database into its Elasticsearch instance.

# Build

```sh
go build
```

# Run

Just run the executable:

```sh
./index_lore
```

If the environment variables already used by LORE are present, the program will connect to the servers as long as PostgreSQL and Elasticsearch are running. The Django application does not need to be running.

The environment variables required are:

- `DATABASE_URL`
- `LORE_DB_DISABLE_SSL`
- `HAYSTACK_URL`

Optional:

- `HAYSTACK_INDEX` (defaults to `haystack` as in `lore/settings.py`)

# Performance

The primary reason for creating this script is that the Python version was using an excessive amount of RAM and would not run at all on Heroku due to RAM limits.

For a set of test data containing 26,156 learning resources on the same machine, using Docker for the servers:

Python, three runs (MB): 699, 657, 671

Go, three runs (MB): 107, 107, 100

Note that adjusting the `maxRecs` constant will cause the application to use more or less RAM. However, a memory leak in the Python version caused it to use over 600 MB even when adding records one at a time.

#  Binary downloads

- [Mac OSX 64-bit](https://aoeus.com/index_lore_OSX)
- [Linux 64-bit](https://aoeus.com/index_lore_Linux64)

# TODO

These are the next planned steps, in dependency order.

- Create Elasticsearch index and mapping if they are missing (currently this application only refreshes the index; it doesn't rebuild it).
- Add command-line flag to select refresh or rebuild.
- Allow reindexing by any of: all, single repository, single course, single learning resource (currently "all" only).
- Add RESTful API to trigger refresh by any of these, by primary key: all, single repository, single course, single learning resource.
- Add command-line flag to run in API mode or as a one-off.
- Add optional callback feature to API; after indexing is complete, JSON containing number of records indexed will be returned.

Once the API is complete, the LORE application can outsource all indexing to this service, reducing overall load.
