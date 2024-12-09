package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"dave.internal/pkg/intBool"
	"dave.internal/pkg/parser"
	"dave.internal/pkg/syslog"

	"github.com/redpanda-data/redpanda/src/transform-sdk/go/transform"
)

type convert func(string) (parser.Bindings, error)
type process func(string) (string, error)

var buffer bytes.Buffer

func convertToProcess(fn convert) process {
	return func(data string) (string, error) {
		bindings, err := fn(data)
		buffer.Reset()
		parser.WriteBindingsAsJson(&buffer, data, bindings, err)
		return buffer.String(), err
	}
}

func pipedStdin(fn process) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		result, err := fn(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Printf("%s\n", result)
	}
	return scanner.Err()
}

func convIntParser() convert {
	p := intBool.IntBoolMappingParser()
	return func(data string) (parser.Bindings, error) {
		return parser.Parse(p.ConfigurationParser, parser.WithState(data))
	}
}

func syslogParserRaw() process {
	p := syslog.SyslogParserRaw()

	return func(data string) (string, error) {
		b, err := parser.Parse(p, parser.WithState(data))
		if err != nil {
			return "", err
		}
		return b.CompactJson(), nil
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
		case "syslogRaw":
			pipedStdin(syslogParserRaw())
		case "intBool":
			pipedStdin(convertToProcess(convIntParser()))
		default:
			fmt.Fprintf(os.Stderr, "Invalid argument %v\n", args[1])
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Invalid argument - %v\n", args)
	}
}
