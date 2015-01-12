package context

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Success(rw http.ResponseWriter, statusCode int, response interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	rw.Write(jsonResponse)
	return nil
}

func Errorf(rw http.ResponseWriter, statusCode int, message string, args ...interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	response := map[string]string{
		"error": fmt.Sprintf(message, args...),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	rw.Write(jsonResponse)
	return nil
}

func OK(rw http.ResponseWriter, response interface{}) error {
	return Success(rw, http.StatusOK, response)
}

func Created(rw http.ResponseWriter, response interface{}) error {
	return Success(rw, http.StatusCreated, response)
}

func NoContent(rw http.ResponseWriter) error {
	return Success(rw, http.StatusNoContent, "")
}

func BadRequest(rw http.ResponseWriter, message string, args ...interface{}) error {
	return Errorf(rw, http.StatusBadRequest, message, args...)
}

func InternalServerError(rw http.ResponseWriter, message string, args ...interface{}) error {
	return Errorf(rw, http.StatusInternalServerError, message, args...)
}
