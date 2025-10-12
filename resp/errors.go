package resp

type ErrorType string

const (
	genericError   = "ERR"
	syntaxError    = "SYNTAX"
	wrongtypeError = "WRONGTYPE"
)

func GenericError(message string) string {
	resp := &RESPValue{
		Type:         SimpleError,
		ErrorType:    genericError,
		ErrorMessage: message,
	}

	return resp.ToRESP()
}

func SyntaxError(message string) string {
	resp := &RESPValue{
		Type:         SimpleError,
		ErrorType:    syntaxError,
		ErrorMessage: message,
	}

	return resp.ToRESP()
}

func WrongTypeError(message string) string {
	resp := &RESPValue{
		Type:         SimpleError,
		ErrorType:    wrongtypeError,
		ErrorMessage: message,
	}

	return resp.ToRESP()
}
