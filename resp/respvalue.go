package resp

import (
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"
)

const LB = "\r\n"

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

type RESPPrefix string

const (
	SimpleStringPrefix RESPPrefix = "+"
	SimpleErrorPrefix  RESPPrefix = "-"
	IntegerPrefix      RESPPrefix = ":"
	BulkStringPrefix   RESPPrefix = "$"
	ArrayPrefix        RESPPrefix = "*"
	NullPrefix         RESPPrefix = "_"
	BooleanPrefix      RESPPrefix = "#"
	DoublePrefix       RESPPrefix = ","
	BigNumPrefix       RESPPrefix = "("
	BulkErrorPrefix    RESPPrefix = "!"
	VerbatimPrefix     RESPPrefix = "="
	MapPrefix          RESPPrefix = "%"
	SetPrefix          RESPPrefix = "~"
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

func (r *RESPValue) ToRESP() string {
	switch r.Type {
	case SimpleString:
		return string(SimpleStringPrefix) + r.String + LB

	case SimpleError:
		return string(SimpleErrorPrefix) + r.ErrorType + " " + r.ErrorMessage + LB

	case Integer:
		return string(IntegerPrefix) + fmt.Sprintf("%d", r.Integer) + LB

	case BulkString:
		if r.IsNull {
			return string(BulkStringPrefix) + "-1" + LB
		}

		length := len(r.String)
		return string(BulkStringPrefix) + fmt.Sprintf("%d", length) + LB + r.String + LB

	case Array:
		if r.IsNull {
			return string(ArrayPrefix) + "-1" + LB
		}

		count := len(r.Array)
		if count == 0 {
			return string(ArrayPrefix) + "0" + LB
		}

		str := string(ArrayPrefix) + fmt.Sprintf("%d", count) + LB
		for i := range count {
			str += r.Array[i].ToRESP()
		}

		return str

	case Null:
		return string(NullPrefix) + LB

	case Boolean:
		if *r.BoolVal {
			return string(BooleanPrefix) + "t" + LB
		} else {
			return string(BooleanPrefix) + "f" + LB
		}

	case Double:
		switch {
		case math.IsInf(r.Double, 1):
			return ",inf\r\n"
		case math.IsInf(r.Double, -1):
			return ",-inf\r\n"
		case math.IsNaN(r.Double):
			return ",nan\r\n"
		default:
			val := fmt.Sprintf("%g", r.Double)
			val = strings.ReplaceAll(val, "e+", "e")
			return strings.ToLower(string(DoublePrefix) + val + LB)
		}

	case BigNum:
		return string(BigNumPrefix) + r.BigNum.String() + LB

	case BulkError:
		length := len(r.ErrorType) + len(r.ErrorMessage) + 1
		return string(BulkErrorPrefix) + fmt.Sprintf("%d", length) + LB + r.ErrorType + " " + r.ErrorMessage + LB

	case Verbatim:
		length := len(r.String) + 4
		return string(VerbatimPrefix) + fmt.Sprintf("%d", length) + LB + r.Encoding + ":" + r.String + LB

	case Map:
		count := len(r.Map)

		str := string(MapPrefix) + fmt.Sprintf("%d", count) + LB
		for key, val := range r.Map {
			str += key.ToRESP() + val.ToRESP()
		}

		return str
	}

	return ""
}
