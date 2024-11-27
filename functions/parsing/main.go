package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strconv"

	"dave.internal/pkg/intBool"
	"dave.internal/pkg/parser"
)

type convert func(string) (parser.Bindings, error)

func pipedStdin(fn convert) error {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		bindings, err := fn(scanner.Text())
		writeBindings(bindings, err)
	}

	return scanner.Err()
}

var buffer bytes.Buffer

func writeBindings(bindings []parser.Binding, err error) {

	buffer.Reset()

	if err != nil {
		buffer.WriteString("{ \"match\": false, \"error\": ")
		buffer.WriteString(fmt.Sprintf("\"%v\" }", err.Error()))
	} else {
		buffer.WriteString("{ \"match\": true, \"bindings\": [")

		first := true
		for _, b := range bindings {
			if !first {
				buffer.WriteString(", ")
			}
			first = false
			switch v := b.Value.(type) {
			case parser.BindingInt:
				buffer.WriteString("\"")
				buffer.WriteString(b.Name)
				buffer.WriteString("\"=\"")
				buffer.WriteString(strconv.Itoa(int(v)))
				buffer.WriteString("\"")
			case parser.BindingBool:
				buffer.WriteString("\"")
				buffer.WriteString(b.Name)
				buffer.WriteString("\"=\"")
				buffer.WriteString(strconv.FormatBool(bool(v)))
				buffer.WriteString("\"")
			case parser.BindingString:
				buffer.WriteString("\"")
				buffer.WriteString(b.Name)
				buffer.WriteString("\"=\"")
				buffer.WriteString(string(v))
				buffer.WriteString("\"")
			}

		}
		buffer.WriteString("] }")
	}
	fmt.Println(buffer.String())
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

func main() {
	args := os.Args
	fmt.Fprintf(os.Stderr, "Args  %v \n", args)

	switch len(args) {
	case 1:
		const data = `[ foo =42, bar=true, baz = false]`
		p := intBool.IntBoolMappingParser()
		writeBindings(parser.Parse(p.ConfigurationParser, data))

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
