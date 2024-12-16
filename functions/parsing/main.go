package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"dave.internal/pkg/intBool"
	"dave.internal/pkg/parser"
	"dave.internal/pkg/syslog"

	"github.com/redpanda-data/redpanda/src/transform-sdk/go/transform"
)

type convert func(string) (parser.Bindings, error)
type process func(string) ([]string, error)

var buffer bytes.Buffer

func convertToProcess(fn convert) process {
	return func(data string) ([]string, error) {
		bindings, err := fn(data)
		buffer.Reset()
		parser.WriteBindingsAsJson(&buffer, data, bindings, err)
		str := buffer.String()
		return []string{str}, err
	}
}

func pipedStdin(fn process) error {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		result, err := fn(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		} else {
			for _, line := range result {
				fmt.Printf("%s\n", line)
			}
		}
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

	return func(data string) ([]string, error) {
		b, err := parser.Parse(p, parser.WithState(data))
		if err != nil {
			return []string{}, err
		}
		return []string{b.CompactJson()}, nil
	}

}

// doTransform is where you read the record that was written, and then you can
// output new records that will be written to the destination topic
func doTransform(e transform.WriteEvent, w transform.RecordWriter) error {
	return w.Write(e.Record())
}

func encode(in string) string {
	byt, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(byt)
}

func multilineParser() process {

	w1 := parser.StartSkipping(parser.Exactly("<"))
	k1 := parser.AppendKeeping(w1, parser.IntParser)
	w2 := parser.AppendSkipping(k1, parser.Exactly(">"))
	s2 := parser.AppendKeeping(w2, parser.Map(parser.OneOf(
		parser.NumberStringParser,
		parser.MonthAsciiParser,
	), func(text string) string { return text }))

	predicate := parser.Apply2(s2, func(p int, month string) bool {
		return true
	})

	p := parser.MultilineParser('<', predicate)

	return func(data string) ([]string, error) {
		b, _, err := p(parser.WithState(data))
		if err != nil {
			return []string{}, err
		}
		return b, nil
	}
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
		case "multiline":
			pipedStdin(multilineParser())
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
