// Copyright (c) 2015, Markus Mäkelä
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of jsontodyncol nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package main

import (
    "fmt"
    "encoding/json"
    "os"
    "strconv"
    "strings"
    "flag"
	"path/filepath"
	"io"
)

// Command line flags
var insert_size = flag.Int("insert-size", 1, "Number of inserted values in each statement")
var help = flag.Bool("help", false, "Show this message")
var pretty = flag.Bool("pretty", false, "Pretty-print output")
var database = flag.String("database", "", "Database name")
var table = flag.String("table", "", "Table name")
var column = flag.String("column", "", "Column name")

// Print usage
func Usage(){
	fmt.Println("JSON to MariaDB Dynamic Column converter 0.1")
	fmt.Println()
	fmt.Println("Usage:", filepath.Base(os.Args[0]), "-table TABLE -column COLUMN [OPTIONS] [FILE]")
	fmt.Println()
	fmt.Println("Reads one or more JSON objects from a file and convert them into MariaDB compatible")
	fmt.Println("INSERT statements. If no input file is provided the standard input is read.")
	fmt.Println()
	fmt.Println("Options:")
	flag.PrintDefaults()
}

// Main function.
// Parses first argument as a file with one or more JSON objects
// and converts them into dynamic column insert statements.
func main(){
    flag.Parse()

    switch {
	case *help:
		Usage()
        os.Exit(0)
	case len(*table) == 0:
        fmt.Fprintln(os.Stderr, "No table name provided! See -help output for more info.")
        os.Exit(1)
	case len(*column) == 0:
		fmt.Fprintln(os.Stderr, "No column provided! See -help output for more info.")
        os.Exit(1)
	default:
	}

	// If no inputs provided, use Stdin
	var input *os.File
	if len(flag.Args()) < 1 {
		input = os.Stdin
	} else {
		file, err := os.Open(flag.Args()[0])
		if err != nil{
			fmt.Fprintln(os.Stderr, "Fatal error:", err)
			os.Exit(1)
		}
		input = file
	}

    decoder := json.NewDecoder(input)
	var err_d error = nil
	values := 0
    for err_d == nil{
        str := "INSERT INTO "
		if len(*database) > 0 {
			str += *database + "."
		}
		str += *table + "(" + *column + ") values"

		n_inserts := *insert_size
		for err_d == nil{
			var obj map[string]interface{}
			if err_d = decoder.Decode(&obj); err_d != nil {
				str = strings.TrimRight(str, ",\n")
				str += ";"

				if err_d != io.EOF {
					fmt.Fprintln(os.Stderr, err_d)
				}

			} else {
				str += " ("
				str += PrintObject(&obj)
				n_inserts--
				values++
				if n_inserts <= 0{
					// Last value
					str += ");"
					break
				}

				// More values to insert
				str += "),"
				if *pretty {
					str += "\n"
				}
			}
		}

		// At least one value was read
		if n_inserts < *insert_size {
			fmt.Println(str)
		}
    }

	if values == 0 {
		fmt.Fprintln(os.Stderr, "No valid values read.")
	}

}

// Print a JSON array as a comma-separated list strings
func PrintList(mylist []interface{}) string{
    str := ""
    for _, elem := range mylist {
		if len(str) > 0 {
			str += ","
		}
        switch v := elem.(type) {
        case string:
            str += fmt.Sprintf("%q", v)

        case bool:
            str += strconv.FormatBool(v)

        case float64:
            str += strconv.FormatFloat(v, 'f', -1, 64)

        default:
			fmt.Fprintln(os.Stderr, "Unknown type:", v)
			// Unknown value
        }
    }
    return str
}

// String cleanup function
func FormatStr(str string) string {
    replacer := strings.NewReplacer("'", "\\'", "\"", "\\\"")
    return strconv.QuoteToASCII(replacer.Replace(str))
}

// Print a JSON object and format it as a COLUMN_CREATE statement
func PrintObject(obj* map[string]interface{}) string{
    str := "COLUMN_CREATE("
    for key, value := range *obj {

        switch conv := value.(type) {
        case []interface{}:
            str += FormatStr(key) + "," + FormatStr(PrintList(conv))

        case map[string]interface{}:
            str += FormatStr(key) + "," +  PrintObject(&conv)

        case string:
            str += FormatStr(key) + "," + FormatStr(conv)

        case float64:
            str += FormatStr(key) + "," + FormatStr(strconv.FormatFloat(conv, 'f', -1, 64))
        }
        str += ","
    }
    str = str[0:len(str) - 1] + ")"
    return str
}
