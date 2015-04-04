package main

import (
	"bytes"
	"fmt"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"strings"
)

var gn_namespace uuid.UUID

func namespace() uuid.UUID {
	if gn_namespace.String() == "00000000-0000-0000-0000-000000000000" {
		gn_namespace = uuid.NewV5(uuid.NamespaceDNS, "globalnames.org")
	}
	return gn_namespace
}

func handler(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Path[1:]
	if input == "" {
		instructions(w)
	} else {
		uuids(w, input)
	}
}

func uuids(w http.ResponseWriter, input string) {
	name_strings := split_names(input)
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	length := len(name_strings) - 1
	for i := range name_strings {
		gn_uuid := uuid.NewV5(namespace(), name_strings[i])
		buffer.WriteString("  {\n")
		buffer.WriteString("    \"name_string\": \"")
		buffer.WriteString(name_strings[i])
		buffer.WriteString("\",\n")
		buffer.WriteString("    \"uuid\": \"")
		buffer.WriteString(gn_uuid.String())
		if length == i {
			buffer.WriteString("\"\n  }\n")
		} else {
			buffer.WriteString("\"\n  },\n")
		}
	}
	buffer.WriteString("]")
	fmt.Fprintf(w, buffer.String())
}

func instructions(w http.ResponseWriter) {
	url := os.Args[1]
	fmt.Fprintf(w, "Enter name string in url like \"%s/Homo%%20sapiens\"\n", url)
	fmt.Fprintf(w, "Or enter serversl name strings divided by pipe character like \"%s/Homo%%20sapiens%%7CPardosa%%20moesta%%7CParus%%20major%%20(Linnaeus,%%201758)\"\n", url)
}

func split_names(input string) []string {
	f := func(c rune) bool {
		var pipe rune = '|'
		return c == pipe
	}
	return strings.FieldsFunc(input, f)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
