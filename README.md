# jsontodyncol - JSON to Dynamic Column converter

Parses one or more JSON objects from a file and transforms them into MySQL INSERT statements with the JSON translated into dymanic columns.

For info about dynaic columns, read https://mariadb.com/kb/en/mariadb/dynamic-columns/ .
## Building

Build with:
```
go build jsontodyncol.go
```

## Usage
```
JSON to MariaDB Dynamic Column converter 0.1
Usage: jsontodyncol -table TABLE -column COLUMN FILE
  -column string
    	Column name
  -database string
    	Database name
  -help
    	Show this message
  -insert-size int
    	Number of inserted values in each statement (default 1)
  -pretty
    	Pretty-print output
  -table string
    	Table name

```
