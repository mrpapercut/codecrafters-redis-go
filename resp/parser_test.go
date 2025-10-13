package resp

import (
	"math"
	"math/big"
	"testing"
)

var parser = GetParser()

func TestParseSimpleString(t *testing.T) {
	raw := "+OK\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != SimpleString {
		t.Errorf("expected Type to be SimpleString, got '%s' instead", getType(parsed.Type))
	}

	if parsed.String != "OK" {
		t.Errorf("expected String to be 'OK', got '%s' instead", parsed.String)
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseSimpleError(t *testing.T) {
	raw := "-ERR message\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != SimpleError {
		t.Errorf("expected Type to be SimpleError, got '%s' instead", getType(parsed.Type))
	}

	if parsed.ErrorType != "ERR" {
		t.Errorf("expected ErrorType to be 'ERR', got '%s' instead", parsed.ErrorType)
	}

	if parsed.ErrorMessage != "message" {
		t.Errorf("expected ErrorMessage to be 'message', got '%s' instead", parsed.ErrorMessage)
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseInteger(t *testing.T) {
	raw := ":1000\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != Integer {
		t.Errorf("expected Type to be Integer, got '%s' instead", getType(parsed.Type))
	}

	if parsed.Integer != 1000 {
		t.Errorf("expected Integer to be 1000, got %d instead", parsed.Integer)
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}

	raw = ":-1000\r\n"
	parsed = parser.ParseCommand(raw)

	if parsed.Type != Integer {
		t.Errorf("expected Type to be Integer, got '%s' instead", getType(parsed.Type))
	}

	if parsed.Integer != -1000 {
		t.Errorf("expected Integer to be -1000, got %d instead", parsed.Integer)
	}

	respString = parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseBulkString(t *testing.T) {
	raw := "$5\r\nhello\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != BulkString {
		t.Errorf("expected Type to be BulkString, got '%s' instead", getType(parsed.Type))
	}

	if parsed.String != "hello" {
		t.Errorf("expected String to be 'hello', got '%s' instead", parsed.String)
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}

	// Empty string
	raw = "$0\r\n\r\n"
	parsed = parser.ParseCommand(raw)

	if parsed.Type != BulkString {
		t.Errorf("expected Type to be BulkString, got '%s' instead", getType(parsed.Type))
	}

	if parsed.String != "" {
		t.Errorf("expected String to be '', got '%s' instead", parsed.String)
	}

	respString = parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}

	// Null-string
	raw = "$-1\r\n"
	parsed = parser.ParseCommand(raw)

	if parsed.Type != BulkString {
		t.Errorf("expected Type to be BulkString, got '%s' instead", getType(parsed.Type))
	}

	if !parsed.IsNull {
		t.Error("expected IsNull to be true, got false instead")
	}

	respString = parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseArray(t *testing.T) {
	// Empty array
	raw := "*0\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != Array {
		t.Errorf("expected Type to be Array, got '%s' instead", getType(parsed.Type))
	}

	if len(parsed.Array) != 0 {
		t.Errorf("expected Array to have length 0, got %d instead", len(parsed.Array))
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}

	// Null-array
	raw = "*-1\r\n"
	parsed = parser.ParseCommand(raw)

	if parsed.Type != Array {
		t.Errorf("expected Type to be Array, got '%s' instead", getType(parsed.Type))
	}

	if !parsed.IsNull {
		t.Error("expected IsNull to be true, got false instead")
	}

	respString = parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}

	// Array of strings
	raw = "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	parsed = parser.ParseCommand(raw)

	if parsed.Type != Array {
		t.Errorf("expected Type to be Array, got '%s' instead", getType(parsed.Type))
	}

	if len(parsed.Array) != 2 {
		t.Errorf("expected Array to have length 2, got %d instead", len(parsed.Array))
	}

	if parsed.Array[0].Type != BulkString {
		t.Errorf("expected Array[0].Type to be BulkString, got '%s' instead", getType(parsed.Array[0].Type))
	}

	if parsed.Array[0].String != "hello" {
		t.Errorf("expected Array[0].String to be 'hello', got '%s' instead", parsed.Array[0].String)
	}

	if parsed.Array[1].Type != BulkString {
		t.Errorf("expected Array[1].Type to be BulkString, got '%s' instead", getType(parsed.Array[1].Type))
	}

	if parsed.Array[1].String != "world" {
		t.Errorf("expected Array[1].String to be 'world', got '%s' instead", parsed.Array[1].String)
	}

	respString = parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}

	// Array with mixed data types
	raw = "*2\r\n:1\r\n$5\r\nhello\r\n"
	parsed = parser.ParseCommand(raw)

	if parsed.Type != Array {
		t.Errorf("expected Type to be Array, got '%s' instead", getType(parsed.Type))
	}

	if len(parsed.Array) != 2 {
		t.Errorf("expected Array to have length 2, got %d instead", len(parsed.Array))
	}

	if parsed.Array[0].Type != Integer {
		t.Errorf("expected Array[0].Type to be Integer, got '%s' instead", getType(parsed.Array[0].Type))
	}

	if parsed.Array[0].Integer != 1 {
		t.Errorf("expected Array[0].Integer to be '1', got '%d' instead", parsed.Array[0].Integer)
	}

	if parsed.Array[1].Type != BulkString {
		t.Errorf("expected Array[1].Type to be BulkString, got '%s' instead", getType(parsed.Array[1].Type))
	}

	if parsed.Array[1].String != "hello" {
		t.Errorf("expected Array[1].String to be 'hello', got '%s' instead", parsed.Array[1].String)
	}

	respString = parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseNestedArrays(t *testing.T) {
	raw := "*2\r\n$5\r\nouter\r\n*2\r\n$7\r\ninner-1\r\n$7\r\ninner-2\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != Array {
		t.Errorf("expected Type to be Array, got %s instead", getType(parsed.Type))
	}

	if len(parsed.Array) != 2 {
		t.Errorf("expected Array to have length 2, got %d instead", len(parsed.Array))
	}

	if parsed.Array[0].Type != BulkString {
		t.Errorf("expected Array[0].Type to be BulkString, got %s instead", getType(parsed.Array[0].Type))
	}

	if parsed.Array[0].String != "outer" {
		t.Errorf("expected Array[0].String to be 'PING', got '%s' instead", parsed.Array[0].String)
	}

	if parsed.Array[1].Type != Array {
		t.Errorf("expected Array[1].Type to be Array, got %s instead", getType(parsed.Array[1].Type))
	}

	if len(parsed.Array[1].Array) != 2 {
		t.Errorf("expected Array[1].Array to have length 2, got %d instead", len(parsed.Array[1].Array))
	}

	if parsed.Array[1].Array[0].String != "inner-1" {
		t.Errorf("expected Array[1].Array[0].String to be 'inner-1', got '%s' instead", parsed.Array[1].Array[0].String)
	}

	if parsed.Array[1].Array[1].String != "inner-2" {
		t.Errorf("expected Array[1].Array[1].String to be 'inner-2', got '%s' instead", parsed.Array[1].Array[1].String)
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseNull(t *testing.T) {
	raw := "_\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != Null {
		t.Errorf("expected Type to be Null, got '%s' instead", getType(parsed.Type))
	}

	if !parsed.IsNull {
		t.Error("expected IsNull to be true, got false instead")
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseBoolean(t *testing.T) {
	values := map[string]bool{
		"#t\r\n": true,
		"#f\r\n": false,
	}

	for raw, expected := range values {
		parsed := parser.ParseCommand(raw)

		if parsed.Type != Boolean {
			t.Errorf("expected Type to be Boolean, got '%s' instead", getType(parsed.Type))
		}

		if parsed.BoolVal == nil {
			t.Error("expected BoolVal to be 'true', got 'nil' instead")
		}

		if expected {
			if !*parsed.BoolVal {
				t.Error("expected BoolVal to be 'true', got 'false' instead")
			}
		} else {
			if *parsed.BoolVal {
				t.Error("expected BoolVal to be 'false', got 'true' instead")
			}
		}

		respString := parsed.ToRESP()

		if respString != raw {
			t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
		}
	}

}

func TestParseDouble(t *testing.T) {
	values := map[string]float64{
		",1.23\r\n":    1.23,
		",+1.23\r\n":   1.23,
		",-1.23\r\n":   -1.23,
		",10\r\n":      10,
		",inf\r\n":     math.Inf(1),
		",-inf\r\n":    math.Inf(-1),
		",1.23E10\r\n": 1.23e10,
		",1.23e10\r\n": 1.23e10,
	}

	for raw, expected := range values {
		parsed := parser.ParseCommand(raw)

		if parsed.Type != Double {
			t.Errorf("expected Type to be Double, got '%s' instead", getType(parsed.Type))
		}

		if parsed.Double != expected {
			t.Errorf("expected Double to be '%f', got '%f' instead", expected, parsed.Double)
		}

		respString := parsed.ToRESP()

		switch raw {
		case ",+1.23\r\n":
			raw = ",1.23\r\n"
		case ",1.23E10\r\n":
			raw = ",1.23e10\r\n"
		}

		if respString != raw {
			t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
		}
	}

	raw := ",nan\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != Double {
		t.Errorf("expected Type to be Double, got '%s' instead", getType(parsed.Type))
	}

	if !math.IsNaN(parsed.Double) {
		t.Error("expected Double to be NaN, but isn't")
	}

	respString := parsed.ToRESP()

	if respString != raw {
		t.Errorf("expected ToRESP() to return '%s', got '%s' instead", raw, respString)
	}
}

func TestParseBigNum(t *testing.T) {
	raw := "(3492890328409238509324850943850943825024385\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != BigNum {
		t.Errorf("expected Type to be BigNum, got '%s' instead", getType(parsed.Type))
	}

	expected, _ := new(big.Int).SetString("3492890328409238509324850943850943825024385", 10)

	if parsed.BigNum.String() != expected.String() {
		t.Errorf("expected BigNum to be '%s', got '%s' instead", expected.String(), parsed.BigNum.String())
	}
}

