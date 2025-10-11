package resp

import (
	"fmt"
	"log"
	"math/big"
	"slices"
	"strconv"
	"strings"
)

type Parser struct{}

type Command struct {
	raw    string
	offset int
}

// const validTokens = `[+\-:$*_#,(!=%|~>]`

type RESPType int

const (
	SimpleString RESPType = iota
	SimpleError
	Integer
	BulkString
	Array
	Null
	Boolean
	Double
	BigNum
	BulkError
	Verbatim
	Map
	Attribute
	Set
	Push
)

type RESPValue struct {
	Type         RESPType
	String       string
	Integer      int64
	Array        []*RESPValue
	IsNull       bool
	BoolVal      *bool
	Double       float64
	BigNum       *big.Int
	ErrorType    string
	ErrorMessage string
	Encoding     string
	Map          map[*RESPValue]*RESPValue
}

func (r *RESPValue) GetValue() any {
	switch r.Type {
	case SimpleString:
		return r.String
	case SimpleError:
		return r.ErrorType + " " + r.ErrorMessage
	case Integer:
		return r.Integer
	case BulkString:
		return r.String
	case Array:
		return r.Array
	case Null:
		return r.IsNull
	case Boolean:
		return &r.BoolVal
	case Double:
		return r.Double
	case BigNum:
		return r.BigNum.String()
	case BulkError:
		return r.ErrorType + " " + r.ErrorMessage
	case Verbatim:
		return r.String
	case Map:
		return r.Map
	default:
		log.Fatal("error: unsupported type")
	}

	return nil
}

func GetParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseCommand(cmd string) *RESPValue {
	command := &Command{
		raw:    cmd,
		offset: 0,
	}

	return p.parseValue(command)
}

func (p *Parser) readLine(cmd *Command) string {
	end := strings.Index(cmd.raw[cmd.offset:], "\r\n")

	line := cmd.raw[cmd.offset : cmd.offset+end]

	cmd.offset = cmd.offset + end + 2

	return line
}

func (p *Parser) parseValue(cmd *Command) *RESPValue {
	if !strings.HasSuffix(cmd.raw, "\r\n") {
		log.Fatalf("error: missing \\r\\n suffix: '%v'", cmd.raw)
	}

	validPrefixes := []string{"+", "-", ":", "$", "*", "_", "#", ",", "(", "!", "=", "%"}
	notImplementedPrefixes := []string{"|", "~", ">"}

	prefix := cmd.raw[cmd.offset : cmd.offset+1]
	if !slices.Contains(validPrefixes, prefix) {
		log.Fatalf("error: invalid prefix: %s", prefix)
	}

	if slices.Contains(notImplementedPrefixes, prefix) {
		log.Fatalf("error: data type not implemented: %s", prefix)
	}

	cmd.offset++

	switch prefix {
	case "+":
		return p.parseSimpleString(cmd)
	case "-":
		return p.parseSimpleError(cmd)
	case ":":
		return p.parseInteger(cmd)
	case "$":
		return p.parseBulkString(cmd)
	case "*":
		return p.parseArray(cmd)
	case "_":
		return p.parseNull(cmd)
	case "#":
		return p.parseBoolean(cmd)
	case ",":
		return p.parseDouble(cmd)
	case "(":
		return p.parseBigNum(cmd)
	case "!":
		return p.parseBulkError(cmd)
	case "=":
		return p.parseVerbatim(cmd)
	case "%":
		return p.parseMap(cmd)
	// case "|":
	// 	return p.parseAttribute(cmd)
	// case "~":
	// 	return p.parseSet(cmd)
	// case ">":
	// 	return p.parsePush(cmd)

	default:
		fmt.Printf("error: invalid RESP type: '%s'\n", prefix)
		return nil
	}
}

func (p *Parser) parseSimpleString(cmd *Command) *RESPValue {
	return &RESPValue{
		Type:   SimpleString,
		String: p.readLine(cmd),
	}
}

func (p *Parser) parseSimpleError(cmd *Command) *RESPValue {
	errorParts := strings.SplitN(p.readLine(cmd), " ", 2)

	return &RESPValue{
		Type:         SimpleError,
		ErrorType:    errorParts[0],
		ErrorMessage: errorParts[1],
	}
}

