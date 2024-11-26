package datetime

import (
	"testing"

	"dave.internal/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestTimeZoneParserFor8601Support(t *testing.T) {
	// Test the timezone parser

	p := timeZone8061Support()

	// Test a valid timezone
	r, err := parser.Parse(p, "+01")
	if err != nil {
		t.Errorf("TimeZone failed: %v", err)
	}
	if r != "+01" {
		t.Errorf("TimeZone failed: %v", r)
	}

	// Test an invalid timezone
	_, err = parser.Parse(p, "+01:00")
	if err != nil {
		t.Errorf("TimeZone failed: %v", err)
	}

	// Test an unrecognized timezone
	_, err = parser.Parse(p, "-0100")
	if err != nil {
		t.Errorf("TimeZone failed: %v", err)
	}

	_, err = parser.Parse(p, "+01:00")
	if err != nil {
		t.Errorf("TimeZone failed: %v", err)
	}

	_, err = parser.Parse(p, "Z")
	if err != nil {
		t.Errorf("TimeZone failed: %v", err)
	}

	//bad
	_, err = parser.Parse(p, "+01:00:00")
	if err == nil {
		t.Errorf("TimeZone failed: %v", err)
	}

	_, err = parser.Parse(p, "Z+01")
	if err == nil {
		t.Errorf("TimeZone failed: %v", err)
	}

	_, err = parser.Parse(p, "Z+01:00")
	if err == nil {
		t.Errorf("TimeZone failed: %v", err)
	}

	_, err = parser.Parse(p, "-01Z")
	if err == nil {
		t.Errorf("TimeZone failed: %v", err)
	}

}

func TestTimePartsFor8601Support(t *testing.T) {

	p := time8601Support()

	// Test a valid time
	r, err := parser.Parse(p, "01:01:01")
	if err != nil {
		t.Errorf("TimeParts failed: %v", err)
	}
	if r != "01:01:01" {
		t.Errorf("TimeParts failed: %v", r)
	}

	_, err = parser.Parse(p, "01:01:01.000")
	if err != nil {
		t.Errorf("TimeParts failed: %v", err)
	}

	_, err = parser.Parse(p, "010101")
	if err != nil {
		t.Errorf("TimeParts failed: %v", err)
	}

	_, err = parser.Parse(p, "010101.000")
	if err != nil {
		t.Errorf("TimeParts failed: %v", err)
	}

	// Test an invalid time
	_, err = parser.Parse(p, "01:01:01-01")
	if err == nil {
		t.Errorf("TimeParts failed: %v", err)
	}

	// Test an unrecognized time
	_, err = parser.Parse(p, "01:01:01.000+01")
	if err == nil {
		t.Errorf("TimeParts failed: %v", err)
	}
}

func TestDatePartsFor8601Support(t *testing.T) {
	// Test the date parts parser

	p := date8601Support()

	// Test a valid date
	r, err := parser.Parse(p, "2020-01-01")
	if err != nil {
		t.Errorf("DateParts failed: %v", err)
	}
	if r != "2020-01-01" {
		t.Errorf("DateParts failed: %v", r)
	}

	r, err = parser.Parse(p, "20200101")
	if err != nil {
		t.Errorf("DateParts failed: %v", err)
	}
	assert.Equal(t, "20200101", r)

	// Test an invalid date
	_, err = parser.Parse(p, "2020-01-01-01")
	if err == nil {
		t.Errorf("DateParts failed: %v", err)
	}

	// Test an odd date
	_, err = parser.Parse(p, "2020-011")
	if err == nil {
		t.Errorf("DateParts failed: %v", err)
	}

	// Test an unrecognized date
	_, err = parser.Parse(p, "+Z")
	if err == nil {
		t.Errorf("DateParts failed: %v", err)
	}

}

func TestISO8601Parser(t *testing.T) {

	p := ISO8601Parser()

	// Test a valid date
	r, err := parser.Parse(p, "2020-01-01T01:01:01+01")
	if err != nil {
		t.Errorf("ISO8601Parser failed: %v", err)
	}
	if r != "2020-01-01 01:01:01+01" {
		t.Errorf("ISO8601Parser failed: %v", r)
	}

	// Test a valid date
	r, err = parser.Parse(p, "2020-01-01 01:01:01Z")
	if err != nil {
		t.Errorf("ISO8601Parser failed: %v", err)
	}
	if r != "2020-01-01 01:01:01Z" {
		t.Errorf("ISO8601Parser failed: %v", r)
	}

	// Test an invalid date
	_, err = parser.Parse(p, "2020-01-01T01:01:01+010")
	if err == nil {
		t.Errorf("ISO8601Parser failed: %v", err)
	}

}