func TestParseBulkError(t *testing.T) {
	raw := "!21\r\nSYNTAX invalid syntax\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != BulkError {
		t.Errorf("expected Type to be BulkError, got '%s' instead", getType(parsed.Type))
	}

	if parsed.ErrorType != "SYNTAX" {
		t.Errorf("expected ErrorType to be 'SYNTAX', got '%s' instead", parsed.ErrorType)
	}

	if parsed.ErrorMessage != "invalid syntax" {
		t.Errorf("expected ErrorMessage to be 'invalid syntax', got '%s' instead", parsed.ErrorMessage)
	}
}

func TestParseVerbatim(t *testing.T) {
	raw := "=15\r\ntxt:Some string\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != Verbatim {
		t.Errorf("expected Type to be Verbatim, got '%s' instead", getType(parsed.Type))
	}

	if parsed.Encoding != "txt" {
		t.Errorf("expected Encoding to be 'txt', got '%s' instead", parsed.Encoding)
	}

	if parsed.String != "Some string" {
		t.Errorf("expected String to be 'Some string', got '%s' instead", parsed.String)
	}
}

func TestParseMap(t *testing.T) {
	raw := "%2\r\n+first\r\n:1\r\n+second\r\n:2\r\n"
	parsed := parser.ParseCommand(raw)

	if parsed.Type != Map {
		t.Errorf("expected Type to be Map, got '%s' instead", getType(parsed.Type))
	}

	expected := map[string]int64{
		"first":  1,
		"second": 2,
	}

	for parsedKey, parsedValue := range parsed.Map {
		if parsedKey.Type != SimpleString {
			t.Fatalf("expected key.Type to be SimpleString, got '%s' instead", getType(parsed.Type))
		}

		if parsedValue.Type != Integer {
			t.Fatalf("expected value.Type to be Integer, got '%s' instead", getType(parsed.Type))
		}

		expectedValue, ok := expected[parsedKey.String]
		if !ok {
			t.Fatalf("key not found: '%s'", parsedKey.String)
		}

		if expectedValue != parsedValue.Integer {
			t.Fatalf("expected value.Integer to be '%d', got '%d' instead", expectedValue, parsedKey.Integer)
		}
	}
}

func TestParsePingCommand(t *testing.T) {
	raw := "*1\r\n$4\r\nPING\r\n"

	parsed := parser.ParseCommand(raw)

	if parsed.Type != Array {
		t.Errorf("expected Type to be Array, got %s instead", getType(parsed.Type))
	}

	if len(parsed.Array) != 1 {
		t.Errorf("expected Array to have length 1, got %d instead", len(parsed.Array))
	}

	if parsed.Array[0].Type != BulkString {
		t.Errorf("expected Array[0].Type to be BulkString, got %s instead", getType(parsed.Array[0].Type))
	}

	if parsed.Array[0].String != "PING" {
		t.Errorf("expected Array[0].String to be 'PING', got '%s' instead", parsed.Array[0].String)
	}
}
