package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityParser(t *testing.T) {
	p, err := Parse(PriorityParser(), WithState("<13>"))
	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Equal(t, Priority{Facility: 1, Severity: 3}, p)
	assert.Equal(t, "{\"fac\":1,\"sev\":3}", p.CompactJson())
	fmt.Println(p.CompactJson())

	p, err = Parse(PriorityParser(), WithState("<165>"))
	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Equal(t, Priority{Facility: 16, Severity: 5}, p)
	assert.Equal(t, "{\"fac\":16,\"sev\":5}", p.CompactJson())
	fmt.Println(p.CompactJson())
}

func TestTag3164Parser(t *testing.T) {
	p, err := Parse(Tag3164Parser(), WithState("myapp[1234]"))
	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Equal(t, tag{AppName: "myapp", Pid: 1234}, p)
	assert.Equal(t, "{\"app\":\"myapp\",\"pid\":1234}", p.CompactJson())
	fmt.Println(p.CompactJson())
}

func TestCompactJsonBinding(t *testing.T) {
	b := Binding{Name: "test", Value: BindingInt(42)}
	assert.Equal(t, "\"test\":42", CompactJsonBinding(b))
	fmt.Println(CompactJsonBinding(b))

	b = Binding{Name: "test", Value: BindingBool(true)}
	assert.Equal(t, "\"test\":true", CompactJsonBinding(b))
	fmt.Println(CompactJsonBinding(b))

	b = Binding{Name: "test", Value: BindingString("foo")}
	assert.Equal(t, "\"test\":\"foo\"", CompactJsonBinding(b))
	fmt.Println(CompactJsonBinding(b))

	b = Binding{Name: "test", Value: BindingBinding{Binding{Name: "foo", Value: BindingInt(42)}}}
	assert.Equal(t, "\"test\":{\"foo\":42}", CompactJsonBinding(b))
	fmt.Println(CompactJsonBinding(b))

	b = Binding{Name: "test", Value: BindingBinding{Binding{Name: "foo", Value: BindingInt(42)}, Binding{Name: "bar", Value: BindingString("baz")}}}
	assert.Equal(t, "\"test\":{\"foo\":42,\"bar\":\"baz\"}", CompactJsonBinding(b))
	fmt.Println(CompactJsonBinding(b))
}

func TestCompactJsonBindings(t *testing.T) {
	b := BindingList{
		Binding{Name: "test", Value: BindingInt(42)},
		Binding{Name: "test2", Value: BindingBool(true)},
		Binding{Name: "test3", Value: BindingString("foo")},
		Binding{Name: "test4", Value: BindingBinding{Binding{Name: "foo", Value: BindingInt(42)}}},
	}
	assert.Equal(t, "{\"test\":42,\"test2\":true,\"test3\":\"foo\",\"test4\":{\"foo\":42}}", CompactJsonBindings(b))
	fmt.Println(CompactJsonBindings(b))
}

func TestCompactJsonColumns(t *testing.T) {
	c := []string{"col1", "col2", "col3"}
	assert.Equal(t, "[\"col1\",\"col2\",\"col3\"]", CompactJsonColumns(c))
	fmt.Println(CompactJsonColumns(c))
}

func TestMessageJson(t *testing.T) {
	m := message{
		Raw: "<13>Oct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.",
		Bindings: BindingList{
			Binding{Name: "test", Value: BindingInt(42)},
			Binding{Name: "test2", Value: BindingBool(true)},
			Binding{Name: "test3", Value: BindingString("foo")},
			Binding{Name: "test4", Value: BindingBinding{Binding{Name: "foo", Value: BindingInt(42)}}},
		},
		Columns: []string{"col1", "col2", "col3"},
	}
	assert.Equal(t, "{\"raw\":\"\\u003c13\\u003eOct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.\",\"bnd\":{\"test\":42,\"test2\":true,\"test3\":\"foo\",\"test4\":{\"foo\":42}},\"col\":[\"col1\",\"col2\",\"col3\"]}", m.CompactJson())
	fmt.Println(m.CompactJson())
}

func TestSyslogParserRaw(t *testing.T) {
	p, err := Parse(SyslogParserRaw(), WithState("<13>Oct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message."))
	assert.NotNil(t, p)
	assert.Nil(t, err)
	assert.Equal(t, SyslogMetadataRaw{
		Format: "raw",
		Message: message{
			Raw:      "<13>Oct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.",
			Bindings: nil,
			Columns:  nil,
		},
	}, p)
	assert.Equal(t, "{\"fmt\":\"raw\",\"msg\":{\"raw\":\"\\u003c13\\u003eOct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.\",\"bnd\":{},\"col\":[]}}", p.CompactJson())
	fmt.Println(p.CompactJson())
}

// func TestSyslogSnipeStructuredData(t *testing.T) {

// 	initial := WithState("<13>Oct 22 12:34:56 myhostname mypp[1234]: - [exampleSDID@32473 iut=\"3\" eventSource=\"Application\" eventID=\"1011\"] This is a sample syslog message. [test1=\"foo\" test2=\"bar\"]")
// 	assert.Equal(t, 172, len(initial.remaining()))

// 	snipes := initial.Snipe(func(r rune) bool {
// 		return r == '['
// 	})

// 	parser := structureDataParser()

// 	for _, s := range snipes {
// 		fmt.Println(s.remaining())
// 		p, err := Parse(parser, s)
// 		if err == nil {
// 			fmt.Println(p.CompactJson())
// 		} else {
// 			fmt.Printf("Err %v\n", err)
// 		}
// 		fmt.Println()
// 	}

// 	// 	Bindings := Bindings{
// 	// 		Binding{Name: "exampleSDID", Value: BindingString("exampleSDID@32473")},
// 	// 		Binding{Name: "iut", Value: BindingString("3")},
// 	// 		Binding{Name: "eventSource", Value: BindingString("Application")},
// 	// 		Binding{Name: "eventID", Value: BindingString("1011")},
// 	// 		Binding{Name: "test1", Value: BindingString("foo")},
// 	// 		Binding{Name: "test2", Value: BindingString("bar")},
// 	// 	}

// 	// assert.Equal(t, 6, len(Bindings))
// }
