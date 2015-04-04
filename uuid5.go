package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/satori/go.uuid"
)

var gnNamespace uuid.UUID
var url string

func namespace() uuid.UUID {
	if gnNamespace.String() == "00000000-0000-0000-0000-000000000000" {
		gnNamespace = uuid.NewV5(uuid.NamespaceDNS, "globalnames.org")
	}
	return gnNamespace
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
	nameStrings := splitNames(input)
	var buffer bytes.Buffer
	buffer.WriteString("[\n")
	length := len(nameStrings) - 1
	for i := range nameStrings {
		gnUUID := uuid.NewV5(namespace(), nameStrings[i])
		buffer.WriteString("  {\n")
		buffer.WriteString("    \"name_string\": \"")
		buffer.WriteString(nameStrings[i])
		buffer.WriteString("\",\n")
		buffer.WriteString("    \"uuid\": \"")
		buffer.WriteString(gnUUID.String())
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
	fmt.Fprintf(w, "Enter name string in url like \"%s/Homo%%20sapiens\"\n", url)
	fmt.Fprintf(w, "Or enter serveral name strings divided by pipe character like \"%s/Homo%%20sapiens%%7CPardosa%%20moesta%%7CParus%%20major%%20(Linnaeus,%%201758)\"\n", url)
}

func splitNames(input string) []string {
	f := func(c rune) bool {
		pipe := '|'
		return c == pipe
	}
	return strings.FieldsFunc(input, f)
}

func main() {
	portPtr := flag.String("port", "8080", "Port to run the server")
	urlPtr := flag.String("url", "http://localhost", "URL to show in help")
	flag.Parse()
	port := ":" + *portPtr
	url = *urlPtr
	fmt.Printf("Server started at %s%s\n\n", url, port)
	fmt.Print("To change url, port use flags:\n\n")
	fmt.Printf("  %s --port 80 --url http://your_url", os.Args[0])
	http.HandleFunc("/", handler)
	http.ListenAndServe(port, nil)
}
