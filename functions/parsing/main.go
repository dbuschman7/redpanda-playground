package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"

	"dave.internal/pkg/intBool"
	"dave.internal/pkg/parser"

	"github.com/redpanda-data/redpanda/src/transform-sdk/go/transform"
)

type convert func(string) (parser.Bindings, error)

var buffer bytes.Buffer

func pipedStdin(fn convert) error {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		var line = scanner.Text()
		bindings, err := fn(line)
		buffer.Reset()
		parser.WriteBindingsAsJson(&buffer, line, bindings, err)
		fmt.Printf("%s\n", buffer.String())
	}
	return scanner.Err()
}

func convIntParser() convert {
	p := intBool.IntBoolMappingParser()
	return func(data string) (parser.Bindings, error) {
		return parser.Parse(p.ConfigurationParser, data)
	}
}

func syslogParser() convert {
	return func(data string) (parser.Bindings, error) {
		return nil, errors.New("not implemented")
	}
}

// doTransform is where you read the record that was written, and then you can
// output new records that will be written to the destination topic
func doTransform(e transform.WriteEvent, w transform.RecordWriter) error {
	return w.Write(e.Record())
}

func main() {
	args := os.Args
	fmt.Fprintf(os.Stderr, "Args(%v) %v \n", len(args), args)

	switch len(args) {
	case 1:
		// Register your transform function.
		// This is a good place to perform other setup too.
		transform.OnRecordWritten(doTransform)
	case 2:
		switch args[1] {
		case "syslog":
			pipedStdin(syslogParser())
		case "intBool":
			pipedStdin(convIntParser())
		default:
			fmt.Fprintf(os.Stderr, "Invalid argument %v\n", args[1])
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Invalid argument - %v\n", args)
	}
}
