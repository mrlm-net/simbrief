package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Number is a numeric OFP value kept exactly as SimBrief sent it.
//
// SimBrief sends numbers as strings, often zero-padded ("0365", "05000") or
// without a leading zero (".57"), and sometimes empty. Number keeps that raw
// text so nothing is lost, and decodes from a JSON string, a JSON number,
// JSON null or XML character data alike. Use Int, Float or Bool to read it.
type Number string

// UnmarshalJSON accepts a JSON string, number, boolean or null.
func (n *Number) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	switch {
	case len(data) == 0 || bytes.Equal(data, []byte("null")):
		*n = ""
		return nil
	case data[0] == '"':
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*n = Number(s)
		return nil
	case bytes.Equal(data, []byte("true")):
		*n = "1"
		return nil
	case bytes.Equal(data, []byte("false")):
		*n = "0"
		return nil
	}
	var f json.Number
	if err := json.Unmarshal(data, &f); err != nil {
		return fmt.Errorf("simbrief: cannot decode %s as a number", data)
	}
	*n = Number(f)
	return nil
}

// String returns the raw text as sent by SimBrief.
func (n Number) String() string { return string(n) }

// IsSet reports whether the value is not empty.
func (n Number) IsSet() bool { return strings.TrimSpace(string(n)) != "" }

// Float returns the value as float64, or 0 when it is empty or not numeric.
func (n Number) Float() float64 {
	f, err := strconv.ParseFloat(strings.TrimSpace(string(n)), 64)
	if err != nil {
		return 0
	}
	return f
}

// Int returns the value as int (fractions are truncated), or 0 when it is
// empty or not numeric. Leading zeros are fine: "0365" is 365.
func (n Number) Int() int {
	s := strings.TrimSpace(string(n))
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return int(n.Float())
}

// Bool reports whether the value is a non-zero number ("1" is true).
func (n Number) Bool() bool { return n.Float() != 0 }
