package httpUtils

import (
    "fmt"
    "net/http"
    "encoding/json"
    "errors"
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, data any) {
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    err := json.NewEncoder(w).Encode(data)

    if err != nil {
        fmt.Printf("[ERROR] Couldn't send json response: %s\n", err.Error())
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

        return fmt.Errorf("\"%s\" should be a %s", typedErr.Field, expectedType)
    case *json.SyntaxError:
        return errors.New("Invalid JSON syntax")
    default:
        fmt.Printf("[ERROR] %w", err)
        return errors.New("Invalid JSON syntax")
    }
}

func ParseJSON(r *http.Request, data any) error {
    if r.Header.Get("Content-Type") != "application/json" {
        return errors.New("Body data should be in JSON format")
    }

    err := json.NewDecoder(r.Body).Decode(data)

    if err != nil {
        return handleError(err)
    }

    return nil
}

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
