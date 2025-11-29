package httpUtils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/FernandoVT10/go-auth/app/utils"
)

type HttpError struct {
    Code int
    Message string
}

func (e *HttpError) Error() string {
    return fmt.Sprintf("%s with status code %d", e.Message, e.Code)
}

func InternalServerError() error {
    return &HttpError{
        Code: http.StatusInternalServerError,
        Message: "Internal Server Error",
    }
}

func HandleError(w http.ResponseWriter, err error) {
    errTyped, ok := err.(*HttpError)
    if !ok {
        SendErrorMsg(w, http.StatusInternalServerError, "Internal Server Error")
        return
    }

    SendErrorMsg(w, errTyped.Code, errTyped.Message)
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, data any) {
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    err := json.NewEncoder(w).Encode(data)

    if err != nil {
        utils.LogError("Couldn't send json response: %s", err)
    }
}

func SendErrorMsg(w http.ResponseWriter, statusCode int, msg string) {
    json := map[string]string{"error": msg}
    SendJSONResponse(w, statusCode, json)
}

func handleError(err error) error {
    switch typedErr := err.(type) {
    case *json.UnmarshalTypeError:
        expectedType := typedErr.Type.String()
        // for better comprehension for the user
        if expectedType == "int" {
            expectedType = "number"
        }

        msg := fmt.Sprintf("\"%s\" should be a %s", typedErr.Field, expectedType)
        return &HttpError{
            Code: http.StatusBadRequest,
            Message: msg,
        }
    case *json.SyntaxError:
        return &HttpError{
            Code: http.StatusBadRequest,
            Message: "Invalid JSON syntax",
        }
    default:
        utils.LogError(err.Error())
        return InternalServerError()
    }
}

func ParseJSON(r *http.Request, data any) error {
    if r.Header.Get("Content-Type") != "application/json" {
        return &HttpError {
            Code: http.StatusBadRequest,
            Message: "Body data should be in JSON format",
        }
    }

    err := json.NewDecoder(r.Body).Decode(data)

    if err != nil {
        return handleError(err)
    }

    return nil
}

func GetAuthToken(r *http.Request) (string, error) {
    header := r.Header.Get("Authorization")
    if header == "" {
        return "", &HttpError{
            Code: http.StatusBadRequest,
            Message: "Authorization header is required",
        }
    }

    splitted := strings.Split(header, " ")
    if len(splitted) != 2 || splitted[0] != "Bearer" {
        return "", &HttpError{
            Code: http.StatusBadRequest,
            Message: "Authorization header should be in the form: Bearer <token>",
        }
    }

    return splitted[1], nil
}
