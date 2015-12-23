# jsontodyncol - JSON to Dynamic Column converter

Parses one or more JSON objects from a file and transforms them into MySQL INSERT statements with the JSON translated into dynamic columns.

For information about dynamic columns, read https://mariadb.com/kb/en/mariadb/dynamic-columns/ .
## Building

Build with:
```
go build jsontodyncol.go
```

## Usage
```
JSON to MariaDB Dynamic Column converter 0.1

Usage: jsontodyncol -table TABLE -column COLUMN [OPTIONS] [FILE]

Reads one or more JSON objects from a file and convert them into MariaDB compatible
INSERT statements. If no input file is provided the standard input is read.

Options:
  -column="": Column name
  -database="": Database name
  -help=false: Show this message
  -insert-size=1: Number of inserted values in each statement
  -pretty=false: Pretty-print output
  -table="": Table name
```
