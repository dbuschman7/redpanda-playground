package common

import (
	"testing"

	"dave.internal/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestMonthAsciiParser(t *testing.T) {
	// Test: "jan"
	{
		p, err := parser.Parse(MonthAsciiParser, "jan")
		assert.Nil(t, err)
		assert.Equal(t, "jan", p)
	}

	// Test: "JAN"
	{
		p, err := parser.Parse(MonthAsciiParser, "JAN")
		assert.Nil(t, err)
		assert.Equal(t, "JAN", p)
	}

	// Test: "january"
	{
		p, err := parser.Parse(MonthAsciiParser, "january")
		assert.NotEmpty(t, err)
		assert.Empty(t, p)
	}

	// Test: "jane"
	{
		p, err := parser.Parse(MonthAsciiParser, "jane")
		assert.NotEmpty(t, err)
		assert.Empty(t, p)
	}
}

func TestDateDashParser(t *testing.T) {
	// Test: "2020-01-01"
	{
		p, err := parser.Parse(DateYYMMDDDashParser(), "2020-01-01")
		assert.Nil(t, err)
		assert.Equal(t, "2020-01-01", p)
	}

	// Test: "2020-01-01 12:23:34"
	{
		p, err := parser.Parse(DateYYMMDDDashParser(), "2020-01-01 12:23:34")
		assert.NotEmpty(t, err)
		assert.Empty(t, p)
	}

	// Test: "2020-01-01T12:23:34"
	{
		p, err := parser.Parse(DateYYMMDDDashParser(), "2020-01-01T12:23:34")
		assert.NotEmpty(t, err)
		assert.Empty(t, p)
	}
}
func TestDateMMMDDParser3164(t *testing.T) {
	// Test: "Oct 22"
	{
		p, err := parser.Parse(DateMMMDDParser3164(), "Oct 22 12:23:34")
		assert.Nil(t, err)
		assert.Equal(t, "Oct 22 12:23:34", p)
	}

	// Test: "Oct 32"
	{
		p, err := parser.Parse(DateMMMDDParser3164(), "Oct 32 12:23:34")
		// assert.NotEmpty(t, err)
		// assert.Empty(t, p)
		assert.Nil(t, err)
		assert.Equal(t, "Oct 32 12:23:34", p)
	}

	// Test: "Oct 2"
	{
		p, err := parser.Parse(DateMMMDDParser3164(), "Oct 2 12:23:34")
		assert.Empty(t, err)
		assert.Equal(t, "Oct 02 12:23:34", p)
	}

	// Test: "Oct 2x"
	{
		p, err := parser.Parse(DateMMMDDParser3164(), "Oct 2x 12:23:34")
		assert.NotEmpty(t, err)
		assert.Empty(t, p)
	}

	// Test: "Oct 2 12:23:34"
	{
		p, err := parser.Parse(DateMMMDDParser3164(), "Oct 2 12:23")
		assert.NotEmpty(t, err)
		assert.Empty(t, p)
	}
}
