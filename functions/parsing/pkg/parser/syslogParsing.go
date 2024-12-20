package parser

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

// RFC 3164 -
//
//	<PRIVAL>TIMESTAMP HOSTNAME TAG: MESSAGE
//	<13>Oct 22 12:34:56 myhostname myapp[1234]: This is a sample syslog message.
//
// RFC 5424 -
//
//	<PRIVAL>VERSION TIMESTAMP HOSTNAME APP-NAME PROCID MSGID [STRUCTURED-DATA] MESSAGE
// <165>1 2003-10-11T22:14:15.003Z myhostname myapp 1234 ID47 - [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] An application event log entry...

type Priority struct {
	Facility int
	Severity int
}

type message struct {
	Raw      string
	Bindings []Binding
	Columns  []string
}

type tag struct {
	AppName string
	Pid     int
}

type structureData struct {
	Bindings BindingList
}

type SyslogMetadata struct {
	Format     string
	Priority   Priority
	Timestamp  string
	Hostname   string
	Tag        tag
	Structured structureData
	Message    message
}

type SyslogMetadataRaw struct {
	Format  string
	Message message
}

// ////////////////////////////////////
// JSON
func encode(in string) string {
	byt, err := json.Marshal(in)
	if err != nil {
		return ""
	}
	return string(byt)
}

func (p Priority) CompactJson() string {
	return fmt.Sprintf("{\"fac\":%d,\"sev\":%d}", p.Facility, p.Severity)
}

func CompactJsonBindings(b BindingList) string {
	var buffer string
	for _, binding := range b {
		if len(buffer) > 0 {
			buffer += ","
		}
		buffer += CompactJsonBinding(binding)
	}
	return fmt.Sprintf("{%s}", buffer)
}

func CompactJsonBinding(b Binding) string {
	switch b.Value.(type) {
	case BindingInt:
		return fmt.Sprintf("\"%s\":%d", b.Name, b.Value)
	case BindingBool:
		return fmt.Sprintf("\"%s\":%t", b.Name, b.Value)
	case BindingString:
		return fmt.Sprintf("\"%s\":%s", b.Name, encode(string(b.Value.(BindingString))))
	case BindingBinding:
		bindings := b.Value.(BindingBinding)

		var buffer string
		for _, binding := range bindings {
			if len(buffer) > 0 {
				buffer += ","
			}
			buffer += CompactJsonBinding(binding)
		}
		return fmt.Sprintf("\"%s\":{%s}", b.Name, buffer)
	}
	return ""
}

func CompactJsonColumns(c []string) string {
	var buffer string
	for _, col := range c {
		if len(buffer) > 0 {
			buffer += ","
		}
		buffer += encode(col)
	}
	return fmt.Sprintf("[%s]", buffer)
}

func (m message) CompactJson() string {
	return fmt.Sprintf("{\"raw\":%s,\"bnd\":%v,\"col\":%v}", encode(m.Raw), CompactJsonBindings(m.Bindings), CompactJsonColumns(m.Columns))
}

func (t tag) CompactJson() string {
	return fmt.Sprintf("{\"app\":%s,\"pid\":%d}", encode(t.AppName), t.Pid)
}

func (s structureData) CompactJson() string {
	return fmt.Sprintf("{\"bnd\":%v}", CompactJsonBindings(s.Bindings))
}

func (s SyslogMetadataRaw) CompactJson() string {
	return fmt.Sprintf("{\"fmt\":\"%s\",\"msg\":%s}", s.Format, s.Message.CompactJson())
}

// ////////////////////////////////////
func PriorityParser() Parser[Priority] {
	w1 := StartSkipping(Exactly("<"))
	k1 := AppendKeeping(w1, IntParser)
	w2 := AppendSkipping(k1, Exactly(">"))
	return Apply(w2, func(p int) Priority {
		var val int = p / 10
		return Priority{Facility: val, Severity: p % 10}
	})
}

// ////////////////////////////////////
// Major Assemblies
// ////////////////////////////////////

func Tag3164Parser() Parser[tag] {
	w1 := StartKeeping(EntityNameParser)
	k1 := AppendSkipping(w1, Exactly("["))
	w2 := AppendKeeping(k1, IntParser)
	k2 := AppendSkipping(w2, Exactly("]"))

	return Apply2(k2, func(appName string, pid int) tag {
		return tag{AppName: appName, Pid: pid}
	})
}

// func structureDataParser() Parser[structureData] {
// 	w1 := StartSkipping(Exactly("["))
// 	k1 := AppendKeeping(w1, NameParser())
// 	w2 := AppendSkipping(k1, Exactly("@"))
// 	k2 := AppendKeeping(w2, IntParser)
// 	w3 := AppendSkipping(k2, Exactly(" "))
// 	k3 := AppendKeeping(w3, BindingsParser())
// 	w4 := AppendSkipping(k3, Exactly("]"))

// 	return Apply3(w4, func(name string, id int, bindings Bindings) structureData {
// 		sdName := Binding{Name: "sd-name", Value: BindingString(name)}
// 		sdId := Binding{Name: "sd-id", Value: BindingInt(id)}

// 		return structureData{Bindings: bindings + Bindings{sdName, sdId}}
// 	})

// }

// ////////////////////////////////////

func SyslogParserRaw() Parser[SyslogMetadataRaw] {

	p := GetString(ConsumeWhile(func(r rune) bool {
		return r != utf8.RuneError
	}))

	return Map(p, func(s string) SyslogMetadataRaw {
		return SyslogMetadataRaw{
			Format:  "raw",
			Message: message{Raw: s},
		}
	})
}