func (p *Parser) parseInteger(cmd *Command) *RESPValue {
	integer, err := strconv.Atoi(p.readLine(cmd))
	if err != nil {
		log.Fatal("error: unable to parse integer")
	}

	return &RESPValue{
		Type:    Integer,
		Integer: int64(integer),
	}
}

func (p *Parser) parseBulkString(cmd *Command) *RESPValue {
	len, err := strconv.Atoi(p.readLine(cmd))
	if err != nil {
		log.Fatal("error: invalid length")
	}

	if len == -1 {
		return &RESPValue{
			Type:   Null,
			IsNull: true,
		}
	}

	str := cmd.raw[cmd.offset : cmd.offset+len]
	cmd.offset += len + 2

	return &RESPValue{
		Type:   BulkString,
		String: str,
	}
}

func (p *Parser) parseArray(cmd *Command) *RESPValue {
	rawCount := p.readLine(cmd)
	count, err := strconv.Atoi(rawCount)
	if err != nil {
		log.Fatalf("error: invalid count: '%s'", rawCount)
	}

	if count == -1 {
		return &RESPValue{
			Type:   Null,
			IsNull: true,
		}
	}

	respValue := &RESPValue{
		Type:  Array,
		Array: make([]*RESPValue, 0),
	}

	for range count {
		respValue.Array = append(respValue.Array, p.parseValue(cmd))
	}

	return respValue
}

func (p *Parser) parseNull(_ *Command) *RESPValue {
	return &RESPValue{
		Type:   Null,
		IsNull: true,
	}
}

func (p *Parser) parseBoolean(cmd *Command) *RESPValue {
	val := p.readLine(cmd)

	resp := &RESPValue{
		Type: Boolean,
	}

	_t := true
	_f := false

	switch val {
	case "t":
		resp.BoolVal = &_t
	case "f":
		resp.BoolVal = &_f
	default:
		log.Fatal("error: invalid boolean value")
	}

	return resp
}

func (p *Parser) parseDouble(cmd *Command) *RESPValue {
	str := p.readLine(cmd)

	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatalf("error: invalid value for float: %s", str)
	}

	return &RESPValue{
		Type:   Double,
		Double: val,
	}
}

func (p *Parser) parseBigNum(cmd *Command) *RESPValue {
	str := p.readLine(cmd)

	bignum, ok := new(big.Int).SetString(str, 10)
	if !ok {
		log.Fatalf("error: invalid bignum: %s", str)
	}

	return &RESPValue{
		Type:   BigNum,
		BigNum: bignum,
	}
}

func (p *Parser) parseBulkError(cmd *Command) *RESPValue {
	// Not using 'length'
	_ = p.readLine(cmd)

	errorParts := strings.SplitN(p.readLine(cmd), " ", 2)

	return &RESPValue{
		Type:         BulkError,
		ErrorType:    errorParts[0],
		ErrorMessage: errorParts[1],
	}
}

func (p *Parser) parseVerbatim(cmd *Command) *RESPValue {
	// Not using 'length'
	_ = p.readLine(cmd)

	str := p.readLine(cmd)
	parts := strings.SplitN(str, ":", 2)

	return &RESPValue{
		Type:     Verbatim,
		Encoding: parts[0],
		String:   parts[1],
	}
}

func (p *Parser) parseMap(cmd *Command) *RESPValue {
	count := p.readLine(cmd)

	resp := make(map[*RESPValue]*RESPValue)

	for range count {
		key := p.parseValue(cmd)
		value := p.parseValue(cmd)

		resp[key] = value
	}

	return &RESPValue{
		Type: Map,
		Map:  resp,
	}
}

func (p *Parser) ToResponse() string {
	return ""
}

func getType(respType RESPType) string {
	types := []string{
		"SimpleString",
		"SimpleError",
		"Integer",
		"BulkString",
		"Array",
		"Null",
		"Boolean",
		"Double",
		"BigNum",
		"BulkError",
		"Verbatim",
		"Map",
		"Attribute",
		"Set",
		"Push",
	}

	t := int(respType)

	if t < 0 || t > len(types) {
		log.Fatalf("error: invalid type %d", t)
	}

	return types[t]
}
