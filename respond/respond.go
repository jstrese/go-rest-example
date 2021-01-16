package respond

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type RespondResponse struct {
	StatusCode int    `json:"StatusCode"`
	Message    string `json:"Message"`
}

func New(args ...interface{}) {
	fmt.Sprintf("Bad Request [%s]: %s", args[0], args[1])
}

func isValidStatusCode(customCode interface{}) (code int, out error) {
	code = 400
	newCode := customCode.(int)

	if newCode >= 100 && newCode <= 599 {
		code = newCode
		return
	} else {
		out = errors.New("Status code must be >= 100 and <= 599")
		return
	}
}

func parseString(val interface{}) string {
	switch val.(type) {
	case string:
		return val.(string)
	case int:
		return strconv.Itoa(val.(int))
	case bool:
		return strconv.FormatBool(val.(bool))
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case uint, uint8, uint16, uint32, uint64, int8,
		int16, int32, int64, uintptr:
		return fmt.Sprintf("%d", val)
	default:
		return fmt.Sprintf("%s", val)
	}
}

// BadRequest( w, message, statusCode )
func BadRequest(w http.ResponseWriter, args ...interface{}) {
	res := RespondResponse{
		StatusCode: 400,
		Message:    "Bad Request",
	}

	for i, arg := range args {
		switch i {
		case 0:
			res.Message = parseString(arg)
		case 1:
			if code, out := isValidStatusCode(arg); out == nil {
				res.StatusCode = code
			} else {
				panic(out)
			}
		default:
			panic(fmt.Sprintf("Respond.BadRequest: Too many parameters (%d) expected 2", len(args)))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}

// NotFound( w, message? )
func NotFound(w http.ResponseWriter, args ...interface{}) {
	res := RespondResponse{
		StatusCode: 404,
		Message:    "Not Found",
	}

	for i, arg := range args {
		switch i {
		case 0:
			res.Message = parseString(arg)
		default:
			panic(fmt.Sprintf("Respond.NotFound: Too many parameters (%d) expected 1", len(args)))
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res)
}

// Continue( w )
func Continue(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(100)
}

// Set ( w, statusCode )
func Set(w http.ResponseWriter, code int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	// Make sure the code is within range since there is no default
	if _, out := isValidStatusCode(code); out != nil {
		panic(out)
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
}

// Response( w, statusCode, message? )
func Response(w http.ResponseWriter, code int, args ...interface{}) {
	var ResponseCode int = code
	var ResponseMessage string

	// Make sure the code is within range since there is no default
	if _, out := isValidStatusCode(code); out != nil {
		panic(out)
	}

	// Make sure we only have 1 other argument
	if len(args) > 1 {
		panic(fmt.Sprintf("Respond.Response: Too many parameters (%d) expected 1", len(args)))
	} else if len(args) == 1 {
		ResponseMessage = parseString(args[0])
	}

	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(ResponseCode)

	if len(ResponseMessage) > 0 {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(
			RespondResponse{
				StatusCode: ResponseCode,
				Message:    ResponseMessage,
			})
	}
}
